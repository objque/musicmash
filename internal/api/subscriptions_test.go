package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/objque/musicmash/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestAPI_SubscribeUser(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureUserExists(&db.User{ID: "objque"}))
	assert.NoError(t, db.DbMgr.EnsureArtistExists(&db.Artist{Name: "Moderat", StoreID: 00001}))

	// action
	artists := []string{"modeRat"}
	buffer, _ := json.Marshal(&artists)
	resp, err := http.Post(fmt.Sprintf("%s/objque/subscriptions", server.URL), "application/json", bytes.NewReader(buffer))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusAccepted)
	body := map[string]interface{}{}
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Len(t, body, 1)
}

func TestAPI_SubscribeUser_UserNotFound(t *testing.T) {
	setup()
	defer teardown()

	// action
	resp, err := http.Post(fmt.Sprintf("%s/objque/subscriptions", server.URL), "application/json", nil)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)
}
