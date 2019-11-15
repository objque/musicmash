package releases

import "github.com/spf13/cobra"

func NewReleaseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Manage releases",
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(
		NewListCommand(),
	)
	return cmd
}
