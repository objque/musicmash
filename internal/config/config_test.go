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

notifier:
  enabled: false
  delay: 30m
  url: http://notify/v1/releases

sentry:
  enabled: false
  key: https://uuid@sentry.io/123456
  environment: production
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
