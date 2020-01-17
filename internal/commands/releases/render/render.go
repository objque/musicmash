package render

import (
	"fmt"
	"os"

	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/olekukonko/tablewriter"
)

var headers = []string{"id", "artist_id", "released", "poster", "title", "itunes_id", "spotify_id", "deezer_id"}

func Releases(releases []*releases.Release, showName bool) error {
	table := tablewriter.NewWriter(os.Stdout)
	if showName {
		headers[1] = "artist_name"
	}
	table.SetHeader(headers)
	table.SetAutoFormatHeaders(false)
	for i := range releases {
		row := []string{fmt.Sprint(releases[i].ID)}

		if showName {
			row = append(row, fmt.Sprint(releases[i].ArtistName))
		} else {
			row = append(row, fmt.Sprint(releases[i].ArtistID))
		}

		row = append(row,
			releases[i].Released.Format("2006-01-02"),
			releases[i].Poster,
			releases[i].Title,
			releases[i].ItunesID,
			releases[i].SpotifyID,
			releases[i].DeezerID,
		)
		table.Append(row)
	}
	table.Render()
	return nil
}
