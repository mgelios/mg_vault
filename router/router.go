package router

import (
	"embed"
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

	//TODO: refactor place to init root category
	storage.InitRootLinkCategoryIfAbsent()

	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(auth.TokenAuth))

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			user := auth.GetUserClaimsFromContext(r)
			linkCategory := model.LinkCategory{}
			linkSubcategories := []model.LinkCategory{}
			if user.Id != "<nil>" {
				linkCategory, _ = storage.GetRootLinkCategory()
				linkSubcategories, _ = storage.GetLinkSubcategoriesByParentId(linkCategory.Id)
			}
			responseModel := model.MainPageResponse{
				User:              user,
				LinkCategory:      linkCategory,
				LinkSubcategories: linkSubcategories,
			}
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

		DefineQuickNotesProtectedRoutes(r)
		DefineProtectedNoteRoutes(r)
		DefineProtectedUserRoutes(r)
		DefineProtectedLinkCategoryRoutes(r)
	})
}
