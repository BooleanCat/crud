package router

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/BooleanCat/crud/crud"
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
	return func(w http.ResponseWriter, req *http.Request) {
		if err := crud.Create(filepath.Join(store, chi.URLParam(req, "filename")), req.Body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func read(store string) func(_ http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := crud.ReadInto(filepath.Join(store, chi.URLParam(req, "filename")), w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func update(store string) func(_ http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		err := crud.Update(filepath.Join(store, chi.URLParam(req, "filename")), req.Body)
		if err == nil {
			return
		}

		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
	}
}

func delete(store string) func(_ http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		err := crud.Delete(filepath.Join(store, chi.URLParam(req, "filename")))
		if err == nil {
			return
		}

		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
	}
}
