package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name       string            `json:"name" yaml:"name"`
	Ports      map[string]string `json:"ports" yaml:"ports"`
	Volumes    map[string]string `json:"volumes" yaml:"volumes"`
	Components []Component       `json:"components" yaml:"components"`
	Commands   map[string]string `json:"commands" yaml:"commands"`
}

func Load(filename string) (Config, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}
	cfg := Config{}
	err = yaml.Unmarshal(bs, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

type Component struct {
	Name        string            `json:"name" yaml:"name"`
	Image       string            `json:"image" yaml:"image"`
	Entrypoint  string            `json:"entrypoint" yaml:"entrypoint"`
	Cmd         []string          `json:"cmd" yaml:"cmd"`
	Mounts      map[string]string `json:"mounts" yaml:"mounts"`
	Environment map[string]string `json:"environment" yaml:"environment"`
	Hosts       map[string]string `json:"hosts" yaml:"hosts"`
}
