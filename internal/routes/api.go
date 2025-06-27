package routes

import (
	"app/internal/errorx"
	"app/internal/module/abouts"
	"app/internal/module/commons/authutil"
	"app/internal/module/commons/middleware"
	"app/internal/module/users"
	pkgHTTP "app/internal/pkg/http"
	"app/pkg/auth"
	"app/pkg/ginx"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type APIRouting struct {
	logger          *zap.Logger
	json            *pkgHTTP.JSON
	aboutController *abouts.AboutController
	userController  *users.UserController
}

var _ ginx.Routing = (*APIRouting)(nil)

func NewAPIRouting(
	logger *zap.Logger,
	json *pkgHTTP.JSON,
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
		jwtAuth(r),
	}
}

func (r APIRouting) Register(router *gin.RouterGroup) {
	{
		router.GET("/abouts", r.aboutController.About)
	}

	{
		router.POST("/users/register", r.userController.Register)
		router.POST("/users/login", r.userController.Login)
		router.GET("/users/profile", r.userController.Profile)
	}
}

func jwtAuth(r APIRouting) gin.HandlerFunc {
	authUser := authutil.AuthenticationUser[uint]{}

	return auth.GinJWTAuthHandlerFunc(
		&authUser,
		auth.WithContextKey("user"),
		auth.WithSkips([]string{
			r.Prefix() + "/users/register",
			r.Prefix() + "/users/login",
		}),
		auth.WithUnauthenticatedFunc(func(ctx *gin.Context, err error) {
			r.json.Fail(ctx, errorx.NewServiceError("登录已过期").WithCode(http.StatusUnauthorized))
			ctx.AbortWithStatus(http.StatusOK)
		}),
	)
}
