package artists

import (
	"fmt"

	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/artists"
	"github.com/spf13/cobra"
)

func NewSearchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "search <artist_name>",
		Short:        "Search artists",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("http://%v:%v", config.Config.HTTP.IP, config.Config.HTTP.Port)
			result, err := artists.Search(api.NewProvider(url, 1), args[0])
			if err != nil {
				return err
			}

			if len(result) == 0 {
				fmt.Println("Artists not found")
				return nil
			}

			fmt.Println("id\tname")
			for _, artist := range result {
				fmt.Println(fmt.Sprintf("%d\t%v", artist.ID, artist.Name))
			}
			return nil
		},
	}
	return cmd
}
