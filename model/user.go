package model

type User struct {
	Id           string `json:"id,omitempty" bson:"_id,omitempty"`
	Username     string `json:"username bson:"username"`
	Email        string `json:"email bson:"email"`
	FirstName    string `json:"first_name" bson:"first_name"`
	SecondName   string `json:"second_name" bson:"second_name"`
	PasswordHash string `json:"password_hash" bson:"password_hash"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Username string ``
}
