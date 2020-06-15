package config

import (
  "log"
  "io/ioutil"
  "gopkg.in/yaml.v2"
)

type ConfigSet []struct {
  Name        string   `yaml:"name"`
  Type        string   `yaml:"type"`
  Rank        int      `yaml:"rank"`
  URL         string   `yaml:"url"`
  File        string   `yaml:"file"`
  Destination string   `yaml:"destination"`
  Ignore      []string `yaml:"ignore"`
}

type Config struct {
  Name        string   `yaml:"name"`
  Type        string   `yaml:"type"`
  Rank        int      `yaml:"rank"`
  URL         string   `yaml:"url"`
  File        string   `yaml:"file"`
  Destination string   `yaml:"destination"`
  Ignore      []string `yaml:"ignore"`
}

func (c *ConfigSet) ParseYAML(data []byte) error {
	return yaml.Unmarshal(data, c)
}

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
