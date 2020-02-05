package notifysettings

import "github.com/spf13/cobra"

func NewNotificationSettingsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "notification-settings",
		Aliases: []string{"ns"},
		Short:   "Manage notifications",
		Args:    cobra.NoArgs,
	}
	cmd.AddCommand(
		NewCreateCommand(),
		NewUpdateCommand(),
		NewListCommand(),
	)
	return cmd
}
