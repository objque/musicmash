package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/musicmash/musicmash/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestAPI_Users_Create(t *testing.T) {
	setup()
	defer teardown()

	// action
	body := CreateUserScheme{UserName: "objque@me"}
	buffer, _ := json.Marshal(&body)
	resp, err := http.Post(fmt.Sprintf("%s/users", server.URL), "application/json", bytes.NewReader(buffer))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusCreated)
	_, err = db.DbMgr.FindUserByName("objque@me")
	assert.NoError(t, err)
}

func TestAPI_Users_Create_EmptyBody(t *testing.T) {
	setup()
	defer teardown()

	// action
	body := map[string]string{"user_name": ""}
	buffer, _ := json.Marshal(&body)
	resp, err := http.Post(fmt.Sprintf("%s/users", server.URL), "application/json", bytes.NewReader(buffer))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAPI_Users_Create_AlreadyExists(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureUserExists("objque@me"))

	// action
	body := CreateUserScheme{UserName: "objque@me"}
	buffer, _ := json.Marshal(&body)
	resp, err := http.Post(fmt.Sprintf("%s/users", server.URL), "application/json", bytes.NewReader(buffer))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
}
