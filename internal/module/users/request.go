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

type ChangePasswordRequest struct {
	OldPassword     string `json:"oldPassword" binding:"required" label:"旧密码"`
	NewPassword     string `json:"newPassword" binding:"required,min=3,max=16" label:"新密码"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=NewPassword" label:"确认密码"`
}
