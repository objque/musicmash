package render

import (
	"os"

	"github.com/musicmash/musicmash/pkg/api/stores"
	"github.com/olekukonko/tablewriter"
)

var headers = []string{"name"}

func Stores(stores []*stores.Store) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetAutoFormatHeaders(false)
	for i := range stores {
		table.Append([]string{
			stores[i].Name,
		})
	}
	table.Render()
	return nil
}
