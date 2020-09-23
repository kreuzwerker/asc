package command

import (
	"github.com/kreuzwerker/asc/console"
	"github.com/kreuzwerker/asc/database"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           app,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {

		database, err := database.Open("asc.db")

		if err != nil {
			return err
		}

		defer database.Close()

		console := console.New(database)

		console.Run()

		return nil

	},
}
