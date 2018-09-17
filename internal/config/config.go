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
	Store    Store      `yaml:"store"`
	Tasks    Tasks      `yaml:"tasks"`
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
	Workers                    int     `yaml:"workers"`
	CountOfSkippedHoursToFetch float64 `yaml:"count_of_skipped_hours_to_fetch"`
}

type Tasks struct {
	Subscriptions SubscriptionsTask `yaml:"subscriptions"`
}

type SubscriptionsTask struct {
	FindArtistWorkers      int `yaml:"find_artist_workers"`
	SubscribeArtistWorkers int `yaml:"subscribe_artist_workers"`
}

type Store struct {
	URL    string `yaml:"url"`
	Region string `yaml:"region"`
	Token  string `yaml:"token"`
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

func (db *DBConfig) GetConnString() (DBType string, ConnString string) {
	if db.DBType != "mysql" {
		panic("Only mysql is currently supported")
	}
	var connString = fmt.Sprintf(
		"%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=UTC",
		db.DBLogin,
		db.DBPass,
		db.DBHost,
		db.DBName)
	return db.DBType, connString
}
