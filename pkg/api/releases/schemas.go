package releases

import (
	"encoding/json"
	"fmt"
	"time"
)

type Release struct {
	ID          uint64
	ArtistID    int64
	ArtistName  string
	Released    time.Time
	Poster      string
	Title       string
	SpotifyID   string
	Type        string
	DurationMs  int64
	TracksCount int32
	Explicit    bool
}

func (r *Release) UnmarshalJSON(bytes []byte) error {
	var release = struct {
		Released    string `json:"released"`
		ArtistName  string `json:"artist_name"`
		Poster      string `json:"poster"`
		Title       string `json:"title"`
		SpotifyID   string `json:"spotify_id"`
		Type        string `json:"type"`
		DurationMs  int64  `json:"duration_ms"`
		ID          uint64 `json:"id"`
		ArtistID    int64  `json:"artist_id"`
		TracksCount int32  `json:"tracks_count"`
		Explicit    bool   `json:"explicit"`
	}{}

	if err := json.Unmarshal(bytes, &release); err != nil {
		return fmt.Errorf("can't unmarshal json: %w", err)
	}

	released, err := time.Parse("2006-01-02", release.Released)
	if err != nil {
		return fmt.Errorf("can't unmarshal released: %w", err)
	}

	// TODO (m.kalinin): is it possible to simplify code?
	r.Released = released
	r.ArtistName = release.ArtistName
	r.Poster = release.Poster
	r.Title = release.Title
	r.SpotifyID = release.SpotifyID
	r.Type = release.Type
	r.DurationMs = release.DurationMs
	r.ID = release.ID
	r.ArtistID = release.ArtistID
	r.TracksCount = release.TracksCount
	r.Explicit = release.Explicit

	return nil
}
