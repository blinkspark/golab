package config

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/blinkspark/golab/util"
)

// InitConfig init config
func InitConfig(configPath string, config interface{}) {
	configFile, err := os.Open(configPath)
	// configFile, err := os.OpenFile(ConfigPath, os.O_RDWR, 0666)
	util.CheckErr(err)
	defer configFile.Close()
	fileReader := bufio.NewReader(configFile)
	str, err := fileReader.ReadString(0)
	err = json.Unmarshal([]byte(str), config)
	util.CheckErr(err)
}
