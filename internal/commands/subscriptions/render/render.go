package render

import (
	"fmt"
	"os"

	"github.com/musicmash/musicmash/pkg/api/subscriptions"
	"github.com/olekukonko/tablewriter"
)

var headers = []string{"artist_id"}

func Subscriptions(subscriptions []*subscriptions.Subscription) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetAutoFormatHeaders(false)
	for i := range subscriptions {
		table.Append([]string{
			fmt.Sprintf("%v", subscriptions[i].ArtistID),
		})
	}
	table.Render()
	return nil
}
