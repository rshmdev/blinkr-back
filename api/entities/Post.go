package entities

type Post struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	PostAt      string   `json:"post_at"`
	Author      string   `json:"author"`
	Media       []string `json:"media"`
	AuthorInfo  User     `json:"author_info"` // Adicione este campo para armazenar as informações do autor

}

type ImagePost struct {
	MediaId string `json:"media_id"`
}

type Replys struct {
	ID          string `json:"id"`
	ReplyFrom   string `json:"replay_from"`
	Author      string `json:"author"`
	Description string `json:"description"`
}
