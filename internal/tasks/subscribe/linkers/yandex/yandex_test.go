package yandex

import (
	"net/http"
	"testing"

	"github.com/musicmash/musicmash/internal/db"

	"github.com/stretchr/testify/assert"
)

func Test_YandexLinker_Reserve(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/api/v2.1/handlers/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"yandexuid": "1234276871451297001"}`))
	})
	task := NewLinker(server.URL)

	// action
	task.reserveArtists([]string{"skrillex", "nero"})

	// assert
	assert.Len(t, task.reservedArtists, 2)
}

func Test_YandexLinker_Free(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/api/v2.1/handlers/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"yandexuid": "1234276871451297001"}`))
	})
	task := NewLinker(server.URL)
	artists := []string{"skrillex", "nero"}
	task.reserveArtists(artists)
	assert.Len(t, task.reservedArtists, 2)

	// action
	task.freeReservedArtists(artists)

	// assert
	assert.Len(t, task.reservedArtists, 0)
}

func Test_YandexLinker_Search_AlreadyExists(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/api/v2.1/handlers/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"yandexuid": "1234276871451297001"}`))
	})
	task := NewLinker(server.URL)
	artists := []string{"skrillex", "nero"}
	db.DbMgr.EnsureArtistExistsInStore("skrillex", "yandex", "xyz")
	db.DbMgr.EnsureArtistExistsInStore("nero", "yandex", "zyx")

	// action
	task.SearchArtists(artists)
}

func Test_YandexLinker_Search(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	mux.HandleFunc("/api/v2.1/handlers/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"yandexuid": "1234276871451297001"}`))
	})
	mux.HandleFunc("/handlers/music-search.jsx", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
    "text": "gorgon city",
    "artists": {
        "items": [{
            "id": 817678,
            "name": "Gorgon City"
        }]
    }
}`))
	})
	task := NewLinker(server.URL)

	// action
	task.SearchArtists([]string{"gorgon city"})

	// assert
	artists, err := db.DbMgr.GetArtistsForStore("yandex")
	assert.NoError(t, err)
	assert.Len(t, artists, 1)
}
