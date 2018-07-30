package itunes

import "time"

type LastRelease struct {
	ID   uint64
	Date time.Time
}

func NewInfo(id, released string) *LastRelease {
	return &LastRelease{}
}

func (r *LastRelease) IsLatest() bool {
	now := time.Now().UTC().Truncate(time.Hour * 24)
	yesterday := now.Add(-time.Hour * 48)
	return r.Date.UTC().After(yesterday)
}
