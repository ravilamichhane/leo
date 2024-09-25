package go_ssr

import (
	"log/slog"
	"os"

	"github.com/ravilmc/leo/reactssr/packages/cache"
	"github.com/ravilmc/leo/reactssr/packages/typeconverter"
	"github.com/ravilmc/leo/reactssr/packages/utils"
)

type Engine struct {
	Config                  *Config
	HotReload               *HotReload
	CacheManager            *cache.Manager
	CachedLayoutCSSFilePath string
}

// New creates a new gossr Engine instance
func New(config Config) (*Engine, error) {
	engine := &Engine{
		Config:       &config,
		CacheManager: cache.NewManager(),
	}
	if err := os.Setenv("APP_ENV", config.AppEnv); err != nil {
		slog.Error("Failed to set APP_ENV environment variable")
	}
	err := config.Validate()
	if err != nil {
		slog.Error("Failed to validate config")
		return nil, err
	}
	utils.CleanCacheDirectories()
	// If using a layout css file, build it and cache it
	if config.LayoutCSSFilePath != "" {
		if err = engine.BuildLayoutCSSFile(); err != nil {
			slog.Error("Failed to build layout css file")
			return nil, err
		}
	}

	// If running in production mode, return and don't start hot reload or type converter
	if os.Getenv("APP_ENV") == "production" {
		slog.Info("Running go-ssr in production mode")
		return engine, nil
	}
	slog.Info("Running go-ssr in development mode")
	slog.Debug("Starting type converter")
	// Start the type converter to convert Go types to Typescript types
	if err := typeconverter.Start(engine.Config.PropsStructsPath, engine.Config.GeneratedTypesPath); err != nil {
		slog.Error("Failed to init type converter")
		return nil, err
	}

	slog.Debug("Starting hot reload server")
	engine.HotReload = newHotReload(engine)
	engine.HotReload.Start()
	return engine, nil
}
