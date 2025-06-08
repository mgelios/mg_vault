package model

type Link struct {
	Id             string `json:"id,omitempty" bson:"_id,omitempty"`
	Url            string `json:"url" bson:"url"`
	Name           string `json:"name" bson:"name"`
	BadgeText      string `json:"badge_text" bson:"badge_text"`
	BadgeColor     string `json:"badge_color" bson:"badge_color"`
	BadgeTextColor string `json:"badge_text_color" bson:"badge_text_color"`
}

type LinkGroup struct {
	Id    string `json:"id,omitempty" bson:"_id,omitempty"`
	Links []Link `json:"links" bson:"links"`
	Name  string `json:"name" bson:"name"`
}

type LinkCategory struct {
	Id         string      `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string      `json:"name" bson:"name"`
	Links      []Link      `json:"links" bson:"links"`
	LinkGroups []LinkGroup `json:"link_groups" bson:"link_groups"`
}
