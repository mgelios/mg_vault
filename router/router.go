package router

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"mg_vault/auth"
	"mg_vault/model"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
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

		router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			var loginRequestBody model.LoginRequest
			err := json.NewDecoder(r.Body).Decode(&loginRequestBody)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(loginRequestBody.Password), bcrypt.DefaultCost)

			// if err != nil {
			// 	http.Error(w, err.Error(), http.StatusBadRequest)
			// }

			// loginRequestBody.Password = string(hashedPassword)
			token := auth.ProcessLoginRequest(loginRequestBody)
			var cookie *http.Cookie = &http.Cookie{
				Name:  "Authorization",
				Value: "Bearer " + token,
				Path:  "/",
			}
			http.SetCookie(w, cookie)
		})
	})

	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(auth.TokenAuth))
		r.Use(jwtauth.Authenticator(auth.TokenAuth))

		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("Protected resource, hi %v", claims["id"])))
		})
	})
}
