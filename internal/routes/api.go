package routes

import (
	"app/internal/module/abouts"
	"app/internal/module/commons/authutil"
	"app/internal/module/commons/middleware"
	"app/internal/module/users"
	"app/internal/pkg/http"
	"app/pkg/auth"
	"app/pkg/ginx"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type APIRouting struct {
	logger          *zap.Logger
	json            *http.JSON
	aboutController *abouts.AboutController
	userController  *users.UserController
}

var _ ginx.Routing = (*APIRouting)(nil)

func NewAPIRouting(
	logger *zap.Logger,
	json *http.JSON,
	aboutController *abouts.AboutController,
	userController *users.UserController,
) *APIRouting {
	return &APIRouting{
		logger:          logger,
		json:            json,
		aboutController: aboutController,
		userController:  userController,
	}
}

func (r APIRouting) Prefix() string {
	return "/api"
}

func (r APIRouting) Middleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.Logger(r.logger),
	}
}

func (r APIRouting) Register(router *gin.RouterGroup) {
	{
		group := router.Group("")

		group.GET("/abouts", r.aboutController.About)

		group.POST("/users/register", r.userController.Register)
		group.POST("/users/login", r.userController.Login)
	}

	{
		group := router.Group("", jwtAuth(r, true))
		group.GET("/users/profile", r.userController.Profile)
	}
}

func jwtAuth(r APIRouting, abort bool) gin.HandlerFunc {
	authUser := authutil.AuthenticationUser[uint]{}

	return auth.GinJWTAuthHandlerFunc(
		&authUser,
		abort,
		auth.WithContextKey(authutil.CtxKeyUser),
		auth.WithUnauthenticatedFunc(authutil.Unauthenticated(r.json)),
	)
}
