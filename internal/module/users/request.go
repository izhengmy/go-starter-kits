package users

type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=10" label:"用户名"`
	Password        string `json:"password" binding:"required,min=3,max=16" label:"密码"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password" label:"确认密码"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" label:"用户名"`
	Password string `json:"password" binding:"required" label:"密码"`
}
