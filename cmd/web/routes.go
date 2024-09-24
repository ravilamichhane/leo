package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ravilmc/leo/node"
)

func (app *application) routes() http.Handler {

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(60 * time.Second))

	mux.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	mux.Get("/*", func(w http.ResponseWriter, r *http.Request) {

		initialProps := InitialProps{
			Name:          "GoでReactをSSRする",
			InitialNumber: 100,
		}
		jsonProps, err := json.Marshal(initialProps)

		vm := node.New(nil)
		vm.Run(app.BackendBundle)

		val := vm.Run(fmt.Sprintf(`renderApp("http://localhost:4000%s")`, r.URL.Path))
		renderedHTML := val.String()

		log.Println(renderedHTML)
		log.Println(val.Error())

		tmpl, err := template.New("webpage").Parse(htmlTemplate)
		if err != nil {
			log.Fatal("Error parsing template:", err)
		}

		w.Header().Set("Content-Type", "text/html")
		data := PageData{
			RenderedContent: template.HTML(renderedHTML),
			InitialProps:    template.JS(jsonProps),
			JS:              template.JS(app.ClientBundle),
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	return mux
}
