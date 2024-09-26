package helpers

import (
	"log/slog"
	"os"
	"strings"
	"text/template"

	"github.com/google/uuid"
)

type PageData struct {
	PackageName    string
	Path           string
	FilePath       string
	GenerationPath string
	Paginate       bool
	ComponentPath  string
	Method         string
	Name           string
	LowerName      string
}

func GetPathInfo(generationName string) PageData {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	route_path := strings.Split(cwd, "routes/")

	if len(route_path) < 2 {
		slog.Error("GENERATION", slog.Any("error", "You must be inside route directory"))
	}
	packagenames := strings.Split(route_path[1], "/")
	packagename := packagenames[len(packagenames)-1]

	return PageData{
		PackageName:    strings.ReplaceAll(packagename, "_", ""),
		Path:           "/" + strings.ReplaceAll(route_path[1], "_", ":"),
		FilePath:       "routes/" + route_path[1] + "/" + generationName,
		GenerationPath: cwd + "/" + generationName,
		Method:         "GET",
		ComponentPath:  "routes/" + route_path[1] + "/page.tsx",
	}
}

func WriteFile(content string, pageInfo *PageData) {

	template, err := template.New(uuid.New().String()).Parse(content)

	if err != nil {
		panic(err)
	}

	f, err := os.Create(pageInfo.GenerationPath)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	file, err := os.OpenFile(pageInfo.GenerationPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		panic(err)
	}

	defer file.Close()
	err = template.Execute(file, pageInfo)

	if err != nil {
		panic(err)
	}
}
