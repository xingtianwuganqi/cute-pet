package settings

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DataBase string `yaml:"database"`
		Charset  string `yaml:"charset"`
	} `yaml:"database"`

	App struct {
		Port      int    `yaml:"port"`
		Debug     bool   `yaml:"debug"`
		LogLevel  string `yaml:"log_level"`
		SecretKey string `yaml:"secret_key"`
		Env       string `yaml:"env"`
	} `yaml:"app"`
	ApiKeys struct {
		Google   string `yaml:"google"`
		Facebook string `yaml:"facebook"`
	} `yaml:"api_keys"`
}

// 配置开发环境
var env = "local"
var Conf Config

func ConfigEnvironment() {
	configFile := "config/dev.yaml" // 默认使用开发环境配置
	if env == "production" {
		configFile = "config/production.yaml"
	} else if env == "local" {
		configFile = "config/local.yaml"
	} else if env == "test" {
		configFile = "config/test.yaml"
	}
	// 读取配置文件
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		// 处理错误
		panic(err)
	}

	// 解析配置
	var configInfo Config
	if err := yaml.Unmarshal(data, &configInfo); err != nil {
		// 处理错误
		panic(err)
	}
	Conf = configInfo
}
