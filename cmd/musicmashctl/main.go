package main

import (
	"github.com/musicmash/musicmash/internal/command/commands"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "musicmash",
	}
	commands.AddCommands(rootCmd)
	_ = rootCmd.Execute()
}
