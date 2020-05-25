package artists

import "github.com/spf13/cobra"

func NewArtistCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "artist",
		Aliases: []string{"artists", "arts", "art"},
		Short:   "Manage artists",
		Args:    cobra.NoArgs,
	}
	cmd.AddCommand(
		NewCreateCommand(),
		NewShowCommand(),
	)
	return cmd
}
