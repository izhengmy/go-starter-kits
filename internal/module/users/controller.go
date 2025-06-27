package users

import (
	"app/internal/module/commons/authutil"
	"app/internal/pkg/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	json    *http.JSON
	service *UserService
}

func NewUserController(json *http.JSON, service *UserService) *UserController {
	return &UserController{
		json:    json,
		service: service,
	}
}

func (c UserController) Register(ctx *gin.Context) {
	var request RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.json.Fail(ctx, err)
		return
	}
	err := c.service.Register(request)
	if err != nil {
		c.json.Fail(ctx, err)
		return
	}
	c.json.Success(ctx, nil)
}

func (c UserController) Login(ctx *gin.Context) {
	var request LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.json.Fail(ctx, err)
		return
	}
	response, err := c.service.Login(request)
	if err != nil {
		c.json.Fail(ctx, err)
		return
	}
	c.json.Success(ctx, response)
}

func (c UserController) Profile(ctx *gin.Context) {
	authUser := authutil.GetAuthenticationUser(ctx)
	response, err := c.service.Profile(*authUser)
	if err != nil {
		c.json.Fail(ctx, err)
		return
	}
	c.json.Success(ctx, response)
}

func (c UserController) ChangePassword(ctx *gin.Context) {
	var request ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.json.Fail(ctx, err)
		return
	}
	authUser := authutil.GetAuthenticationUser(ctx)
	err := c.service.ChangePassword(*authUser, request)
	if err != nil {
		c.json.Fail(ctx, err)
		return
	}
	c.json.Success(ctx, nil)
}
