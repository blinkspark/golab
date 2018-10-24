package blink_config

import (
	"github.com/blinkspark/golab/util"
	"testing"
)

func TestJsonConfig_Save(t *testing.T) {
	config := JsonConfig{"Test": "Hello"}
	err := config.Save("config.json")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("config.json")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(config)
	}
}

func TestJsonConfig_Get(t *testing.T) {
	config, err := LoadConfig("config.json")
	util.CheckErr(err)
	str := config.Get("Test").(string)
	if str != "Hello" {
		t.Fail()
	} else {
		t.Log(str)
	}
}

func TestJsonConfig_Set(t *testing.T) {
	config, err := LoadConfig("config.json")
	util.CheckErr(err)
	config.Set("Test2", "test2")
	str := config.Get("Test2").(string)
	if str != "test2" {
		t.Fail()
	} else {
		t.Log(str)
	}
}
