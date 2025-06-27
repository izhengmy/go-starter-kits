//go:build wireinject
// +build wireinject

package pkg

import (
	"app/internal/pkg/http"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	http.NewJSON,
)
