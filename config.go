package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ConfigFile struct {
	Routes map[string]Route
}

type Route struct {
	Servers  []string
	Strategy string
}

var config = ConfigFile{}

var defaultPath = filepath.Join(getAppPath(), "config/config.yml")

func init() {
	err := yaml.Unmarshal(openConfigFile(defaultPath), &config)
	if err != nil {
		log.Fatalf("Error unmarshalling %s", defaultPath)
	}
}

func getAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}

func openConfigFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error opening %s: %s", filename, err)
		return nil
	}
	return data
}

