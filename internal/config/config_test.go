package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {
	err := Load([]byte(`
stores:
    yandex:
        url: https://music.yandex.ru
        fetch_workers: 5

    itunes:
        url: https://api.music.apple.com
        fetch_workers: 25
        meta:
          token: "dd214b951ab64de0af62d53678750c90"
          region: "us"
`))

	assert.NoError(t, err)
	assert.Len(t, Config.Stores, 2)

	assert.Equal(t, "https://music.yandex.ru", Config.Stores["yandex"].URL)
	assert.Equal(t, 5, Config.Stores["yandex"].FetchWorkers)

	assert.Equal(t, "https://api.music.apple.com", Config.Stores["itunes"].URL)
	assert.Equal(t, 25, Config.Stores["itunes"].FetchWorkers)
	assert.Len(t, Config.Stores["itunes"].Meta, 2)
	assert.Equal(t, "dd214b951ab64de0af62d53678750c90", Config.Stores["itunes"].Meta["token"])
	assert.Equal(t, "us", Config.Stores["itunes"].Meta["region"])
}
