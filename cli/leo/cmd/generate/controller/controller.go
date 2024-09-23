package controller

/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/

import (
	"github.com/spf13/cobra"
)

// modelCmd represents the model command
var ControllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "Generate controller.",
	Long:  `Generate controller.`,
	Run: func(cmd *cobra.Command, args []string) {

		gnerateResource, err := cmd.Flags().GetBool("resource")
		if err != nil {
			panic(err)
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			panic(err)
		}
		generateController(name, gnerateResource)

	},
}

func init() {
	ControllerCmd.Flags().BoolP("resource", "r", false, "Generate resource endpoints")
	ControllerCmd.Flags().StringP("name", "n", "", "Name of Controller")
	ControllerCmd.MarkFlagsOneRequired("name")
}
