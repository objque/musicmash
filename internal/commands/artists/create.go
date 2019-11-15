package artists

import (
	"fmt"

	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/artists"
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {
	var artist artists.Artist
	cmd := &cobra.Command{
		Use:          "create [OPTIONS] <name>",
		Short:        "Create new artist",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("http://%v:%v", config.Config.HTTP.IP, config.Config.HTTP.Port)
			artist.Name = args[0]
			err := artists.Create(api.NewProvider(url, 1), &artist)
			if err != nil {
				return err
			}

			fmt.Println("id\tname")
			fmt.Println(fmt.Sprintf("%d\t%v", artist.ID, artist.Name))
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&artist.Poster, "poster", "", "Url to artist photo")
	flags.UintVar(&artist.Followers, "followers", 0, "Count of artists followers")
	flags.IntVar(&artist.Popularity, "popularity", 0, "Artist popularity")

	return cmd
}
