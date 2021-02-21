package subscriptions

import (
	"fmt"
	"strconv"

	"github.com/musicmash/musicmash/pkg/api/subscriptions"
)

func parseArtists(args []string) ([]*subscriptions.Subscription, error) {
	subs := make([]*subscriptions.Subscription, len(args))
	for i := range args {
		id, err := strconv.ParseInt(args[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("artist_id must be int, but provided: %s", args[i])
		}

		subs[i] = &subscriptions.Subscription{ArtistID: id}
	}
	return subs, nil
}
