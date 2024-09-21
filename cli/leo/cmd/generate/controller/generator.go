package controller

import (
	_ "embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/google/uuid"
	"github.com/ravilmc/leo/helpers"
)

//go:embed templates/controller.txt
var controller string

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

type ControllerData struct {
	PackageName    string
	ControllerName string
	Name           string
	LowerName      string
}

func generateController(controllerName string) {
	data := ControllerData{
		Name:           helpers.Capitalise(controllerName),
		ControllerName: helpers.Capitalise(controllerName) + "Controller",
		PackageName:    strings.ToLower(controllerName) + "controller",
		LowerName:      strings.ToLower(controllerName),
	}
	controllerPath := "app/controllers/" + strings.ToLower(controllerName) + "controller/"

	mainfile := controllerPath + strings.ToLower(controllerName) + "controller.go"
	getfile := controllerPath + "Get" + data.Name + "ById.go"
	createfile := controllerPath + "Create" + data.Name + ".go"
	updatefile := controllerPath + "Update" + data.Name + ".go"
	deletefile := controllerPath + "Delete" + data.Name + ".go"
	getAllfile := controllerPath + "GetAll" + data.Name + ".go"

	err := os.MkdirAll(filepath.Dir(mainfile), os.ModePerm)
	if err != nil {
		panic(err)
	}

	WriteFile(mainfile, controller, data)
	WriteFile(getfile, getRequest, data)
	WriteFile(createfile, createRequest, data)
	WriteFile(updatefile, updateRequest, data)
	WriteFile(deletefile, deleteRequest, data)
	WriteFile(getAllfile, getAllRequest, data)

}

func WriteFile(path string, content string, data any) {

	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToCapitalize": func(s string) string {
			return strings.ToUpper(string(s[0])) + s[1:]
		},
	}
	// Write file

	template, err := template.New(uuid.New().String()).Funcs(
		funcMap).Parse(content)

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
