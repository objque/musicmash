package artists

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/artists"
	"github.com/spf13/cobra"
)

func NewAssociateCommand() *cobra.Command {
	var association artists.Association
	cmd := &cobra.Command{
		Use:          "associate [OPTIONS] <artist_id>",
		Short:        "Associate artist with store",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			artistID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				fmt.Println(errors.New("artist id must be int value"))
				os.Exit(2)
			}

			association.ArtistID = artistID
			url := fmt.Sprintf("http://%v:%v", config.Config.HTTP.IP, config.Config.HTTP.Port)
			err = artists.Associate(api.NewProvider(url, 1), &association)
			if err != nil {
				return err
			}

			fmt.Println("Artist successfully associated")
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&association.StoreName, "store-name", "", "Store name")
	flags.StringVar(&association.StoreID, "store-id", "", "Artist ID in the store")
	_ = cmd.MarkFlagRequired("store-name")
	_ = cmd.MarkFlagRequired("store-id")

	return cmd
}
