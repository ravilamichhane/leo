package api

/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/

import (
	"bufio"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// modelCmd represents the model command
var ApiCmd = &cobra.Command{
	Use:   "api",
	Short: "Generate controller.",
	Long:  `Generate controller.`,
	Run: func(cmd *cobra.Command, args []string) {

		slog.Info("----------------------------------------------")
		slog.Info("Enter the name of the route you want to create.")
		slog.Info("----------------------------------------------")
		reader := bufio.NewReader(os.Stdin)
		name, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		name = strings.ReplaceAll(name, "\n", "")

		slog.Info("----------------------------------------------")
		slog.Info("Enter the methodG of the route")
		slog.Info("----------------------------------------------")
		method, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		method = strings.ReplaceAll(method, "\n", "")

		shouldPaginate := false

		switch strings.ToUpper(method) {
		case "GET":
			slog.Info("----------------------------------------------")
			slog.Info("Is it a paginated route? (y/n) default: n")
			slog.Info("----------------------------------------------")
			paginate, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			paginate = strings.ReplaceAll(paginate, "\n", "")

			switch strings.ToUpper(paginate) {
			case "Y":
				shouldPaginate = true
			}
		}

		generateRoute(name, strings.ToUpper(method), shouldPaginate)
		// generateController(name, gnerateResource)

	},
}

func init() {
	ApiCmd.Flags().StringP("method", "m", "GET", "Method of Controller")
	ApiCmd.Flags().StringP("name", "n", "", "Name of Controller")
	ApiCmd.Flags().BoolP("paginated", "p", false, "Paginated Controller")
}
