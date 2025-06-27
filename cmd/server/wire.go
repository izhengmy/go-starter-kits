//go:build wireinject
// +build wireinject

package main

import (
	"app/internal/config"
	"app/internal/module/abouts"
	"app/internal/module/users"
	"app/internal/pkg"
	"app/internal/routes"
	"app/pkg/provider"
	"app/pkg/server"

	"github.com/google/wire"
)

func wireServer(conf config.Config) (*server.Server, func(), error) {
	panic(wire.Build(
		provider.NewGin,
		provider.NewServer,
		provider.NewZapLogger,
		provider.NewGORMDataSources,
		provider.NewTranslator,
		routes.ProviderSet,
		pkg.ProviderSet,
		abouts.ProviderSet,
		users.ProviderSet,
	))
}
