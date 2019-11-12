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
  refetch_after_hours: 1
  delay: 8

stores:
  itunes:
    url: https://api.music.apple.com
    fetch_workers: 5
    release_url: https://itunes.apple.com/us/album/%s
    artist_url: https://itunes.apple.com/us/artist/%s
    name: Apple Music
    fetch: true

sentry:
  enabled: false
  key: https://uuid@sentry.io/123456
  environment: production

notifier:
  enabled: true
  telegram_token: 12345:xxxx_yyy_token
  delay: 1
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
