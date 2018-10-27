package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {
	err := Load([]byte(`
stores:
    - type: yandex
      url: "https://music.yandex.ru"

    - type: itunes
      url: "https://api.music.apple.com"
      meta:
        token: "kpXVCIsInR5cCI6IkpXVCIsInR5cCI6IkpXVCIsInR5cCI6I"
        region: "us"
`))

	assert.NoError(t, err)
	assert.Len(t, Config.Stores, 2)
}
