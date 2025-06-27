//go:build wireinject
// +build wireinject

package abouts

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	wire.NewSet(
		NewAboutController,
	),
)
