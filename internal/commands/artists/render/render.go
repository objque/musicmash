package render

import (
	"fmt"
	"os"

	"github.com/musicmash/musicmash/pkg/api/artists"
	"github.com/olekukonko/tablewriter"
)

var (
	artistHeaders = []string{"id", "name", "poster", "followers", "popularity"}
	albumHeaders  = []string{"id", "name"}
)

func Artist(artist *artists.Artist) error {
	return Artists([]*artists.Artist{artist})
}

func Artists(artists []*artists.Artist) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(artistHeaders)
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

func Albums(albums []*artists.Album) error {
	_, _ = fmt.Fprintln(os.Stdout, "Albums:")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(albumHeaders)
	table.SetAutoFormatHeaders(false)
	for i := range albums {
		table.Append([]string{
			fmt.Sprintf("%v", albums[i].ID),
			albums[i].Name,
		})
	}
	table.Render()
	return nil
}
