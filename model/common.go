package model

type MainPageResponse struct {
	User         UserClaims   `json:"user"`
	LinkCategory LinkCategory `json:"link_category"`
}
