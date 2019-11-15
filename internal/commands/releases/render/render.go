package render

import (
	"fmt"
	"os"

	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/olekukonko/tablewriter"
)

var headers = []string{"artist_id", "title", "poster", "released", "store_name", "store_id"}

func Releases(releases []*releases.Release) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetAutoFormatHeaders(false)
	for i := range releases {
		table.Append([]string{
			fmt.Sprint(releases[i].ArtistID),
			releases[i].Title,
			releases[i].Poster,
			releases[i].Released.Format("2006-01-02T15:04:05"),
			releases[i].StoreName,
			releases[i].StoreID,
		})
	}
	table.Render()
	return nil
}
