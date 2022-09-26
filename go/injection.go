package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type IgnitionConfig struct {
	Storage struct {
		DirectoriesList []struct {
			Path string `json:"path"`
			Mode int    `json:"mode"`
		} `json:"directories"`
		FilesList []struct {
			Path     string `json:"path"`
			Mode     int    `json:"mode"`
			Contents struct {
				Compression string `json:"compression"`
				Source      string `json:"source"`
			} `json:"contents"`
		} `json:"files"`
	} `json:"storage"`
	Systemd struct {
		UnitsList []struct {
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

	read_file, err := os.ReadFile(filename)
	check(err)

	err = json.Unmarshal([]byte(string(read_file)), &ignitionConfig)
	check(err)

	fmt.Printf("ignitionConfig: %v", ignitionConfig)
}
