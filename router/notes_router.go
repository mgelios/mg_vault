package router

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"mg_vault/auth"
	"mg_vault/model"
	"mg_vault/storage"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func DefineProtectedNoteRoutes(r chi.Router) {
	r.Get("/notes", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserClaimsFromContext(r)
		path := strings.Split(r.URL.Query().Get("path"), ",")
		response := model.UserNotesResponse{}
		response.Notes, _ = storage.GetAllNotesForUserInPath(user.Id, path)
		response.Tree, _ = storage.GetNotesTreeForUser(user.Id)
		response.User = user
		if err := templates.ExecuteTemplate(w, "notes.html", response); err != nil {
			slog.Error(err.Error())
		}
	})
	r.Get("/notes/view", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserClaimsFromContext(r)
		response := model.UserNoteResponse{
			User: user,
		}
		response.Note, _ = storage.GetNoteById(r.URL.Query().Get("note_id"))
		if err := templates.ExecuteTemplate(w, "view_note.html", response); err != nil {
			slog.Error(err.Error())
		}
	})
	r.Get("/notes/create", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserClaimsFromContext(r)
		response := model.UserNoteResponse{
			User: user,
			Note: model.Note{},
		}
		if err := templates.ExecuteTemplate(w, "create_note.html", response); err != nil {
			slog.Error(err.Error())
		}
	})
	r.Get("/notes/edit", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserClaimsFromContext(r)
		response := model.UserNoteResponse{
			User: user,
		}
		response.Note, _ = storage.GetNoteById(r.URL.Query().Get("note_id"))
		if err := templates.ExecuteTemplate(w, "edit_note.html", response); err != nil {
			slog.Error(err.Error())
		}
	})
	r.Get("/api/v1/notes/export", func(w http.ResponseWriter, r *http.Request) {
		buffer := &bytes.Buffer{}
		user := auth.GetUserClaimsFromContext(r)
		response := model.UserNotesResponse{}
		response.Notes, _ = storage.GetAllNotesForUser(user.Id)
		response.Tree, _ = storage.GetNotesTreeForUser(user.Id)
		response.User = user
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err := json.NewEncoder(buffer).Encode(response)
		if err != nil {
			slog.Error("Error during encoding notes to buffer")
		}
		w.Write(buffer.Bytes())
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
	r.Delete("/api/v1/notes", func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserClaimsFromContext(r)
		storage.DeleteNoteByUserAndId(r.URL.Query().Get("note_id"), user.Id)
		w.Header().Add("HX-Redirect", "/notes")
	})
}
