package main

import (
	"net/http"

	"github.com/ravilmc/leo/goreact"
)

type Props struct {
	Name string `json:"name"`
}

func (app *application) showHome(w http.ResponseWriter, r *http.Request) {
	app.renderReact(w, goreact.RenderConfig{
		File:  "Pages/App.tsx",
		Title: "Echo example app",
		MetaTags: map[string]string{
			"og:title":    "Echo example app",
			"description": "Hello world!",
		},
		Props: &Props{
			Name: "ravi",
		},
	})
}
