package config

import (
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var Config *AppConfig

type AppConfig struct {
	HTTP     HTTPConfig `yaml:"http"`
	DB       DBConfig   `yaml:"db"`
	Log      LogConfig  `yaml:"log"`
	Fetching Fetching   `yaml:"fetching"`
	Stores   []*Store   `yaml:"stores"`
	Notifier Notifier   `yaml:"notifier"`
}

type HTTPConfig struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
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

type Fetching struct {
	// TODO (m.kalinin): extraxt workers into the store struct
	Workers                    int     `yaml:"workers"`
	CountOfSkippedHoursToFetch float64 `yaml:"count_of_skipped_hours_to_fetch"`
}

type Store struct {
	Name string `yaml:"type"`
	URL  string `yaml:"url"`
	Meta Meta   `yaml:"meta"`
}

type Meta map[string]string

type Notifier struct {
	TelegramToken       string  `yaml:"telegram_token"`
	CountOfSkippedHours float64 `yaml:"count_of_skipped_hours"`
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

func (db *DBConfig) GetConnString() (dialect, connString string) {
	if db.DBType != "mysql" {
		panic("Only mysql is currently supported")
	}
	connString = fmt.Sprintf(
		"%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=UTC",
		db.DBLogin,
		db.DBPass,
		db.DBHost,
		db.DBName)
	return db.DBType, connString
}
