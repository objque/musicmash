package subscriptions

import (
	"fmt"
	"strconv"
)

func parseArtists(args []string) ([]int64, error) {
	artists := make([]int64, len(args))
	for i := range args {
		id, err := strconv.ParseInt(args[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("artist_id must be int, but provided: %s", args[i])
		}

		artists[i] = id
	}
	return artists, nil
}
