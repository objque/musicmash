package config

import (
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var Config *AppConfig

type AppConfig struct {
	DB       DBConfig   `yaml:"db"`
	Log      LogConfig  `yaml:"log"`
	HTTP     HTTPConfig `yaml:"http"`
	Fetching Fetching   `yaml:"fetching"`
}

type LogConfig struct {
	File          string `yaml:"file"`
	Level         string `yaml:"level"`
	SyslogEnabled bool   `yaml:"syslog_enable"`
}

type DBConfig struct {
	DBType  string `yaml:"dbtype"`
	DBHost  string `yaml:"dbhost"`
	DBName  string `yaml:"dbname"`
	DBLogin string `yaml:"dblogin"`
	DBPass  string `yaml:"dbpass"`
	Log     bool   `yaml:"log"`
}

type HTTPConfig struct {
	IP   string
	Port int
}

type Fetching struct {
	CountOfSkippedHoursToFetch float64
}

func InitConfig(filepath string) error {
	data, err := ioutil.ReadFile(filepath)

	if err != nil {
		return err
	}

	if err := Load(data); err != nil {
		return err
	}

	log.Infof("Config loaded from %v.", filepath)
	return nil
}

func Load(data []byte) error {
	cfg := AppConfig{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}
	Config = &cfg
	return nil
}

func (cfg *AppConfig) GetConnString() (DBType string, ConnString string) {
	if cfg.DB.DBType != "mysql" {
		panic("Only mysql is currently supported")
	}
	var connString = fmt.Sprintf(
		"%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=UTC",
		cfg.DB.DBLogin,
		cfg.DB.DBPass,
		cfg.DB.DBHost,
		cfg.DB.DBName)
	return cfg.DB.DBType, connString
}
