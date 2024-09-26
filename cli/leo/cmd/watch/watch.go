package watch

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/ravilmc/leo/tygo"
	"github.com/spf13/cobra"
)

var WatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch changes",
	Long:  "Watch changes",
	Run: func(cmd *cobra.Command, args []string) {

		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		folder := cwd + "/app/routes"

		watcher, err := fsnotify.NewWatcher()

		if err != nil {
			slog.Error("START", slog.Any("error", err))
		}
		defer watcher.Close()

		done := make(chan bool)

		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if strings.Contains(event.Name, ".") {
						if strings.Contains(event.Name, ".go") {
							if strings.Contains(event.Name, "/routes/api/") {
								dir := path.Dir(event.Name)

								splitted := strings.Split(dir, "/routes/api")
								slog.Debug("START", slog.String("dir", dir))
								slog.Debug("START", slog.Any("splitted", splitted))

								path := "github.com/ravilmc/leoapp/app/routes/api" + splitted[1]

								slog.Debug("START", slog.String("path", path))

								slog.Info("generating api types and fetcher")
								gen := tygo.New(&tygo.Config{
									Packages: []*tygo.PackageConfig{
										{
											Path:       path,
											OutputPath: dir + "/api.ts",
											Frontmatter: `
							import { Fetcher, handleResponseError } from "./fetcher";
							import {useForm} from "react-hook-form";
							import { SafeParse } from "./safeparse";
							import { generateFormData } from "./formdata";
							import { useMutation } from "@tanstack/react-query";
							import { Form } from "@/components/ui/form";
							import { Button } from "@/components/ui/button";
							import { DateInput } from "@/components/forms/date-input";
							import { HTMLInput } from "@/components/forms/html-input";
							import { ImageInput } from "@/components/forms/image-upload";
							import { SelectInput } from "@/components/forms/select-input";
							import { SwitchInput } from "@/components/forms/switch-input";
							import { TextInput } from "@/components/forms/text-input";
							
											`,
											TypeMappings: map[string]string{
												"time.Time": "string",
											},
										},
									},
								})

								err = gen.Generate()

								if err == nil {

									cmd := exec.Command("./node_modules/.bin/prettier", dir+"/api.ts", "--write")
									stdout, err := cmd.Output()

									if err != nil {
										fmt.Println(err.Error())
										return
									}

									fmt.Println(string(stdout))
								} else {
									slog.Error("generating api types", slog.Any("error", err))
								}

								// Print the output

							} else {

								slog.Info("generating page types")

							}
						}
					}

					// Execute the task on every event
					slog.Info("FileChanged", slog.String("path", event.Name))
					if strings.Contains(event.Name, ".go") {

					}

				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.Println("error:", err)
				}
			}
		}()

		err = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error walking directory: %w", err)
			}
			// Only add directories
			if info.IsDir() {
				fmt.Println("Watching directory:", path)
				err = watcher.Add(path)
				if err != nil {
					return fmt.Errorf("error adding directory to watcher: %w", err)
				}
			}
			return nil
		})

		slog.Error("WATCH", slog.Any("error", err))

		// Block main goroutine until done is closed
		<-done

		os.Exit(0)

	},
}
