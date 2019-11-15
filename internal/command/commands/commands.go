package commands

import (
	"github.com/musicmash/musicmash/internal/command/artists"
	"github.com/musicmash/musicmash/internal/command/stores"
	"github.com/musicmash/musicmash/internal/command/subscriptions"
	"github.com/spf13/cobra"
)

func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		artists.NewArtistCommand(),
		stores.NewStoreCommand(),
		subscriptions.NewSubscriptionCommand(),
	)
}
