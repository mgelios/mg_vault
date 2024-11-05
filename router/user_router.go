package router

import (
	"log/slog"
	"mg_vault/auth"
	"mg_vault/model"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func DefinePublicUserRoutes(router *chi.Mux) {
	router.Post("/api/v1/user/login", auth.ProcessLoginRequest)
	router.Get("/user/login", func(w http.ResponseWriter, r *http.Request) {
		if err := templates.ExecuteTemplate(w, "login.html", ""); err != nil {
			slog.Error(err.Error())
		}
	})
}

func DefineProtectedUserRoutes(router chi.Router) {
	router.Post("/api/v1/user/logout", auth.ProcessLogoutRequest)
	router.Get("/user/profile", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserClaimsFromContext(r)
		response := model.ProfilePageResponse{
			User: user,
		}
		if err := templates.ExecuteTemplate(w, "profile.html", response); err != nil {
			slog.Error(err.Error())
		}
	})
}
