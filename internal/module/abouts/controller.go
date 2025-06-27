package abouts

import (
	"app/internal/config"
	"app/internal/pkg/http"

	"github.com/gin-gonic/gin"
)

type AboutController struct {
	json   *http.JSON
	config config.Config
}

func NewAboutController(json *http.JSON, config config.Config) *AboutController {
	return &AboutController{
		json:   json,
		config: config,
	}
}

func (c AboutController) About(ctx *gin.Context) {
	c.json.Success(ctx, AboutResponse{
		Name:   c.config.App.Name,
		Env:    c.config.App.Env,
		Locale: c.config.App.Locale,
	})
}
