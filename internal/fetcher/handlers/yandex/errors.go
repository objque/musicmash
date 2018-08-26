package yandex

import "errors"

var (
	ArtistNotFoundErr  = errors.New("artist not found")
	ReleaseNotFoundErr = errors.New("release not found")
)
