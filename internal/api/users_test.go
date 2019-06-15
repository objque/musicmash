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

func TestAPI_Users_Create(t *testing.T) {
	setup()
	defer teardown()

	// action
	body := CreateUserScheme{UserName: testutil.UserObjque}
	buffer, _ := json.Marshal(&body)
	resp, err := http.Post(fmt.Sprintf("%s/users", server.URL), "application/json", bytes.NewReader(buffer))
	defer func() { _ = resp.Body.Close() }()

	// assert
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusCreated)
	_, err = db.DbMgr.FindUserByName(testutil.UserObjque)
	assert.NoError(t, err)
}

func TestAPI_Users_Create_EmptyBody(t *testing.T) {
	setup()
	defer teardown()

	// action
	body := map[string]string{"user_name": ""}
	buffer, _ := json.Marshal(&body)
	resp, err := http.Post(fmt.Sprintf("%s/users", server.URL), "application/json", bytes.NewReader(buffer))
	defer func() { _ = resp.Body.Close() }()

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAPI_Users_Create_AlreadyExists(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	assert.NoError(t, db.DbMgr.EnsureUserExists(testutil.UserObjque))

	// action
	body := CreateUserScheme{UserName: testutil.UserObjque}
	buffer, _ := json.Marshal(&body)
	resp, err := http.Post(fmt.Sprintf("%s/users", server.URL), "application/json", bytes.NewReader(buffer))
	defer func() { _ = resp.Body.Close() }()

	// assert
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
}

func TestAPI_Users_Get(t *testing.T) {
	setup()
	defer teardown()

	// action
	assert.NoError(t, db.DbMgr.EnsureUserExists("objque@me"))
	resp, err := http.Get(fmt.Sprintf("%s/users/%s", server.URL, "objque@me"))
	defer func() { _ = resp.Body.Close() }()

	// assert
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func TestAPI_Users_Get_NotFound(t *testing.T) {
	setup()
	defer teardown()

	// action
	resp, err := http.Get(fmt.Sprintf("%s/users/%s", server.URL, "objque@me"))
	defer func() { _ = resp.Body.Close() }()

	// assert
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)
}
