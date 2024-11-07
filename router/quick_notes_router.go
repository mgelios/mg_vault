package router

import (
	"encoding/json"
	"log/slog"
	"mg_vault/auth"
	"mg_vault/model"
	"mg_vault/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func DefineQuickNotesProtectedRoutes(r chi.Router) {
	r.Get("/qnotes", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserClaimsFromContext(r)
		response := model.UserQuckNotesResponse{}
		response.Notes, _ = storage.GetAllQuickNotesForUser(user.Id)
		response.User = user
		if err := templates.ExecuteTemplate(w, "quick_notes.html", response); err != nil {
			slog.Error(err.Error())
		}
	})
	r.Get("/qnotes/create", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserClaimsFromContext(r)
		response := model.UserQuckNoteEditResponse{
			User: user,
			Note: model.QuickNote{},
		}
		if err := templates.ExecuteTemplate(w, "create_quick_note.html", response); err != nil {
			slog.Error(err.Error())
		}
	})
	r.Get("/qnotes/edit", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserClaimsFromContext(r)
		response := model.UserQuckNoteEditResponse{
			User: user,
		}
		response.Note, _ = storage.GetQuickNoteForUserWithId(user.Id, r.URL.Query().Get("qnote_id"))
		if err := templates.ExecuteTemplate(w, "edit_quick_note.html", response); err != nil {
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
	r.Put("/api/v1/qnotes", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserClaimsFromContext(r)
		var quickNote model.QuickNote
		err := json.NewDecoder(r.Body).Decode(&quickNote)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		quickNote.Author = user.Id

		storage.UpdateQuickNote(quickNote)
		w.Header().Add("HX-Redirect", "/qnotes")
	})
}
