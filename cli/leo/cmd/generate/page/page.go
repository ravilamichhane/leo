package page

import (
	"log/slog"

	"github.com/spf13/cobra"
)

var PageCmd = &cobra.Command{
	Use:   "page",
	Short: "Generate page",
	Long:  "Generate page",
	Run: func(cmd *cobra.Command, args []string) {

		slog.Info("Page Generation started")
		generatePage()
		slog.Info("Page Generation completed")
	},
}
