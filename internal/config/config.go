package config

import (
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

var Config *AppConfig

func New() *AppConfig {
	return &AppConfig{
		HTTP: HTTPConfig{
			IP:   "127.0.0.1",
			Port: 8844,
		},
		DB: DBConfig{
			Type: "sqlite3",
			Host: "./musicmash.sqlite3",
			Log:  false,
		},
		Log: LogConfig{
			File:  "./musicmash.log",
			Level: "INFO",
		},
		Fetching: FetchingConfig{
			RefetchAfterHours:   1,
			CountOfSkippedHours: 8,
		},
		Stores: StoresConfig{
			"itunes": {
				Name:         "Apple Music",
				URL:          "https://api.music.apple.com",
				FetchWorkers: 5,
				ReleaseURL:   "https://itunes.apple.com/us/album/%s",
				ArtistURL:    "https://itunes.apple.com/us/artist/%s",
				Fetch:        true,
			},
		},
		Sentry: SentryConfig{
			Enabled:     false,
			Key:         "https://uuid@sentry.io/123456",
			Environment: "production",
		},
		Notifier: NotifierConfig{
			TelegramToken:       "12345:xxxx_yyy_token",
			CountOfSkippedHours: 1,
		},
	}
}

func (c *AppConfig) LoadFromFile(configPath string) error {
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Error(err)
		return err
	}

	return c.LoadFromBytes(b)
}

func (c *AppConfig) LoadFromBytes(val []byte) error {
	if err := yaml.Unmarshal(val, c); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (db *DBConfig) GetConnString() (dialect, connString string) {
	const mysqlConnectionFormat = "%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=UTC"
	switch db.Type {
	case "mysql":
		return db.Type, fmt.Sprintf(mysqlConnectionFormat, db.Login, db.Pass, db.Host, db.Name)
	case "sqlite3":
		return db.Type, db.Host
	default:
		panic("Only mysql/sqlite3 is currently supported")
	}
}
