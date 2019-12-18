package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultConfig = `
---
http:
  ip: 127.0.0.1
  port: 8844

db:
  type: sqlite3
  host: ./musicmash.sqlite3
  log: false

log:
  file: ./musicmash.log
  level: INFO

fetcher:
  enabled: true
  delay: 8

stores:
  itunes:
    url: https://api.music.apple.com
    fetch_workers: 5
    release_url: https://itunes.apple.com/us/album/%s
    name: Apple Music
    fetch: true

sentry:
  enabled: false
  key: https://uuid@sentry.io/123456
  environment: production

notifier:
  enabled: true
  delay: 1

proxy:
  enabled: false
  type: socks5
  host: example.com:1080
  user_name: musicmash
  password: 1s9J-9j2sa-Zkks
`

func TestConfig_Load(t *testing.T) {
	// arrange
	config := New()

	// action
	result := &AppConfig{}
	err := result.LoadFromBytes([]byte(defaultConfig))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, config, result)
}
