package main

import (
	"os"

	"github.com/musicmash/musicmash/internal/commands"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/log"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "musicmashctl",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logLevel := "info"
			if isDebugMode, _ := cmd.Flags().GetBool("debug"); isDebugMode {
				logLevel = "debug"
			}
			logPath, _ := cmd.Flags().GetString("log-path")
			log.SetLogFormatter(&log.DefaultFormatter)
			log.ConfigureStdLogger(logLevel, logPath)

			config.Config = config.New()
			configPath, _ := cmd.Flags().GetString("config")
			if configPath != "" {
				if err := config.Config.LoadFromFile(configPath); err != nil {
					exitWithError(err)
				}
			}

			// prioritise cli flags
			if ip, _ := cmd.Flags().GetString("http-ip"); ip != "" {
				config.Config.HTTP.IP = ip
			}
			if port, _ := cmd.Flags().GetInt("http-port"); port != 0 {
				config.Config.HTTP.Port = port
			}
		},
	}
	commands.AddCommands(rootCmd)

	rootCmd.PersistentFlags().String("config", "", "Path to musicmash.yaml")
	rootCmd.PersistentFlags().String("http-ip", "", "API ip address")
	rootCmd.PersistentFlags().Int("http-port", 0, "API port")
	rootCmd.PersistentFlags().Bool("debug", false, "Set log-level as debug")
	rootCmd.PersistentFlags().String("log-path", "", "Path to file for output logs")
	_ = rootCmd.Execute()
}

func exitWithError(err error) {
	log.Error(err)
	os.Exit(2)
}
