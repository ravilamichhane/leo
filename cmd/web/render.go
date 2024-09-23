package main

import (
	"net/http"

	"github.com/ravilmc/leo/goreact"
)

func (app *application) renderReact(w http.ResponseWriter, config goreact.RenderConfig) {

	response := app.engine.RenderRoute(config)

	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/html")
	w.Write(response)

}
