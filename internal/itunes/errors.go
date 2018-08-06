package itunes

import "errors"

var (
	ArtistInactiveErr = errors.New("the last release from an artist was long ago or page is broken")
	ArtistNotFoundErr = errors.New("artists not found")
)
