package releases

import (
	"fmt"
	"strconv"

	"github.com/musicmash/musicmash/internal/commands/releases/render"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewByCommand() *cobra.Command {
	var (
		showNames  bool
		showPoster bool
	)
	cmd := &cobra.Command{
		Use:          "by <artist_id>",
		Short:        "List releases from the artist",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("http://%v:%v", config.Config.HTTP.IP, config.Config.HTTP.Port)

			artistID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return errors.New("artist_id should be int")
			}

			result, err := releases.By(api.NewProvider(url, 1), artistID)
			if err != nil {
				return err
			}

			if len(result) == 0 {
				fmt.Println("Artist hasn't released anything yet")
				return nil
			}

			return render.Releases(result, showNames, showPoster)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&showNames, "names", true, "Replace artist_id with artist_name")
	flags.BoolVar(&showPoster, "show-poster", false, "Show poster column")

	return cmd
}
