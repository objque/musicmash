package subscriptions

import (
	"net/http"
	"testing"

	v2 "github.com/objque/musicmash/internal/clients/itunes"
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
	assert.NoError(t, db.DbMgr.EnsureUserExists(userID))
	mux.HandleFunc("/v1/catalog/us/search", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
          {
            "results": {
              "artists": {
                "data": [
                  {
                    "attributes": {
                      "name": "Moderat"
                    },
                    "id": "1416749924"
                  }
                ]
              }
            }
          }`))
	})

	// action
	done, stateID := FindArtistsAndSubscribeUserTask(userID, artists, v2.NewProvider(server.URL, "xxx"))
	<-done

	// assert
	state, err := db.DbMgr.GetState(stateID)
	assert.NoError(t, err)
	assert.Equal(t, db.CompleteState, state.Value)
	subs, err := db.DbMgr.FindAllUserSubscriptions(userID)
	assert.NoError(t, err)
	assert.Len(t, subs, 2)
}
