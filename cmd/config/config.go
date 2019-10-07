package config

import (
	"io"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type ApproveConfig struct {
	Repo   string        `yaml:"repo"`
	Owners []OwnerConfig `yaml:"owners"`
}

type OwnerConfig struct {
	Name     string   `yaml:"name"`
	Patterns []string `yaml:"patterns"`
}

func LoadConfigFromFile(file string) (*ApproveConfig, error) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	return LoadConfig(f)
}

func LoadConfig(r io.Reader) (*ApproveConfig, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var cfg ApproveConfig
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
