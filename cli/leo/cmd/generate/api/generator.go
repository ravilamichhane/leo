package api

import (
	_ "embed"
	"log/slog"
	"os"
	"strings"

	"github.com/ravilmc/leo/cli/leo/helpers"
)

//go:embed templates/get.txt
var getRequest string

//go:embed templates/create.txt
var createRequest string

//go:embed templates/update.txt
var updateRequest string

//go:embed templates/delete.txt
var deleteRequest string

//go:embed templates/getAll.txt
var getAllRequest string

type RouteData struct {
	PackageName string
	Name        string
	Path        string
	Method      string
	LowerName   string
}

func generateRoute(name string, method string, paginated bool) {

	apipagedata := helpers.GetPathInfo(name + ".go")

	apipagedata.Name = name
	apipagedata.LowerName = strings.ToLower(name[:1]) + name[1:]

	switch method {
	case "GET":
		switch paginated {
		case true:
			helpers.WriteFile(getAllRequest, &apipagedata)
		case false:
			helpers.WriteFile(getRequest, &apipagedata)
		}
	case "POST":
		helpers.WriteFile(createRequest, &apipagedata)
	case "PUT":
		helpers.WriteFile(updateRequest, &apipagedata)
	case "DELETE":
		helpers.WriteFile(deleteRequest, &apipagedata)
	default:
		slog.Error("GENERATION", slog.Any("error", "Method not supported"))
		os.Exit(1)
	}

}

// func WriteFile(path string, content string, data RouteData) {

// 	// Write file

// 	template, err := template.New(uuid.New().String()).Parse(content)

// 	if err != nil {
// 		slog.Error("ERROR PARSING TEMPLATE", slog.Any("error", err))
// 		os.Exit(1)
// 	}

// 	f, err := os.Create(path)

// 	if err != nil {

// 		panic(err)
// 	}

// 	defer f.Close()

// 	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

// 	if err != nil {
// 		panic(err)
// 	}

// 	defer file.Close()
// 	err = template.Execute(file, data)

// 	if err != nil {
// 		slog.Error("ERROR WRITING FILE", slog.Any("error", err))
// 	}
// }
