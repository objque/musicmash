package stores

import (
	"fmt"

	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/stores"
	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "list",
		Short:        "List stores",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("http://%v:%v", config.Config.HTTP.IP, config.Config.HTTP.Port)
			result, err := stores.List(api.NewProvider(url, 1))
			if err != nil {
				return err
			}

			for _, store := range result {
				fmt.Println(store.Name)
			}
			return nil
		},
	}
	return cmd
}
