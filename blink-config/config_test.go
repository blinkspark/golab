package blink_config

import "testing"

func TestJsonConfig_Save(t *testing.T) {
	config := JsonConfig{"Test": "Hello"}
	err := config.Save("config.json")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadConfig(t *testing.T) {
	config := LoadConfig("./config.json")
	if t == nil {
		t.Fail()
	} else {
		t.Log(config)
	}
}

func TestJsonConfig_Get(t *testing.T) {
	config := LoadConfig("config.json")
	str := config.Get("Test").(string)
	if str != "Hello" {
		t.Fail()
	} else {
		t.Log(str)
	}
}

func TestJsonConfig_Set(t *testing.T) {
	config := LoadConfig("config.json")
	config.Set("Test2", "test2")
	str := config.Get("Test2").(string)
	if str != "test2" {
		t.Fail()
	} else {
		t.Log(str)
	}
}
