package users

type LoginResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"tokenType"`
}

type ProfileResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}
