package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Router creates a new router and returns it.
func Router(store string) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Heartbeat("/ping"))
	router.Route("/Create", func(r chi.Router) {
		r.Get("/{filename}", create(store))
	})

	router.Route("/Read", func(r chi.Router) {
		r.Get("/{filename}", read(store))
	})

	return router
}

func create(store string) func(_ http.ResponseWriter, req *http.Request) {
	return func(_ http.ResponseWriter, req *http.Request) {
		content, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}
		if err = ioutil.WriteFile(filepath.Join(store, chi.URLParam(req, "filename")), content, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func read(store string) func(_ http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		content, err := ioutil.ReadFile(filepath.Join(store, chi.URLParam(req, "filename")))
		if err != nil {
			panic(err)
		}
		w.Write([]byte(fmt.Sprintf("Content: %s", content)))
	}
}
