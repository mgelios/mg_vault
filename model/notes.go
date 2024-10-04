package model

type QuickNote struct {
	Id      string `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string `json:"name" bson:"name"`
	Content string `json:"content" bson:"content"`
}

type Note struct {
	Id      string   `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string   `json:"name" bson:"name"`
	Content string   `json:"content" bson:"content"`
	Path    []string `json:"path" bson:"path"`
	Tags    []string `json:"tags" bson:"tags"`
}
