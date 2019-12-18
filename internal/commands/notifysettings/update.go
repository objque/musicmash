package notifysettings

import (
	"fmt"

	"github.com/musicmash/musicmash/internal/commands/notifysettings/render"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/notifysettings"
	"github.com/spf13/cobra"
)

func NewUpdateCommand() *cobra.Command {
	setting := notifysettings.Settings{}
	cmd := &cobra.Command{
		Use:          "update [OPTIONS] <username>",
		Short:        "Update new notifications settings",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("http://%v:%v", config.Config.HTTP.IP, config.Config.HTTP.Port)
			err := notifysettings.Update(api.NewProvider(url, 1), args[0], &setting)
			if err != nil {
				return err
			}

			return render.Setting(&setting)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&setting.Service, "service", "telegram", "Service which will be used for notifications")
	flags.StringVar(&setting.Data, "data", "", "Data for notifier. e.g chat_id for telegram")
	_ = cmd.MarkFlagRequired("data")
	return cmd
}
