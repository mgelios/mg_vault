package router

import (
	"log/slog"
	"mg_vault/auth"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func DefineUserRoutes(router *chi.Mux) {
	router.Route("/api/v1/user", func(router chi.Router) {
		router.Post("/login", auth.ProcessLoginRequest)
	})

	router.Route("/user", func(router chi.Router) {
		router.Get("/login", func(w http.ResponseWriter, r *http.Request) {
			if err := templates.ExecuteTemplate(w, "login.html", ""); err != nil {
				slog.Error(err.Error())
			}
		})
	})
}
