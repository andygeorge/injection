package main

import (
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

const (
	DefaultDirectoryMode       os.FileMode = 0755
	DefaultFileMode            os.FileMode = 0644
	DefaultSystemdUnitFileMode os.FileMode = 0644
	DefaultSystemdUnitPath     string      = "/etc/systemd/system"
)

type IgnitionConfig struct {
	Storage struct {
		Directories []struct {
			Path string `json:"path"`
			Mode int    `json:"mode"`
		} `json:"directories"`
		Files []struct {
			Path     string `json:"path"`
			Mode     int    `json:"mode"`
			Contents struct {
				Compression string `json:"compression"`
				Source      string `json:"source"`
			} `json:"contents"`
		} `json:"files"`
	} `json:"storage"`
	Systemd struct {
		Units []struct {
			Name     string `json:"name"`
			Enabled  bool   `json:"enabled"`
			Contents string `json:"contents"`
		} `json:"units"`
	} `json:"systemd"`
}

func main() {
	var filename string
	var ignitionConfig IgnitionConfig
	var err error

	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Specify a filename!")
		os.Exit(1)
	} else {
		filename = args[0]
	}

	readFile, err := os.ReadFile(filename)
	handleError(err, "Error reading file")

	err = json.Unmarshal([]byte(string(readFile)), &ignitionConfig)
	handleError(err, "Error unmarshaling json")

	err = WriteDirectories(ignitionConfig)
	handleError(err, "Error in WriteDirectories")

	err = WriteFiles(ignitionConfig)
	handleError(err, "Error in WriteFiles")

	err = WriteUnits(ignitionConfig)
	handleError(err, "Error in WriteUnits")
}

func OpenFile(path string, mode os.FileMode) (*os.File, error) {
	options := os.O_WRONLY | os.O_TRUNC | os.O_CREATE

	return os.OpenFile(path, options, mode)
}

func WriteDirectories(ignitionConfig IgnitionConfig) error {
	var err error

	for _, directoryConfig := range ignitionConfig.Storage.Directories {
		path := directoryConfig.Path
		mode := DefaultDirectoryMode
		fmt.Println(path)

		if directoryConfig.Mode != 0 {
			mode = os.FileMode(directoryConfig.Mode)
		}

		if _, err = os.Stat(path); os.IsNotExist(err) {
			err = os.MkdirAll(path, mode)
		} else {
			err = os.Chmod(path, mode)
		}
	}

	return err
}

func WriteFiles(ignitionConfig IgnitionConfig) error {
	var err error
	var targetFile *os.File

	for _, fileConfig := range ignitionConfig.Storage.Files {
		path := fileConfig.Path
		mode := DefaultFileMode
		compression := fileConfig.Contents.Compression
		source := fileConfig.Contents.Source
		var rawData, unescapedData string
		var decodedGzipData []byte
		fmt.Println(path)

		if fileConfig.Mode != 0 {
			mode = os.FileMode(fileConfig.Mode)
		}

		if compression == "gzip" {
			idx := strings.Index(source, ";base64,")
			rawData = source[idx+8:]

			gz, err := decodeBase64Data(rawData)
			if err != nil {
				return err
			}
			decodedGzipData, err = decodeGzipData(string(gz))
			if err != nil {
				return err
			}
			targetFile, err = OpenFile(path, mode)
			if err != nil {
				return err
			}
			defer targetFile.Close()
			fmt.Fprintf(targetFile, "%s", string(decodedGzipData))
		} else {
			idx := strings.Index(source, ",")
			rawData = source[idx+1:]

			targetFile, err = OpenFile(path, mode)
			if err != nil {
				return err
			}
			defer targetFile.Close()
			unescapedData, err = url.QueryUnescape(rawData)
			if err != nil {
				return err
			}
			fmt.Fprintf(targetFile, "%s", unescapedData)
		}
	}

	return err
}

func WriteUnits(ignitionConfig IgnitionConfig) error {
	var err error
	var targetFile *os.File

	for _, unitConfig := range ignitionConfig.Systemd.Units {
		path := DefaultSystemdUnitPath + "/" + unitConfig.Name
		mode := DefaultSystemdUnitFileMode
		fmt.Println(path)

		unitEnabledString := "disable"
		if unitConfig.Enabled {
			unitEnabledString = "enable"
		}

		targetFile, err = OpenFile(path, mode)
		if err != nil {
			handleError(err, "Error opening file")
			err = nil
			continue
		}
		fmt.Fprintf(targetFile, "%s", unitConfig.Contents)
		defer targetFile.Close()

		cmd := exec.Command("systemctl", unitEnabledString, unitConfig.Name)
		err = cmd.Run()
		if err != nil {
			handleError(err, "Error running systemctl")
			err = nil
			continue
		}
	}

	return err
}

func decodeBase64Data(data string) ([]byte, error) {
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("unable to decode base64: %q", err)
	}

	return decodedData, nil
}

func decodeGzipData(data string) ([]byte, error) {
	reader, err := gzip.NewReader(strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

func handleError(err error, message string) {
	if err != nil {
		if message != "" {
			fmt.Printf("%s", message)
		}
		fmt.Printf(": %q\n", err)
	}
}
