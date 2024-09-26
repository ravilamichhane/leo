package routetree

import (
	"bufio"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type HandlerData struct {
	Name   string
	Import string
}

//go:embed templates/routetree.txt
var routetree string

var RouteTreeCmd = &cobra.Command{
	Use:   "routetree",
	Short: "Generate route tree.",
	Long:  `Generate route tree.`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		file, err := os.Open("go.mod")
		if err != nil {
			fmt.Println("No mod file found:", err)
			return
		}
		defer file.Close()

		// Create a new scanner
		scanner := bufio.NewScanner(file)
		basepackageName := ""
		if scanner.Scan() {
			line := scanner.Text()
			// slog.Debug("line", slog.String("line", line))
			// Check if the line contains the package name
			if strings.Contains(line, "module") {
				// Split the line by spaces
				parts := strings.Split(line, " ")
				// Get the package name
				basepackageName = parts[1]
			}

		}

		if err := scanner.Err(); err != nil {
			slog.Error("Error reading file:", slog.Any("err", err))
			return
		}

		slog.Debug("name:" + basepackageName)

		if basepackageName == "" {
			slog.Error("Invalid package")
			return
		}
		handlers := make([]HandlerData, 0)

		err = filepath.Walk(cwd+"/app/routes", func(path string, info os.FileInfo, err error) error {

			if strings.Contains(path, "/app/routes/api") {
				return nil
			}

			if !info.IsDir() {
				name := info.Name()
				if name == "handler.go" {
					splitted := strings.Split(path, "/app/routes")
					// slog.Debug("splitted", slog.Any("splitted", splitted))
					if len(splitted) < 2 {
						slog.Error("Invalid path")
						os.Exit(1)
					}

					path = strings.ReplaceAll(splitted[1], "/handler.go", "")
					handlers = append(handlers, HandlerData{
						Name:   strings.TrimSuffix("routes_"+strings.ReplaceAll(strings.TrimPrefix(path, "/"), "/", "_"), "_"),
						Import: basepackageName + "/app/routes" + path,
					})

				}

			}

			if err != nil {
				slog.Error("error adding directory to watcher", slog.Any("err", err))
				return fmt.Errorf("error adding directory to watcher: %w", err)
			}

			template, err := template.New(uuid.New().String()).Parse(routetree)

			if err != nil {
				panic(err)
			}

			f, err := os.Create(cwd + "/routes.go")

			if err != nil {
				panic(err)
			}

			defer f.Close()

			file, err := os.OpenFile(cwd+"/routes.go", os.O_APPEND|os.O_WRONLY, os.ModeAppend)

			if err != nil {
				panic(err)
			}

			defer file.Close()
			err = template.Execute(file, handlers)

			if err != nil {
				panic(err)
			}

			log.Println(handlers)
			return nil
		})
		if err != nil {
			return
		}

	},
}
