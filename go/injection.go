package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
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
	check(err)

	err = json.Unmarshal([]byte(string(readFile)), &ignitionConfig)
	check(err)

	err = WriteDirectories(ignitionConfig)
	check(err)

	err = WriteFiles(ignitionConfig)
	check(err)

	err = WriteUnits(ignitionConfig)
	check(err)
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
		fmt.Println(path)

		if fileConfig.Mode != 0 {
			mode = os.FileMode(fileConfig.Mode)
		}

		targetFile, err = OpenFile(path, mode)
		fmt.Fprintf(targetFile, "%s", fileConfig.Contents.Source)
		defer targetFile.Close()
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
		fmt.Fprintf(targetFile, "%s", unitConfig.Contents)
		defer targetFile.Close()
		check(err)

		cmd := exec.Command("systemctl", unitEnabledString, unitConfig.Name)
		err = cmd.Run()
		check(err)
	}

	return err
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
