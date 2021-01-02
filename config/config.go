package config

import (
	"io/ioutil"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Username        string
	Password        string
	Host            string
	Port            int
	RootDatabase    string `yaml:"rootDatabase"`
	ServiceDatabase string `yaml:"serviceDatabase"`
}

func ReadConfig() (*Config, error) {
	buf, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(buf, &config); err != nil {
		panic(err)
	}
	return &config, nil
}
