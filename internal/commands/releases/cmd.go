package releases

import "github.com/spf13/cobra"

func NewReleaseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "release",
		Aliases: []string{"releases", "rels", "rel"},
		Short:   "Manage releases",
		Args:    cobra.NoArgs,
	}
	cmd.AddCommand(
		NewByCommand(),
		NewForCommand(),
	)
	return cmd
}
