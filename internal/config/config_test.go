package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {
	err := Load([]byte(`
artists: http://artists

stores:
    yandex:
        url: https://music.yandex.ru
        fetch_workers: 5
        name: "Yandex.Music"
        fetch: false

    itunes:
        url: https://api.music.apple.com
        fetch_workers: 25
        name: "Apple Music"
        fetch: true
        meta:
          token: "dd214b951ab64de0af62d53678750c90"
          region: "us"
notifier:
  count_of_skipped_hours: 8
  telegram_token: "12340255:BBBZZZJJJJJAAAEEEEE"
`))

	assert.NoError(t, err)
	assert.Len(t, Config.Stores, 2)

	assert.Equal(t, "http://artists", Config.Artists)

	assert.Equal(t, "https://music.yandex.ru", Config.Stores["yandex"].URL)
	assert.Equal(t, 5, Config.Stores["yandex"].FetchWorkers)
	assert.Equal(t, "Yandex.Music", Config.Stores["yandex"].Name)
	assert.False(t, Config.Stores["yandex"].Fetch)

	assert.Equal(t, "https://api.music.apple.com", Config.Stores["itunes"].URL)
	assert.Equal(t, 25, Config.Stores["itunes"].FetchWorkers)
	assert.Len(t, Config.Stores["itunes"].Meta, 2)
	assert.Equal(t, "dd214b951ab64de0af62d53678750c90", Config.Stores["itunes"].Meta["token"])
	assert.Equal(t, "us", Config.Stores["itunes"].Meta["region"])
	assert.Equal(t, "us", Config.Stores["itunes"].Meta["region"])
	assert.Equal(t, "Apple Music", Config.Stores["itunes"].Name)
	assert.True(t, Config.Stores["itunes"].Fetch)

	assert.Equal(t, float64(8), Config.Notifier.CountOfSkippedHours)
	assert.Equal(t, "12340255:BBBZZZJJJJJAAAEEEEE", Config.Notifier.TelegramToken)
}
