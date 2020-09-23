package command

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kreuzwerker/asc/database"
	"github.com/kreuzwerker/asc/transfer"
	"github.com/pkg/errors"
	progressbar "github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{

	Use:   "import [file]",
	Short: "Import pricing information from file",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		matches, err := filepath.Glob(args[0])

		if err != nil {
			return err
		}

		database, err := database.New("asc.db")

		if err != nil {
			return err
		}

		defer database.Close()

		for _, match := range matches {

			fi, err := os.Stat(match)

			if err != nil {
				return errors.Wrapf(err, "failed to get file stats for %q", match)
			}

			fh, err := os.Open(match)

			if err != nil {
				return errors.Wrapf(err, "failed to open file %q", match)
			}

			bar := progressbar.DefaultBytes(
				fi.Size(),
				fmt.Sprintf("importing %q", match),
			)

			if err := transfer.Import(io.TeeReader(fh, bar), func(header *transfer.Header, sku *transfer.SKU) error {
				return database.Add(header, sku)
			}); err != nil {
				return errors.Wrapf(err, "failed to import %q", match)
			}

			if err := database.Commit(); err != nil {
				return errors.Wrapf(err, "failed to commit file %q to batch", match)
			}

			if err := fh.Close(); err != nil {
				return errors.Wrapf(err, "failed to close file %q", match)
			}

		}

		return nil

	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
