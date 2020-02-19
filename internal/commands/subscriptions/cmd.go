package subscriptions

import "github.com/spf13/cobra"

func NewSubscriptionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "subscription",
		Aliases: []string{"subscriptions", "subs", "sub"},
		Short:   "Manage subscriptions",
		Args:    cobra.NoArgs,
	}
	cmd.AddCommand(
		NewCreateCommand(),
		NewDeleteCommand(),
		NewListCommand(),
	)
	return cmd
}
