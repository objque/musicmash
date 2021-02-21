package open

import (
	"fmt"
	"os/exec"

	"github.com/musicmash/musicmash/internal/log"
	"github.com/spf13/cobra"
)

func NewOpenCommand() *cobra.Command {
	var web bool
	cmd := &cobra.Command{
		Use:          "open",
		Short:        "Open releases in the app/web",
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			schema := "itmss"
			if web {
				schema = "https"
			}

			url := fmt.Sprintf("%s://%s/%s", schema, "itunes.apple.com/us/album", args[0])
			log.Debugf("Executing 'open %v'", url)
			//nolint:gosec
			return exec.Command("open", url).Run()
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&web, "web", false, "Open release in the browser")
	return cmd
}
