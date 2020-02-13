package associations

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/associations"
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "create [OPTIONS] <artist_id> <store_name> <store_id>",
		Short:        "Associate artist with store",
		Args:         cobra.ExactArgs(3),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			artistID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				fmt.Println(errors.New("artist id must be int value"))
				os.Exit(2)
			}

			association := associations.Association{ArtistID: artistID, StoreName: args[1], StoreID: args[2]}
			url := fmt.Sprintf("http://%v:%v", config.Config.HTTP.IP, config.Config.HTTP.Port)
			err = associations.Create(api.NewProvider(url, 1), &association)
			if err != nil {
				return err
			}

			fmt.Println("Artist successfully associated")
			return nil
		},
	}
	return cmd
}
