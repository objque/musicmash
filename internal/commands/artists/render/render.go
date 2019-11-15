package render

import (
	"fmt"
	"os"

	"github.com/musicmash/musicmash/pkg/api/artists"
	"github.com/olekukonko/tablewriter"
)

var headers = []string{"id", "name", "poster", "followers", "popularity"}

func Artist(artist *artists.Artist) error {
	return Artists([]*artists.Artist{artist})
}

func Artists(artists []*artists.Artist) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetAutoFormatHeaders(false)
	for i := range artists {
		table.Append([]string{
			fmt.Sprintf("%v", artists[i].ID),
			artists[i].Name,
			artists[i].Poster,
			fmt.Sprintf("%v", artists[i].Followers),
			fmt.Sprintf("%v", artists[i].Popularity),
		})
	}
	table.Render()
	return nil
}
