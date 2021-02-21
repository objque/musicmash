package associations

import (
	"fmt"

	"github.com/musicmash/musicmash/internal/commands/associations/render"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/associations"
	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	opts := associations.ListOpts{}
	cmd := &cobra.Command{
		Use:          "list [OPTIONS]",
		Short:        "Show associations",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("http://%v:%v", config.Config.HTTP.IP, config.Config.HTTP.Port)
			associations, err := associations.List(api.NewProvider(url, 1), &opts)
			if err != nil {
				return err
			}

			if len(associations) == 0 {
				fmt.Println("Artists associations not found")
				return nil
			}

			return render.Associations(associations)
		},
	}

	flags := cmd.Flags()
	flags.Int64Var(&opts.ArtistID, "artist-id", 0, "Filter associations by artist_id")
	flags.StringVar(&opts.StoreName, "store-name", "", "Filter associations by store-name")

	return cmd
}
