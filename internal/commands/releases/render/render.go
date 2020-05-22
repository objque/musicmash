package render

import (
	"fmt"
	"os"

	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/olekukonko/tablewriter"
)

var headers = []string{"id", "poster", "released", "artist_id", "title", "type", "explicit", "itunes_id", "spotify_id", "deezer_id"}

type Options struct {
	ShowNames   bool
	ShowPosters bool
}

//nolint:gocyclo,gocognit,golint
func Releases(releases []*releases.Release, opts Options) error {
	table := tablewriter.NewWriter(os.Stdout)
	if opts.ShowNames {
		headers[3] = "artist_name"
	}
	if !opts.ShowPosters {
		// cut id
		headers = headers[1:]
		// replace poster with id
		headers[0] = "id"
	}
	table.SetHeader(headers)
	table.SetAutoFormatHeaders(false)
	for i := range releases {
		row := []string{fmt.Sprint(releases[i].ID)}

		if opts.ShowPosters {
			row = append(row, fmt.Sprint(releases[i].Poster))
		}

		row = append(row, releases[i].Released.Format("2006-01-02"))

		if opts.ShowNames {
			row = append(row, fmt.Sprint(releases[i].ArtistName))
		} else {
			row = append(row, fmt.Sprint(releases[i].ArtistID))
		}

		var spotifyID, deezerID string
		if releases[i].SpotifyID != nil {
			spotifyID = *releases[i].SpotifyID
		}
		if releases[i].DeezerID != nil {
			deezerID = *releases[i].DeezerID
		}

		row = append(row,
			releases[i].Title,
			releases[i].Type,
			fmt.Sprintf("%v", releases[i].Explicit),
			*releases[i].ItunesID,
			spotifyID,
			deezerID,
		)
		table.Append(row)
	}
	table.Render()
	return nil
}
