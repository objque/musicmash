package releases

import (
	"fmt"
	"time"

	"github.com/musicmash/musicmash/internal/commands/releases/render"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/spf13/cobra"
)

const layout = "2006-01-02T15:04:05"

func NewListCommand() *cobra.Command {
	var since string
	cmd := &cobra.Command{
		Use:          "list [OPTIONS]",
		Short:        "List releases",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("http://%v:%v", config.Config.HTTP.IP, config.Config.HTTP.Port)
			t, err := time.Parse(layout, since)
			if err != nil {
				return fmt.Errorf("since must be in format: %v", layout)
			}

			result, err := releases.Get(api.NewProvider(url, 1), t)
			if err != nil {
				return err
			}

			if len(result) == 0 {
				fmt.Println(fmt.Sprintf("No releases were found since %s", since))
				return nil
			}

			return render.Releases(result)
		},
	}

	now := time.Now().UTC().Truncate(time.Hour * 24).Format(layout)
	flags := cmd.Flags()
	flags.StringVarP(&since, "since", "s", now, "Filter releases by date")

	return cmd
}
