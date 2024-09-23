package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/ravilmc/leo/goreact"
)

const port = ":4000"

type application struct {
	engine *goreact.Engine
}

var APP_ENV string

func main() {

	app := &application{}

	engine, err := goreact.New(goreact.Config{
		AppEnv:             APP_ENV,
		AssetRoute:         "/assets",
		FrontendDir:        "../frontend",
		GeneratedTypesPath: "../frontend/generated.d.ts",
		TailwindConfigPath: "../frontend/tailwind.config.js",
		PropsStructsPath:   "./models/props.go",
		LayoutFilePath:     "Layout.tsx",
		LayoutCSSFilePath:  "Layout.css",
	})

	if err != nil {
		panic(err)
	}

	app.engine = engine
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	slog.SetDefault(logger)

	srv := &http.Server{
		Addr:              port,
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	slog.Debug("START", slog.String("port", port))
	slog.Info("START", slog.String("message", fmt.Sprintf("Server is running on port %s", port)))
	slog.Info("START", slog.String("url", fmt.Sprintf("http://localhost%s", port)))

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
