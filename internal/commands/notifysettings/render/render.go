package render

import (
	"os"

	"github.com/musicmash/musicmash/pkg/api/notifysettings"
	"github.com/olekukonko/tablewriter"
)

var headers = []string{"service", "data"}

func Setting(artist *notifysettings.Settings) error {
	return Settings([]*notifysettings.Settings{artist})
}

func Settings(settings []*notifysettings.Settings) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetAutoFormatHeaders(false)
	for i := range settings {
		table.Append([]string{
			settings[i].Service,
			settings[i].Data,
		})
	}
	table.Render()
	return nil
}
