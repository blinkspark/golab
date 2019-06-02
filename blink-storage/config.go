package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Port int `json:"port"`
}

var defaultConfig *Config

func NewDefaltConfig() *Config {
	if defaultConfig == nil {
		defaultConfig = &Config{
			Port: 22333,
		}
	}
	return defaultConfig
}

func ReadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c *Config

	err = json.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func WriteConfig(path string, c *Config) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, os.ModePerm)
}

func ReadConfigWithLazyCreate(path string) (*Config, error) {
	config, err := ReadConfig(configPath)
	if err != nil {
		config = NewDefaltConfig()
		err = WriteConfig(configPath, config)
		if err != nil {
			return nil, err
		}
	}
	return config, err
}
