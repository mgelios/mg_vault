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

func DefineProtectedLinkCategoryRoutes(r chi.Router) {
	r.Get("/links/category", OpenViewLinkCategoryPage)
	r.Get("/links/edit", OpenEditLinkCategoryPage)

	r.Post("/api/v1/links/category", CreateLinkCategory)

	r.Put("/api/v1/links/category", UpdateLinkCategory)

	r.Delete("/api/v1/links/category", DeleteLinkCategory)
}

func OpenViewLinkCategoryPage(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserClaimsFromContext(r)
	linkCategory := model.LinkCategory{}
	linkSubcategories := []model.LinkCategory{}
	if user.Id != "<nil>" {
		linkCategory, _ = storage.GetLinkCategoryById(r.URL.Query().Get("category_id"))
		linkSubcategories, _ = storage.GetLinkSubcategoriesByParentId(linkCategory.Id)
	}
	responseModel := model.MainPageResponse{
		User:              user,
		LinkCategory:      linkCategory,
		LinkSubcategories: linkSubcategories,
	}
	if err := templates.ExecuteTemplate(w, "view_link_category.html", responseModel); err != nil {
		slog.Error(err.Error())
	}
}

func OpenEditLinkCategoryPage(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserClaimsFromContext(r)
	response := model.LinkCategoryPageResponse{
		User: user,
	}
	linkCategory, _ := storage.GetLinkCategoryById(r.URL.Query().Get("category_id"))
	response.LinkCategory = linkCategory

	if err := templates.ExecuteTemplate(w, "edit_link_category.html", response); err != nil {
		slog.Error(err.Error())
	}
}

func CreateLinkCategory(w http.ResponseWriter, r *http.Request) {
	newLinkCategoryId, err := storage.CreateLinkCategoryWithParent(r.URL.Query().Get("parent_category_id"))

	if err != nil {
		slog.Error(err.Error())
	}

	w.Header().Add("HX-Redirect", "/links/edit?category_id="+newLinkCategoryId)
}

func UpdateLinkCategory(w http.ResponseWriter, r *http.Request) {
	var linkCategory model.LinkCategory
	err := json.NewDecoder(r.Body).Decode(&linkCategory)

	message, _ := json.Marshal(&linkCategory)
	slog.Info(string(message))

	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	storage.UpdateLinkCategory(linkCategory)
	w.Header().Add("HX-Redirect", "/")
}

func DeleteLinkCategory(w http.ResponseWriter, r *http.Request) {

}
