package config

import (
	"fmt"
	"github.com/widuu/goini"
	"os"
	"strconv"
)

var (
	configData []map[string]map[string]string
)

func InitConfig(curMode string) {
	var fileName = fmt.Sprintf("./conf/%s.ini", curMode)
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

func GetConfigInt64(section string, key string) int64 {
	v := GetConfig(section, key)
	if v == "" {
		return 0
	}
	n, _ := strconv.ParseInt(v, 10, 64)
	return n
}

func GetConfigFloat64(section string, key string) float64 {
	v := GetConfig(section, key)
	if v == "" {
		return 0
	}
	n, _ := strconv.ParseFloat(v, 64)
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

