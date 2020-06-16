package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// User defined configuration file struct, allows multiple resources.
type ConfigSet []struct {
	Name        string   `yaml:"name"`
	Type        string   `yaml:"type"`
	Rank        int      `yaml:"rank"`
	URL         string   `yaml:"url"`
	File        string   `yaml:"file"`
	Destination string   `yaml:"destination"`
	Ignore      []string `yaml:"ignore"`
}

// Single resource configuration.
type Config struct {
	Name        string   `yaml:"name"`
	Type        string   `yaml:"type"`
	Rank        int      `yaml:"rank"`
	URL         string   `yaml:"url"`
	File        string   `yaml:"file"`
	Destination string   `yaml:"destination"`
	Ignore      []string `yaml:"ignore"`
}

// Parses local `config.yaml` array.
func (c *ConfigSet) ParseYAML(data []byte) error {
	return yaml.Unmarshal(data, c)
}

// Load the local `config.yaml` file.
func (p ConfigSet) Load(file string) ConfigSet {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	var config ConfigSet
	if err := config.ParseYAML(data); err != nil {
		log.Fatal(err)
	}

	return config
}
