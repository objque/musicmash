package render

import (
	"fmt"
	"os"

	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/olekukonko/tablewriter"
)

var headers = []string{"id", "poster", "released", "artist_id", "title", "itunes_id", "spotify_id", "deezer_id"}

func Releases(releases []*releases.Release, showName bool) error {
	table := tablewriter.NewWriter(os.Stdout)
	if showName {
		headers[3] = "artist_name"
	}
	table.SetHeader(headers)
	table.SetAutoFormatHeaders(false)
	for i := range releases {
		row := []string{fmt.Sprint(releases[i].ID), releases[i].Poster, releases[i].Released.Format("2006-01-02")}

		if showName {
			row = append(row, fmt.Sprint(releases[i].ArtistName))
		} else {
			row = append(row, fmt.Sprint(releases[i].ArtistID))
		}

		row = append(row,
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
