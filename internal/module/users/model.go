package users

import (
	"app/internal/module/commons/authutil"
	"app/pkg/gormx"
)

type User struct {
	gormx.Model
	Username string
	Password string
}

func (u User) ToAuthenticationUser() authutil.AuthenticationUser[uint] {
	return authutil.AuthenticationUser[uint]{
		ID:       u.ID,
		Username: u.Username,
	}
}
