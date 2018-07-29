package itunes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Parse_Time(t *testing.T) {
	setup()
	defer teardown()

	// action
	res, err := parseTime("Jul 18, 2018")

	// assert
	assert.NoError(t, err)
	assert.Equal(t, 18, res.Day())
	assert.Equal(t, "July", res.Month().String())
	assert.Equal(t, 2018, res.Year())
}
