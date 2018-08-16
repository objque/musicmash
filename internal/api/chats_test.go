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

func TestAPI_Chats_Create(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureUserExists("objque@me"))

	// action
	body := AddUserChatScheme{ChatID: 10004}
	buffer, _ := json.Marshal(&body)
	resp, err := http.Post(fmt.Sprintf("%s/objque@me/chats", server.URL), "application/json", bytes.NewReader(buffer))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusCreated)
	chatID, err := db.DbMgr.FindChatByUserID("objque@me")
	assert.NoError(t, err)
	assert.Equal(t, int64(10004), *chatID)
}

func TestAPI_Chats_Create_UserNotFound(t *testing.T) {
	setup()
	defer teardown()

	// action
	body := AddUserChatScheme{ChatID: 10004}
	buffer, _ := json.Marshal(&body)
	resp, err := http.Post(fmt.Sprintf("%s/objque@me/chats", server.URL), "application/json", bytes.NewReader(buffer))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)
}
