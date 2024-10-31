package model

type QuickNote struct {
	Id      string `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string `json:"name" bson:"name"`
	Content string `json:"content" bson:"content"`
	Author  string `json:"author" bson:"author"`
}

type UserQuckNotesResponse struct {
	Notes []QuickNote `json:"qnotes"`
}

type Note struct {
	Id      string   `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string   `json:"name" bson:"name"`
	Content string   `json:"content" bson:"content"`
	Path    []string `json:"path" bson:"path"`
	Tags    []string `json:"tags" bson:"tags"`
	Author  string   `json:"author" bson:"author"`
}
