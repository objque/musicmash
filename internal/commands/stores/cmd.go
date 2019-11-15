package stores

import "github.com/spf13/cobra"

func NewStoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store",
		Short: "Manage stores",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(
		NewListCommand(),
		NewCreateCommand(),
	)
	return cmd
}
