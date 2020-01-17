package subscriptions

import "github.com/spf13/cobra"

func NewSubscriptionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "subscription",
		Short:   "Manage subscriptions",
		Aliases: []string{"subs"},
		Args:    cobra.NoArgs,
	}
	cmd.AddCommand(
		NewCreateCommand(),
		NewDeleteCommand(),
		NewListCommand(),
	)
	return cmd
}
