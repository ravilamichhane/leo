/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package fetchers

import (
	"github.com/spf13/cobra"
)

// modelCmd represents the model command
var FetcherCmd = &cobra.Command{
	Use:   "fetchers",
	Short: "Generate fetchers for frontend.",
	Long:  `Generate fetchers for frontend.`,
	Run: func(cmd *cobra.Command, args []string) {
		controller, err := cmd.Flags().GetString("controller")
		if err != nil {
			panic(err)
		}

		generate(controller)

	},
}

func init() {
	FetcherCmd.Flags().StringP("controller", "c", "", "Name of controller")
	FetcherCmd.MarkFlagsOneRequired("controller")

}
