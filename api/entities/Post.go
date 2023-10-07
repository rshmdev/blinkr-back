package entities

type Post struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	PostAt      string `json:"post_at"`
	Author      string `json:"author"`
}
