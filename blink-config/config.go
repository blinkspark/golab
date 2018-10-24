package blink_config

import (
	"bufio"
	"encoding/json"
	"github.com/blinkspark/golab/util"
	"io/ioutil"
	"os"
)

// JsonConfig main data type of this config pack
type JsonConfig map[string]interface{}

func (c *JsonConfig) Get(key string) (interface{}) {
	return (*c)[key]
}

func (c *JsonConfig) Set(key string, value interface{}) {
	(*c)[key] = value
}

func (c *JsonConfig) Save(path string) (error) {
	data, err := json.Marshal(*c)
	if err != nil {
		return err
	}

	ioutil.WriteFile(path, data, 0666)
	return nil
}

func LoadConfig(path string) (JsonConfig) {
	f, err := os.Open(path)
	util.CheckErr(err)
	defer f.Close()

	reader := bufio.NewReader(f)
	data, err := ioutil.ReadAll(reader)
	util.CheckErr(err)

	var config JsonConfig
	json.Unmarshal(data, &config)
	return config
}
