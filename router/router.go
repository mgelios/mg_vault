package router

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"mg_vault/auth"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth/v5"
)

var templates *template.Template
var fileServer http.Handler

func RunServer(templateFolder embed.FS, staticContentFolder embed.FS) {
	templates = template.Must(template.ParseFS(templateFolder, "templates/*"))
	var staticContent, _ = fs.Sub(staticContentFolder, "static")
	fileServer = http.FileServer(http.FS(staticContent))

	router := chi.NewRouter()
	router.Use(httprate.LimitAll(100, time.Second))
	definePublicEndpoints(router)
	defineSecuredEndpoints(router)
	slog.Debug("Applies handler to the router")

	env := os.Getenv("MG_ENV")
	if env == "" || env == "dev" {
		slog.Info("Setting up dev env")
		err := http.ListenAndServe(":80", router)
		if err != nil {
			slog.Error(err.Error())
		}
	} else if env == "prod" {
		slog.Info("Setting up prod env")
		certfilePath := os.Getenv("MG_VAULT_CERT_PATH")
		keyfilePath := os.Getenv("MG_VAULT_KEY_PATH")
		err := http.ListenAndServeTLS(":443", certfilePath, keyfilePath, router)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

func definePublicEndpoints(router *chi.Mux) {
	slog.Info("Init of fileserver finished")
	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

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
