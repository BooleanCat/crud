package router

import (
	"fmt"
	"io"
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

	router.Route("/Update", func(r chi.Router) {
		r.Post("/{filename}", update(store))
	})

	router.Route("/Delete", func(r chi.Router) {
		r.Delete("/{filename}", delete(store))
	})

	return router
}

func create(store string) func(_ http.ResponseWriter, req *http.Request) {
	return func(_ http.ResponseWriter, req *http.Request) {
		writeFile(store, req)
	}
}

func read(store string) func(_ http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		content, err := ioutil.ReadFile(filepath.Join(store, chi.URLParam(req, "filename")))
		if err != nil {
			panic(err)
		}

		_, err = io.WriteString(w, fmt.Sprintf("Content: %s", content))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func update(store string) func(_ http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		_, err := os.Stat(filepath.Join(store, chi.URLParam(req, "filename")))
		isFileExists(err, w)
		writeFile(store, req)
	}
}

func delete(store string) func(_ http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		err := os.Remove(filepath.Join(store, chi.URLParam(req, "filename")))
		isFileExists(err, w)
	}
}

func isFileExists(err error, w http.ResponseWriter) {
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func writeFile(store string, req *http.Request) {
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	if err = ioutil.WriteFile(filepath.Join(store, chi.URLParam(req, "filename")), content, os.ModePerm); err != nil {
		panic(err)
	}
}
