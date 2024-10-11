package router

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"mg_vault/auth"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth/v5"
)

var templates *template.Template

func RunServer(templateFolder embed.FS) {
	templates = template.Must(template.ParseFS(templateFolder, "templates/*"))
	router := chi.NewRouter()
	router.Use(httprate.LimitAll(100, time.Second))
	definePublicEndpoints(router)
	defineSecuredEndpoints(router)
	slog.Debug("Applies handler to the router")

	err := http.ListenAndServe(":8000", router)
	if err != nil {
		slog.Error(err.Error())
	}
}

func definePublicEndpoints(router *chi.Mux) {
	fs := http.FileServer(http.Dir("static"))

	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	router.Route("/api/v1/user", func(router chi.Router) {
		router.Post("/login", auth.ProcessLoginRequest)
	})

	router.Route("/index", func(router chi.Router) {
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			if err := templates.ExecuteTemplate(w, "index.html", ""); err != nil {
				slog.Error(err.Error())
			}
		})
	})

	router.Route("/user", func(router chi.Router) {
		router.Get("/login", func(w http.ResponseWriter, r *http.Request) {
			if err := templates.ExecuteTemplate(w, "login.html", ""); err != nil {
				slog.Error(err.Error())
			}
		})
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
