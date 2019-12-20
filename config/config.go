package config

import (
	"fmt"
	"github.com/widuu/goini"
	"os"
	"strconv"
)

var (
	Timeout    int
	Trace      bool
	configData []map[string]map[string]string
)

func InitConfig(curMode, fileName string) {
	if fileName == "" {
		fileName = fmt.Sprintf("../conf/%s.ini", curMode)
	}
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			panic("configuration file " + fileName + " is not exist")
		}
		panic("configuration file " + fileName + " is privilge mode is not right")
	}

	conf := goini.SetConfig(fileName)
	configData = conf.ReadList()
}

func GetConfig(section string, key string) string {
	for _, v := range configData {
		if _, ok := v[section]; ok {
			return v[section][key]
		}
	}
	return ""
}

func GetConfigInt(section string, key string) int {
	v := GetConfig(section, key)
	if v == "" {
		return 0
	}
	n, _ := strconv.Atoi(v)
	return n
}

func GetConfigBool(section string, key string) bool {
	v := GetConfig(section, key)
	if v == "" {
		return false
	}
	n, _ := strconv.ParseBool(v)
	return n
}

func GetSection(section string) map[string]string {
	for _, v := range configData {
		if _, ok := v[section]; ok {
			return v[section]
		}
	}
	return nil
}
