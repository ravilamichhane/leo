/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
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
		controller, err := cmd.Flags().GetString("controller")
		if err != nil {
			panic(err)
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			panic(err)
		}
		generate(controller, name)

	},
}

func init() {
	FormsCmd.Flags().StringP("controller", "c", "", "Name of controller")
	FormsCmd.Flags().StringP("name", "n", "", "Name of Method")
	FormsCmd.MarkFlagsOneRequired("controller")
	FormsCmd.MarkFlagsOneRequired("name")

}
