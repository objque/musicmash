package subscriptions

import (
	"fmt"
	"os"

	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/subscriptions"
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "create <user_name> [<artist_id> ... <artist_id_n>]",
		Short:        "Subscribe user for artists",
		Args:         cobra.MinimumNArgs(2),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			userName := args[0]
			artists, err := parseArtists(args[1:])
			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}

			url := fmt.Sprintf("http://%v:%v", config.Config.HTTP.IP, config.Config.HTTP.Port)
			err = subscriptions.Create(api.NewProvider(url, 1), userName, artists)
			if err != nil {
				return err
			}

			fmt.Println("User has been subscribed")
			return nil
		},
	}
	return cmd
}
