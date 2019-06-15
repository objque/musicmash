package config

import (
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

var Config *AppConfig

type AppConfig struct {
	HTTP     HTTPConfig        `yaml:"http"`
	DB       DBConfig          `yaml:"db"`
	Log      LogConfig         `yaml:"log"`
	Fetching Fetching          `yaml:"fetching"`
	Stores   map[string]*Store `yaml:"stores"`
	Notifier Notifier          `yaml:"notifier"`
	Sentry   Sentry            `yaml:"sentry"`
}

type HTTPConfig struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

type LogConfig struct {
	File  string `yaml:"file"`
	Level string `yaml:"level"`
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
	RefetchAfterHours   float64 `yaml:"refetch_after_hours"`
	CountOfSkippedHours float64 `yaml:"count_of_skipped_hours"`
}

type Store struct {
	URL          string `yaml:"url"`
	FetchWorkers int    `yaml:"fetch_workers"`
	Meta         Meta   `yaml:"meta"`
	ReleaseURL   string `yaml:"release_url"`
	ArtistURL    string `yaml:"artist_url"`
	Name         string `yaml:"name"`
	Fetch        bool   `json:"fetch"`
}

type Meta map[string]string

type Notifier struct {
	TelegramToken       string  `yaml:"telegram_token"`
	CountOfSkippedHours float64 `yaml:"count_of_skipped_hours"`
}

type Sentry struct {
	Enabled bool   `yaml:"enabled"`
	Key     string `yaml:"key"`
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
