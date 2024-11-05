package router

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"mg_vault/auth"
	"mg_vault/model"
	"mg_vault/storage"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth/v5"
)

var templates *template.Template
var fileServer http.Handler
var router *chi.Mux

func InitServer(templateFolder embed.FS, staticContentFolder embed.FS) *chi.Mux {
	templates = template.Must(template.ParseFS(templateFolder, "templates/*.html", "templates/*/*.html"))
	var staticContent, _ = fs.Sub(staticContentFolder, "static")
	fileServer = http.FileServer(http.FS(staticContent))
	router = chi.NewRouter()
	router.Use(httprate.LimitAll(100, time.Second))
	definePublicEndpoints(router)
	defineSecuredEndpoints(router)
	return router
}

func redirectToHttps(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
}

func RunServer() {
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
		go http.ListenAndServe(":80", http.HandlerFunc(redirectToHttps))
		err := http.ListenAndServeTLS(":443", certfilePath, keyfilePath, router)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

func definePublicEndpoints(router *chi.Mux) {
	slog.Info("Starting init of public endpoints")
	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(auth.TokenAuth))

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			user := auth.GetUserClaimsFromContext(r)
			responseModel := model.MainPageResponse{User: user}
			if err := templates.ExecuteTemplate(w, "index.html", responseModel); err != nil {
				slog.Error(err.Error())
			}
		})
	})
	DefinePublicUserRoutes(router)
}

func defineSecuredEndpoints(router *chi.Mux) {
	slog.Info("Starting init of secured endpoints")
	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(auth.TokenAuth))
		r.Use(jwtauth.Authenticator(auth.TokenAuth))
		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			user := auth.GetUserClaimsFromContext(r)
			w.Write([]byte(fmt.Sprintf("Protected resource, hi %v", user.Username)))
		})
		r.Get("/qnotes/create", func(w http.ResponseWriter, r *http.Request) {
			user := auth.GetUserClaimsFromContext(r)
			response := model.UserQuckNotesResponse{
				User: user,
			}
			if err := templates.ExecuteTemplate(w, "edit_quick_note.html", response); err != nil {
				slog.Error(err.Error())
			}
		})
		r.Get("/qnotes", func(w http.ResponseWriter, r *http.Request) {
			user := auth.GetUserClaimsFromContext(r)
			response := model.UserQuckNotesResponse{}
			response.Notes, _ = storage.GetAllQuickNotesForUser(user.Id)
			response.User = user
			if err := templates.ExecuteTemplate(w, "quick_notes.html", response); err != nil {
				slog.Error(err.Error())
			}
		})
		r.Post("/api/v1/qnotes", func(w http.ResponseWriter, r *http.Request) {
			user := auth.GetUserClaimsFromContext(r)
			var quickNote model.QuickNote
			err := json.NewDecoder(r.Body).Decode(&quickNote)
			if err != nil {
				slog.Error(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			quickNote.Author = user.Id

			storage.CreateQuickNote(quickNote)
			w.Header().Add("HX-Redirect", "/qnotes")
		})
		DefineProtectedUserRoutes(r)
	})
}
