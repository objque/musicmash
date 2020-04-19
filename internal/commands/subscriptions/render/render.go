package render

import (
	"fmt"
	"os"

	"github.com/musicmash/musicmash/pkg/api/subscriptions"
	"github.com/olekukonko/tablewriter"
)

var headers = []string{"artist_id", "artist_name"}

func Subscriptions(subscriptions []*subscriptions.Subscription, showPoster bool) error {
	table := tablewriter.NewWriter(os.Stdout)
	if showPoster {
		headers = append(headers, "artist_poster")
	}
	table.SetHeader(headers)
	table.SetAutoFormatHeaders(false)
	for i := range subscriptions {
		row := []string{
			fmt.Sprintf("%v", subscriptions[i].ArtistID),
			subscriptions[i].ArtistName,
		}
		if showPoster {
			row = append(row, subscriptions[i].ArtistPoster)
		}
		table.Append(row)
	}
	table.Render()
	return nil
}
