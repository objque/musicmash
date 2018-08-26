package yandex

import (
	"net/http"
	"testing"

	"github.com/objque/musicmash/internal/clients/yandex"
	"github.com/stretchr/testify/assert"
)

func TestYandexFinder_RemoveReleaseType(t *testing.T) {
	// arrange
	cases := []struct {
		In, Out string
	}{
		{In: "MDR (Remixes) - Single", Out: "MDR (Remixes)"},
		{In: "MDR (Remixes) - EP", Out: "MDR (Remixes)"},
	}

	for _, testCase := range cases {
		// action
		result := removeReleaseType(testCase.In)

		// assert
		assert.Equal(t, testCase.Out, result)
	}
}

func TestYandexFinder_SplitArtists(t *testing.T) {
	// arrange
	cases := []struct {
		In  string
		Len int
	}{
		{In: "Wolfgang Muthspiel, Ambrose Akinmusire, Brad Mehldau, Larry Grenadier & Eric Harland", Len: 5},
	}

	for _, testCase := range cases {
		// action
		result := splitArtists(testCase.In)

		// assert
		assert.Len(t, result, testCase.Len)
	}
}

func TestYandexFinder_SearchArtistID(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	// mock yandex auth
	mux.HandleFunc("/api/v2.1/handlers/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"yandexuid": "1234276871451297001"}`))
	})
	// mock search api
	mux.HandleFunc("/handlers/music-search.jsx", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
            "text": "gorgon city",
            "artists": {
                "items": [{
                    "id": 817678,
                    "name": "skrillex"
                },{
                    "id": 678817,
                    "name": "Gorgon City"
                }]
            }
        }`))
	})
	client := yandex.New(server.URL)

	// action
	artistID, err := searchArtistID(client, "gorgon city")

	// assert
	assert.NoError(t, err)
	assert.Equal(t, 678817, artistID)
}

func TestYandexFinder_SearchArtistID_NotFound(t *testing.T) {
	setup()
	defer teardown()

	// arrange
	// mock yandex auth
	mux.HandleFunc("/api/v2.1/handlers/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"yandexuid": "1234276871451297001"}`))
	})
	// mock search api
	mux.HandleFunc("/handlers/music-search.jsx", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
            "text": "gorgon city",
            "artists": {
                "items": [{
                    "id": 817678,
                    "name": "skrillex"
                }]
            }
        }`))
	})
	client := yandex.New(server.URL)

	// action
	artistID, err := searchArtistID(client, "gorgon city")

	// assert
	assert.Error(t, err)
	assert.Equal(t, ArtistNotFoundErr, err)
	assert.Equal(t, 0, artistID)
}
