package associations

import "github.com/spf13/cobra"

func NewAssociationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "associations",
		Aliases: []string{"assoc", "ass", "as"},
		Short:   "Manage artist associations",
		Args:    cobra.NoArgs,
	}
	cmd.AddCommand(
		NewCreateCommand(),
		NewForCommand(),
		NewListCommand(),
	)
	return cmd
}
