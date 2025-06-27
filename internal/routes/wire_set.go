//go:build wireinject
// +build wireinject

package routes

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	New,
	NewAPIRouting,
)
