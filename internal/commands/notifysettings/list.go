package notifysettings

import (
	"fmt"
	"os"

	"github.com/musicmash/musicmash/internal/commands/notifysettings/render"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/notifysettings"
	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "list <user_name>",
		Short:        "List of user notification settings",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("http://%v:%v", config.Config.HTTP.IP, config.Config.HTTP.Port)
			result, err := notifysettings.List(api.NewProvider(url, 1), args[0])
			if err != nil {
				return err
			}

			if len(result) == 0 {
				fmt.Println("User doesn't have settings")
				os.Exit(0)
			}

			return render.Settings(result)
		},
	}
	return cmd
}
