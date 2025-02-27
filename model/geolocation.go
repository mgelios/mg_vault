package model

type Geolocation struct {
	ID          string  `json:"id,omitempty" bson:"_id,omitempty"`
	Latitude    float64 `json:"latitude" bson:"latitude"`
	Longitude   float64 `json:"longitude" bson:"longitude"`
	Description string  `json:"description" bson:"description"`
	Author      string  `json:"author" bson:"author"`
}

type Route struct {
	ID          int           `json:"id" bson:"id"`
	Name        string        `json:"name" bson:"name"`
	Description string        `json:"description" bson:"description"`
	Waypoints   []Geolocation `json:"waypoints" bson:"waypoints"`
	Author      string        `json:"author" bson:"author"`
}
