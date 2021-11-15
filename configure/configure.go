package configure

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

const (
	configDir = "etc"
	ssqConfigFilename = "ssq.yaml"
)

// ServerConfig 服务配置
type ServerConfig struct {
	Port int `yaml:"port"`
}

// LogConfig 日志配置
type LogConfig struct {
	Filename string `yaml:"filename"`
}

// MysqlConfig 数据库配置
type MysqlConfig struct {
	Url string `yaml:"url"`
}

// LocalConfigYaml 配置结构
type LocalConfigYaml struct {
	Reload bool
	Server ServerConfig `yaml:"server"`
	Log LogConfig `yaml:"log"`
	Mysql MysqlConfig `yaml:"mysql"`
}

// GlobalConfig 全局配置
var GlobalConfig LocalConfigYaml

var RootDir string

// 加载配置
func loadConfig() {
	file, err := os.Open(path.Join(configDir, ssqConfigFilename))
	if err != nil {
		panic("open config file error:" + err.Error())
	}
	defer file.Close()

	configFileData, err := ioutil.ReadAll(file)
	if err != nil {
		panic("read config file error: " + err.Error())
	}

	localConfigYaml := new(LocalConfigYaml)
	if err := yaml.Unmarshal(configFileData, localConfigYaml); err != nil {
		panic("yaml.Unmarshal error:" + err.Error())
	}

	GlobalConfig = *localConfigYaml
}

// 获取服务根目录
func getRootDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic("getRootDir error: " + err.Error())
	}
	dir = filepath.Join(dir, "../") + "/"
	return dir
}

func init() {
	RootDir = getRootDir()
	loadConfig()
}