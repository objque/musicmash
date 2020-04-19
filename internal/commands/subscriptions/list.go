package subscriptions

import (
	"fmt"
	"os"

	"github.com/musicmash/musicmash/internal/commands/subscriptions/render"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/subscriptions"
	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	var showPoster bool
	cmd := &cobra.Command{
		Use:          "list <username>",
		Short:        "List of user subscriptions",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("http://%v:%v", config.Config.HTTP.IP, config.Config.HTTP.Port)
			result, err := subscriptions.List(api.NewProvider(url, 1), args[0])
			if err != nil {
				return err
			}

			if len(result) == 0 {
				fmt.Println("User doesn't have subscriptions")
				os.Exit(0)
			}

			return render.Subscriptions(result, showPoster)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&showPoster, "show-poster", showPoster, "Show poster column")
	return cmd
}
