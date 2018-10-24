package blink_config

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
)

// JsonConfig main data type of this config pack
type JsonConfig map[string]interface{}

func (c *JsonConfig) Get(key string) interface{} {
	return (*c)[key]
}

func (c *JsonConfig) Set(key string, value interface{}) {
	(*c)[key] = value
}

func (c *JsonConfig) Save(path string) error {
	data, err := json.Marshal(*c)
	if err != nil {
		return err
	}

	ioutil.WriteFile(path, data, 0666)
	return nil
}

func LoadConfig(path string) (JsonConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var config JsonConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func NewConfig() JsonConfig {
	return JsonConfig{}
}
