package associations

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/musicmash/musicmash/internal/commands/associations/render"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/associations"
	"github.com/spf13/cobra"
)

func NewForCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "for <artist_id>",
		Short:        "Show artist associations",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			artistID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				fmt.Println(errors.New("artist id must be int value"))
				os.Exit(2)
			}

			url := fmt.Sprintf("http://%v:%v", config.Config.HTTP.IP, config.Config.HTTP.Port)
			associations, err := associations.List(api.NewProvider(url, 1), &associations.ListOpts{ArtistID: artistID})
			if err != nil {
				return err
			}

			if len(associations) == 0 {
				fmt.Println("Artist has no associations with stores")
				return nil
			}

			return render.Associations(associations)
		},
	}
	return cmd
}
