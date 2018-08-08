package subscriptions

import (
	"net/http"
	"testing"

	"github.com/objque/musicmash/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestSubscriptions_FindArtistsAndSubscribeUserTask(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	const userID = "objque@me"
	artists := []string{"King Curtis", "modeRAT"}
	assert.NoError(t, db.DbMgr.EnsureArtistExists(&db.Artist{Name: "King Curtis", StoreID: 0001}))
	assert.NoError(t, db.DbMgr.EnsureUserExists(&db.User{ID: userID}))
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
          "resultCount": 2,
          "results": [
            {
              "artistId": 3316749924,
              "artistName": "Party Favor & Moderat"
            },
            {
              "artistId": 1416749924,
              "artistName": "Moderat"
            }
          ]
        }`))
	})

	// action
	done, stateID := FindArtistsAndSubscribeUserTask(userID, artists)
	<-done

	// assert
	state, err := db.DbMgr.GetState(stateID)
	assert.NoError(t, err)
	assert.Equal(t, db.CompleteState, state.Value)
}
