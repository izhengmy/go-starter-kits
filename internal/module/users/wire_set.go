//go:build wireinject
// +build wireinject

package users

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	wire.NewSet(
		NewUserController,
	),
	wire.NewSet(
		NewUserService,
	),
	wire.NewSet(
		NewUserRepository,
	),
)
