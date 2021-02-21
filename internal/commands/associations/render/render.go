package render

import (
	"fmt"
	"os"

	"github.com/musicmash/musicmash/pkg/api/associations"
	"github.com/olekukonko/tablewriter"
)

var artistHeaders = []string{"artist_id", "store_name", "store_id"}

func Associations(associations []*associations.Association) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(artistHeaders)
	table.SetAutoFormatHeaders(false)
	for i := range associations {
		table.Append([]string{
			fmt.Sprintf("%v", associations[i].ArtistID),
			associations[i].StoreName,
			associations[i].StoreID,
		})
	}
	table.Render()
	return nil
}
