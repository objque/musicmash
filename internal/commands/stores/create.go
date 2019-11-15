package stores

import (
	"fmt"

	"github.com/musicmash/musicmash/internal/commands/stores/render"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/stores"
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "create <name>",
		Short:        "Create new store",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("http://%v:%v", config.Config.HTTP.IP, config.Config.HTTP.Port)
			store, err := stores.Create(api.NewProvider(url, 1), args[0])
			if err != nil {
				return err
			}

			return render.Store(store)
		},
	}
	return cmd
}
