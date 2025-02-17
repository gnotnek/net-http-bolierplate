package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Run() {
	var command = &cobra.Command{
		Use:   "project name",
		Short: "A brief description of your application",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command.AddCommand(apiCmd())

	if err := command.Execute(); err != nil {
		log.Fatal().Msgf("failed to execute command, err: %v", err.Error())
	}
}
