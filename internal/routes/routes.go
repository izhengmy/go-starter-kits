package routes

import "app/pkg/ginx"

func New(apiRouting *APIRouting) ginx.Routes {
	return ginx.Routes{
		apiRouting,
	}
}
