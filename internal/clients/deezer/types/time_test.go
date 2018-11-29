package types

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func TestDeezerTime_ParseLayout(t *testing.T) {
	// arrange
	input := map[string]time.Time{
		`"2018-11-14"`: newDate(2018, 11, 14),
		`"2018-11"`:    newDate(2018, 11, 01),
		`"2018"`:       newDate(2018, 01, 01),
		`"0000-00-00"`: newDate(0001, 01, 01),
	}

	for layout, want := range input {
		//action
		actual := Time{}
		assert.NoError(t, json.Unmarshal([]byte(layout), &actual))

		// assert
		assert.Equal(t, want.Year(), actual.Value.Year())
		assert.Equal(t, want.Month(), actual.Value.Month())
		assert.Equal(t, want.Day(), actual.Value.Day())
	}
}
