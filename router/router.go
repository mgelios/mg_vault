package router

import (
	"fmt"
	"log/slog"
	"mg_vault/auth"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func RunServer() {
	router := chi.NewRouter()
	definePublicEndpoints(router)
	defineSecuredEndpoints(router)
	slog.Debug("Applies handler to the router")

	err := http.ListenAndServe(":8000", router)
	if err != nil {
		slog.Error(err.Error())
	}
}

func definePublicEndpoints(router *chi.Mux) {
	router.Route("/api/v1/user", func(router chi.Router) {
		router.Post("/login", auth.ProcessLoginRequest)
	})

	router.Route("/index", func(router chi.Router) {

	})
}

func defineSecuredEndpoints(router *chi.Mux) {
	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(auth.TokenAuth))
		r.Use(jwtauth.Authenticator(auth.TokenAuth))
		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("Protected resource, hi %v", claims["id"])))
		})
	})
}
