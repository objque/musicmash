package search

import (
	"fmt"
	"strings"

	"github.com/musicmash/musicmash/internal/commands/artists/render"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/search"
	"github.com/spf13/cobra"
)

func NewSearchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "search <entity>",
		Short:        "Search something",
		Long:         "You may provide artist name as one or many arguments which will be joined into one query",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("http://%v:%v", config.Config.HTTP.IP, config.Config.HTTP.Port)
			result, err := search.Do(api.NewProvider(url, 1), strings.Join(args, " "))
			if err != nil {
				return err
			}

			if len(result.Artists) == 0 {
				fmt.Println("Artists not found")
				return nil
			}

			return render.Artists(result.Artists)
		},
	}
	return cmd
}
