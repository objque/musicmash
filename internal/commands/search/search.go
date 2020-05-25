package search

import (
	"fmt"
	"strings"

	artsrender "github.com/musicmash/musicmash/internal/commands/artists/render"
	relsrender "github.com/musicmash/musicmash/internal/commands/releases/render"
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

			_ = artsrender.Artists(result.Artists)
			_ = relsrender.Releases(result.Releases, relsrender.Options{ShowNames: true, ShowPosters: false})
			return nil
		},
	}
	return cmd
}
