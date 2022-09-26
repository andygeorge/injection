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
	SystemdUnitPath            string      = "/etc/systemd/system"
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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println()

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

	for _, directoryConfig := range ignitionConfig.Storage.Directories {
		path := directoryConfig.Path
		fmt.Println(path)

		mode := DefaultDirectoryMode
		if directoryConfig.Mode != 0 {
			mode = os.FileMode(directoryConfig.Mode)
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			err = os.MkdirAll(path, mode)
		} else {
			err = os.Chmod(path, mode)
		}
		check(err)
	}

	for _, fileConfig := range ignitionConfig.Storage.Files {
		fmt.Println(fileConfig.Path)

		mode := DefaultFileMode
		if fileConfig.Mode != 0 {
			mode = os.FileMode(fileConfig.Mode)
		}

		targetFile, err := os.OpenFile(fileConfig.Path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, mode)
		check(err)
		fmt.Fprintf(targetFile, "%s", fileConfig.Contents.Source)
		targetFile.Close()
	}

	for _, unitConfig := range ignitionConfig.Systemd.Units {
		filePath := SystemdUnitPath + "/" + unitConfig.Name
		fmt.Println(filePath)

		unitEnabledString := "disable"
		if unitConfig.Enabled {
			unitEnabledString = "enable"
		}

		targetFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, DefaultSystemdUnitFileMode)
		check(err)
		fmt.Fprintf(targetFile, "%s", unitConfig.Contents)

		cmd := exec.Command("systemctl", unitEnabledString, unitConfig.Name)
		err = cmd.Run()
		check(err)
	}
}
