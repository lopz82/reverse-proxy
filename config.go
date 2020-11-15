package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type ConfigFile struct {
	Routes map[string]Route
}

type Route struct {
	Servers  []string
	Strategy string
}

var config ConfigFile = ConfigFile{}

func init() {
	filename := "config.yml"
	err := yaml.Unmarshal(openConfigFile(filename), &config)
	if err != nil {
		log.Fatalf("Error unmarshalling %s", filename)
	}
}

func openConfigFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error opening %s", filename)
		return nil
	}
	return data
}
