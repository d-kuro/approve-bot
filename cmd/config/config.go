package config

import (
	"io/ioutil"

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

func GetConfig(cfgPath string) (*ApproveConfig, error) {
	buf, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}
	var cfg ApproveConfig
	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
