package entities

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Birthday  string `json:"birthday"`
	Avatar    string `json:"avatar"`
}

type LoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
