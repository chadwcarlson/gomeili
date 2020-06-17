package config

import (
  "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Index struct {
		Name       string `yaml:"name"`
		PrimaryKey string `yaml:"primaryKey"`
    UID        string `yaml:"uid"`
	} `yaml:"index"`
	Attributes struct {
		Displayed  []string `yaml:"displayed"`
		Searchable []string `yaml:"searchable"`
	} `yaml:"attributes"`
	RankingRules []string `yaml:"rankingRules"`
  Synonyms     map[string][]string `yaml:synonyms`
}

// Parses local server config yaml.
func (c *Config) ParseYAML(data []byte) error {
	return yaml.Unmarshal(data, c)
}

// Load the local server config file.
func (p Config) Load(file string) Config {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	if err := config.ParseYAML(data); err != nil {
		log.Fatal(err)
	}

	return config
}
