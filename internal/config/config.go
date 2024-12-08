package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

const yamlConfigFilePath = "config/config.yaml"

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Logging struct {
		Level string `yaml:"level"`
	} `yaml:"logging"`
	InputFilePath string `yaml:"input_file_path"`
}

var (
	configInstance *Config
	configOnce     sync.Once
)

func GetConfig() (*Config, error) {
	var loadErr error
	configOnce.Do(func() {
		file, err := os.ReadFile(yamlConfigFilePath)
		if err != nil {
			loadErr = fmt.Errorf("error reading config file: %w", err)
			return
		}

		configInstance = &Config{}
		if err := yaml.Unmarshal(file, configInstance); err != nil {
			loadErr = fmt.Errorf("error parsing config file: %w", err)
			return
		}
	})

	if loadErr != nil {
		return nil, loadErr
	}

	return configInstance, nil
}
