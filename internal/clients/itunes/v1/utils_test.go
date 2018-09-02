package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_ParseTime_Layout1(t *testing.T) {
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

func TestClient_ParseTime_Layout2(t *testing.T) {
	setup()
	defer teardown()

	// action
	res, err := parseTime("10 Nov 2017")

	// assert
	assert.NoError(t, err)
	assert.Equal(t, 10, res.Day())
	assert.Equal(t, "November", res.Month().String())
	assert.Equal(t, 2017, res.Year())
}
