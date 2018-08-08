package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/objque/musicmash/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestAPI_State(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	const stateID = "xxx"
	assert.NoError(t, db.DbMgr.UpdateState(stateID, db.ProcessingState))

	// action
	resp, err := http.Get(fmt.Sprintf("%s/states/%s", server.URL, stateID))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	state := db.State{}
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&state))
	assert.Equal(t, stateID, state.ID)
	assert.Equal(t, db.ProcessingState, state.Value)
}

func TestAPI_State_NotFound(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	const stateID = "xxx"

	// action
	resp, err := http.Get(fmt.Sprintf("%s/states/%s", server.URL, stateID))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)
}
