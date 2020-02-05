package stores

import "github.com/spf13/cobra"

func NewStoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "store",
		Aliases: []string{"stores"},
		Short:   "Manage stores",
		Args:    cobra.NoArgs,
	}
	cmd.AddCommand(
		NewListCommand(),
	)
	return cmd
}
