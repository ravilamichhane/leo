package page

import (
	_ "embed"
	"log/slog"

	"github.com/ravilmc/leo/cli/leo/helpers"
)

//go:embed templates/page.txt
var page string

//go:embed templates/handler.txt
var handler string

func generatePage() {
	pageInfo := helpers.GetPathInfo("page.tsx")
	handlerInfo := helpers.GetPathInfo("handler.go")

	slog.Debug("Page Generation started")

	helpers.WriteFile(page, &pageInfo)
	helpers.WriteFile(handler, &handlerInfo)
}
