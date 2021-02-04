package releases

import (
	"fmt"
	"time"

	"github.com/jinzhu/now"
	"github.com/musicmash/musicmash/internal/commands/releases/render"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/utils/ptr"
	"github.com/musicmash/musicmash/pkg/api"
	"github.com/musicmash/musicmash/pkg/api/releases"
	"github.com/spf13/cobra"
)

func parseDate(date string) (*time.Time, error) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

//nolint:gocyclo,gocognit,golint
func NewListCommand() *cobra.Command {
	var since, till string
	renderOpts := render.Options{}
	// dirty hack cause cobra can't handle nil as default for int like types
	opts := releases.GetOptions{
		Offset: ptr.Uint(0),
		Limit:  ptr.Uint(100),
	}
	cmd := &cobra.Command{
		Use:          "list",
		Short:        "List of releases with powerful filters",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			// if artist_id not provided, also set limit for since
			if opts.Since != nil {
				since = now.BeginningOfWeek().Format("2006-01-02")
			}

			if opts.Since, err = parseDate(since); since != "" && err != nil {
				return err
			}

			if opts.Till, err = parseDate(till); till != "" && err != nil {
				return err
			}

			url := fmt.Sprintf("http://%v:%v", config.Config.HTTP.IP, config.Config.HTTP.Port)
			rels, err := releases.List(api.NewProvider(url, 1), &opts)
			if err != nil {
				return err
			}

			return render.Releases(rels, renderOpts)
		},
	}

	flags := cmd.Flags()
	flags.Uint64Var(opts.Limit, "limit", 100, "Limit of rows")
	flags.Uint64Var(opts.Offset, "offset", 0, "Offset for rows")
	flags.StringVar(&since, "since", "", "Filter by release date. Format: YYYY-MM-DD")
	flags.StringVar(&till, "till", "", "Filter by release date. Format: YYYY-MM-DD")
	flags.StringVar(&opts.SortType, "sort-type", "desc", "Sort type for releases {asc,desc}")
	flags.StringVar(&opts.ReleaseType, "type", "", "Filter by release-type {album,song,music-video}")
	flags.StringVar(&opts.UserName, "user", "", "Filter by user subscriptions")
	flags.BoolVar(&renderOpts.ShowNames, "names", true, "Replace artist_id with artist_name")
	flags.BoolVar(&renderOpts.ShowPosters, "posters", false, "Show poster column")
	// todo: add filter by explicit

	return cmd
}
