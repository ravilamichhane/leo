package generate

import (
	"log"

	"github.com/ravilmc/leo/cli/leo/cmd/generate/controller"
	"github.com/ravilmc/leo/cli/leo/cmd/generate/model"
	"github.com/spf13/cobra"
)

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code",
	Long:  `Generate code`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		switch args[0] {
		case "model":
			if len(args) < 2 {
				cmd.Help()
				return
			}
			modelName := args[1]
			generateModel(modelName)
		default:
			cmd.Help()
		}
	},
}

func init() {
	GenerateCmd.AddCommand(model.ModelCmd)
	GenerateCmd.AddCommand(controller.ControllerCmd)
}

func generateModel(modelName string) {

	log.Println("Generating model", modelName)

}
