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

func NewShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "show <artist_id>",
		Short:        "Show artist info",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("http://%v:%v", config.Config.HTTP.IP, config.Config.HTTP.Port)
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				fmt.Println(errors.New("id must be int value"))
				os.Exit(2)
			}

			artist, err := artists.Get(api.NewProvider(url, 1), id)
			if err != nil {
				return err
			}

			fmt.Println("id\tname")
			fmt.Println(fmt.Sprintf("%d\t%v", artist.ID, artist.Name))
			return nil
		},
	}

	return cmd
}
