package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultConfig = `
---
http:
  ip: 10.0.0.240
  port: 8844

db:
  type: mysql
  host: 10.0.0.42
  name: musicmash
  login: musicmash
  pass: musicmash
  log: false

log:
  file: /var/log/musicmash/musicmash.log
  level: INFO

fetching:
  refetch_after_hours: 1
  count_of_skipped_hours: 8

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

notifier:
  telegram_token: 12345:xxxx_yyy_token
  count_of_skipped_hours: 1
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
