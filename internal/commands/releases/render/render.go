package render

import (
	"fmt"
	"os"

	"github.com/musicmash/musicmash/pkg/api/artists"
	"github.com/olekukonko/tablewriter"
)

var headers = []string{"id", "artist_id", "released", "poster", "title", "itunes_id", "spotify_id", "deezer_id"}

func Releases(releases []*artists.Release) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetAutoFormatHeaders(false)
	for i := range releases {
		table.Append([]string{
			fmt.Sprint(releases[i].ID),
			fmt.Sprint(releases[i].ArtistID),
			releases[i].Released.Format("2006-01-02"),
			releases[i].Poster,
			releases[i].Title,
			releases[i].ItunesID,
			releases[i].SpotifyID,
			releases[i].DeezerID,
		})
	}
	table.Render()
	return nil
}
