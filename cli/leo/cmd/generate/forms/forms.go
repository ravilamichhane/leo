package forms

import (
	"github.com/spf13/cobra"
)

// modelCmd represents the model command
var FormsCmd = &cobra.Command{
	Use:   "form",
	Short: "Generate fetchers for frontend.",
	Long:  `Generate fetchers for frontend.`,
	Run: func(cmd *cobra.Command, args []string) {

		generate()

	},
}

func init() {

}
