package releases

import (
	"fmt"
	"time"

	"github.com/jinzhu/now"
	"github.com/musicmash/musicmash/internal/commands/releases/render"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/spf13/cobra"
)

const layout = "2006-01-02"

func processDate(value string) *time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		return nil
	}
	return &t
}

func NewForCommand() *cobra.Command {
	var (
		since      string
		till       string
		showNames  bool
		showPoster bool
	)
	cmd := &cobra.Command{
		Use:          "show",
		Short:        "Show user releases",
		Aliases:      []string{"for"},
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			url := fmt.Sprintf("http://%v:%v", config.Config.HTTP.IP, config.Config.HTTP.Port)
			opts := releases.GetOptions{
				Since: processDate(since),
				Till:  processDate(till),
			}

			result, err := releases.For(api.NewProvider(url, 1), args[0], &opts)
			if err != nil {
				return err
			}
			return render.Releases(result, showNames, showPoster)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&since, "since", now.BeginningOfWeek().Format(layout), "Start date of feed. Format: yyyy-dd-mm")
	flags.StringVar(&till, "till", "", "End date of feed. Format: yyyy-dd-mm")
	flags.BoolVar(&showNames, "names", true, "Replace artist_id with artist_name")
	flags.BoolVar(&showPoster, "show-poster", false, "Show poster column")

	return cmd
}
