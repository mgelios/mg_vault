package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Link struct {
	Url            string `json:"url" bson:"url"`
	Name           string `json:"name" bson:"name"`
	BadgeText      string `json:"badge_text" bson:"badge_text"`
	BadgeColor     string `json:"badge_color" bson:"badge_color"`
	BadgeTextColor string `json:"badge_text_color" bson:"badge_text_color"`
}

type LinkGroup struct {
	Name  string `json:"name" bson:"name"`
	Links []Link `json:"links" bson:"links"`
}

type LinkCategory struct {
	Id         string             `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	Links      []Link             `json:"links" bson:"links"`
	LinkGroups []LinkGroup        `json:"link_groups" bson:"link_groups"`
	Parent     primitive.ObjectID `json:"parent" bson:"parent"`
}

type LinkCategoryUpdate struct {
	Name       string      `json:"name" bson:"name"`
	Links      []Link      `json:"links" bson:"links"`
	LinkGroups []LinkGroup `json:"link_groups" bson:"link_groups"`
}

type LinkCategoryPageResponse struct {
	User              UserClaims     `json:"user"`
	LinkCategory      LinkCategory   `json:"link_category"`
	LinkSubcategories []LinkCategory `json:"link_subcategories"`
}
