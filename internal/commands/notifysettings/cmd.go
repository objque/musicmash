package notifysettings

import "github.com/spf13/cobra"

func NewNotificationSettingsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "notification-settings",
		Short:   "Manage notifications",
		Aliases: []string{"ns"},
		Args:    cobra.NoArgs,
	}
	cmd.AddCommand(
		NewCreateCommand(),
		NewUpdateCommand(),
		NewListCommand(),
	)
	return cmd
}
