package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAPI_Chats_Create(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureUserExists(testutil.UserObjque))

	// action
	body := AddUserChatScheme{ChatID: 10004}
	buffer, _ := json.Marshal(&body)
	resp, err := http.Post(fmt.Sprintf("%s/%s/chats", server.URL, testutil.UserObjque), "application/json", bytes.NewReader(buffer))
	defer func() { _ = resp.Body.Close() }()

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	chatID, err := db.DbMgr.FindChatByUserName(testutil.UserObjque)
	assert.NoError(t, err)
	assert.Equal(t, int64(10004), *chatID)
}

func TestAPI_Chats_Create_UserNotFound(t *testing.T) {
	setup()
	defer teardown()

	// action
	body := AddUserChatScheme{ChatID: 10004}
	buffer, _ := json.Marshal(&body)
	resp, err := http.Post(fmt.Sprintf("%s/%s/chats", server.URL, testutil.UserObjque), "application/json", bytes.NewReader(buffer))
	defer func() { _ = resp.Body.Close() }()

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
