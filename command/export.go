package command

import (
	"github.com/spf13/cobra"
	"github.com/kreuzwerker/asc/transfer"
)

var exportCmd = &cobra.Command{

	Use:   "export",
	Short: "Exports pricing information from AWS",
	RunE: func(cmd *cobra.Command, args []string) error {
		return transfer.Export("export")
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
