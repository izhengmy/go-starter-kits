package ginx

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Routing interface {
	Prefix() string
	Middleware() []gin.HandlerFunc
	Register(group *gin.RouterGroup)
}

type Routes []Routing

func New(logger *zap.Logger, routes Routes) *gin.Engine {
	engine := gin.New()

	engine.Use(ginzap.RecoveryWithZap(logger, true))

	for _, r := range routes {
		prefix := r.Prefix()
		middleware := r.Middleware()

		group := engine.Group(prefix)

		if len(middleware) > 0 {
			group.Use(middleware...)
		}

		r.Register(group)
	}

	return engine
}
