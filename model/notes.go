package model

type QuickNote struct {
	Id      string `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string `json:"name" bson:"name"`
	Content string `json:"content" bson:"content"`
	Author  string `json:"author" bson:"author"`
}

type QuickNoteUpdate struct {
	Name    string `json:"name" bson:"name"`
	Content string `json:"content" bson:"content"`
	Author  string `json:"author" bson:"author"`
}

type UserQuckNotesResponse struct {
	Notes []QuickNote `json:"qnotes"`
	User  UserClaims  `json:"user"`
}

type UserQuckNoteEditResponse struct {
	Note QuickNote  `json:"qnote"`
	User UserClaims `json:"user"`
}

type Note struct {
	Id      string   `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string   `json:"name" bson:"name"`
	Content string   `json:"content" bson:"content"`
	Path    []string `json:"path" bson:"path"`
	Tags    []string `json:"tags" bson:"tags"`
	Author  string   `json:"author" bson:"author"`
}

type NoteUpdate struct {
	Name    string   `json:"name" bson:"name"`
	Content string   `json:"content" bson:"content"`
	Path    []string `json:"path" bson:"path"`
	Tags    []string `json:"tags" bson:"tags"`
	Author  string   `json:"author" bson:"author"`
}

type UserNotesResponse struct {
	User  UserClaims `json:"user"`
	Notes []Note     `json:"notes"`
	Tree  NotesTree  `json:tree`
}
type UserNoteResponse struct {
	User UserClaims `json:"user"`
	Note Note       `json:"note"`
}

type NotesTree struct {
	Id     string        `json:"id,omitempty" bson:"_id,omitempty"`
	Author string        `json:"author" bson:"author"`
	Root   NotesTreeNode `json:"root" bson:"root"`
}

type NotesTreeUpdate struct {
	Author string        `json:"author" bson:"author"`
	Root   NotesTreeNode `json:"root" bson:"root"`
}

type NotesTreeNode struct {
	ChildNodes  map[string]*NotesTreeNode `json:"child_nodes" bson:"child_nodes"`
	Entries     int                       `json:"entries" bson:"entries"`
	Breadcrumbs []string                  `json:"breadcrumbs" bson:"breadcrumbs"`
}
