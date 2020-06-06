package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultConfig = `
---
http:
  ip: 0.0.0.0
  port: 8844

db:
  host: musicmash.db
  port: 5432
  log: false
  auto_migrate: false
  migrations_dir: migrations

log:
  file: musicmash.log
  level: INFO

fetcher:
  enabled: false
  delay: 1h

sentry:
  enabled: false
  key: https://uuid@sentry.io/123456
  environment: production

notifier:
  enabled: false
  delay: 1h

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
