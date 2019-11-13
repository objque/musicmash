package commands

import (
	"github.com/musicmash/musicmash/internal/command/stores"
	"github.com/spf13/cobra"
)

func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		stores.NewStoreCommand(),
	)
}
