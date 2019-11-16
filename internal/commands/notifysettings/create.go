package notifysettings

import (
	"fmt"

	"github.com/musicmash/musicmash/internal/commands/notifysettings/render"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/notifysettings"
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {
	var userName string
	cmd := &cobra.Command{
		Use:          "create [OPTIONS] <service> <data>",
		Short:        "Create new notifications settings",
		Args:         cobra.ExactArgs(2),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("http://%v:%v", config.Config.HTTP.IP, config.Config.HTTP.Port)
			setting := notifysettings.Settings{Service: args[0], Data: args[1]}
			err := notifysettings.Create(api.NewProvider(url, 1), userName, &setting)
			if err != nil {
				return err
			}

			return render.Setting(&setting)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&userName, "username", "", "Username which will store this settings")
	_ = cmd.MarkFlagRequired("username")
	return cmd
}
