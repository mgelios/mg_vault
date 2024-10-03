package router

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RunServer() {
	router := chi.NewRouter()
	Handler(router)
	slog.Debug("Applies handler to the router")

	err := http.ListenAndServe(":8000", router)

	if err != nil {
		panic(err)
	}
}

func Handler(router *chi.Mux) {
	router.Route("/api/v1", func(router chi.Router) {
		router.Get("/get", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Some response"))
		})
	})
}
