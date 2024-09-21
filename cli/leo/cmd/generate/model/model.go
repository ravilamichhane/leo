/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package model

import (
	"fmt"

	"github.com/spf13/cobra"
)

// modelCmd represents the model command
var ModelCmd = &cobra.Command{
	Use:   "model",
	Short: "Generate model.",
	Long:  `Generate model.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("model called")
	},
}

func init() {

}
