package page

import (
	_ "embed"
	"html/template"
	"log/slog"
	"os"
	"strings"

	"github.com/google/uuid"
)

//go:embed templates/page.txt
var page string

//go:embed templates/handler.txt
var handler string

type PageData struct {
	PackageName string
	Path        string
	FilePath    string
}

func generatePage() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	path := strings.Split(cwd, "routes/")

	if len(path) < 2 {
		slog.Error("GENERATION", slog.Any("error", "You must be inside route directory"))
	}

	packagenames := strings.Split(path[1], "/")
	packagename := packagenames[len(packagenames)-1]

	data := PageData{
		Path:        "/" + path[1],
		PackageName: packagename,
		FilePath:    "routes/" + path[1] + "/page.tsx",
	}

	slog.Debug("Start", "dir", path[1])
	slog.Debug("Page Generation started")

	WriteFile(cwd+"/page.tsx", page, data)
	WriteFile(cwd+"/handler.go", handler, data)
}

func WriteFile(path string, content string, data any) {

	template, err := template.New(uuid.New().String()).Parse(content)

	if err != nil {
		panic(err)
	}

	f, err := os.Create(path)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		panic(err)
	}

	defer file.Close()
	err = template.Execute(file, data)

	if err != nil {
		panic(err)
	}
}
