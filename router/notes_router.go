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

func DefineProtectedNoteRoutes(r chi.Router) {
	r.Get("/notes", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserClaimsFromContext(r)
		response := model.UserNotesResponse{}
		response.Notes, _ = storage.GetAllNotesForUser(user.Id)
		response.User = user
		if err := templates.ExecuteTemplate(w, "notes.html", response); err != nil {
			slog.Error(err.Error())
		}
	})
	r.Get("/notes/create", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserClaimsFromContext(r)
		response := model.UserNoteEditResponse{
			User: user,
			Note: model.Note{},
		}
		if err := templates.ExecuteTemplate(w, "edit_note.html", response); err != nil {
			slog.Error(err.Error())
		}
	})
	r.Get("/notes/edit", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserClaimsFromContext(r)
		response := model.UserNoteEditResponse{
			User: user,
		}
		response.Note, _ = storage.GetNoteForUserWithId(user.Id, r.URL.Query().Get("note_id"))
		if err := templates.ExecuteTemplate(w, "edit_note.html", response); err != nil {
			slog.Error(err.Error())
		}
	})
	r.Post("/api/v1/notes", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserClaimsFromContext(r)
		var note model.Note
		err := json.NewDecoder(r.Body).Decode(&note)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		note.Author = user.Id
		storage.CreateNote(note)
		w.Header().Add("HX-Redirect", "/notes")
	})
	r.Put("/api/v1/notes", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserClaimsFromContext(r)
		var note model.Note
		err := json.NewDecoder(r.Body).Decode(&note)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		note.Author = user.Id
		storage.UpdateNote(note)
		w.Header().Add("HX-Redirect", "/notes")
	})
}
