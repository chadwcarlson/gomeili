package config

import (

)

type Config struct {
	Index struct {
		Name       string `yaml:"name"`
		PrimaryKey string `yaml:"primaryKey"`
	} `yaml:"index"`
	Attributes struct {
		Displayed  []string `yaml:"displayed"`
		Searchable []string `yaml:"searchable"`
	} `yaml:"attributes"`
	RankingRules []string `yaml:"rankingRules"`
	Synonyms     struct {
		Twoway bool `yaml:"twoway"`
		List   []struct {
			Routes      []string `yaml:"routes,omitempty"`
			Services    []string `yaml:"services,omitempty"`
			Application []string `yaml:"application,omitempty"`
			MultiApp    []string `yaml:"multi-app,omitempty"`
			Regions     []string `yaml:"regions,omitempty"`
		} `yaml:"list"`
	} `yaml:"synonyms"`
}
