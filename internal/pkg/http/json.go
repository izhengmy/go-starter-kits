package http

import (
	"app/internal/errorx"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type jsonResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    any                `json:"data,omitempty"`
	Errors  *map[string]string `json:"errors,omitempty"`
}

type JSON struct {
	translator ut.Translator
}

func NewJSON(translator ut.Translator) *JSON {
	return &JSON{
		translator: translator,
	}
}

func (j JSON) Success(ctx *gin.Context, data any, args ...string) {
	message := "success"

	if len(args) > 0 && args[0] != "" {
		message = args[0]
	}

	if data == nil {
		data = map[string]any{}
	}

	ctx.JSON(http.StatusOK, jsonResponse{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

func (j JSON) Fail(ctx *gin.Context, err error, args ...string) {
	_ = ctx.Error(err)

	var response jsonResponse
	var serviceError *errorx.ServiceError
	var validationErrors validator.ValidationErrors
	switch {
	case errors.As(err, &serviceError):
		response = j.handleServiceError(serviceError)
	case errors.As(err, &validationErrors):
		response = j.handleValidationErrors(validationErrors)
	default:
		response = jsonResponse{
			Code:    http.StatusInternalServerError,
			Message: "服务器错误",
			Errors:  &map[string]string{},
		}
	}

	if len(args) > 0 && args[0] != "" {
		response.Message = args[0]
	}

	ctx.JSON(http.StatusOK, response)
}

func (j JSON) handleServiceError(err *errorx.ServiceError) jsonResponse {
	return jsonResponse{
		Code:    err.Code,
		Message: err.Message,
		Errors:  &map[string]string{},
	}
}

func (j JSON) handleValidationErrors(errs validator.ValidationErrors) jsonResponse {
	messages := map[string]string{}
	for k, e := range errs.Translate(j.translator) {
		if idx := strings.Index(k, "."); idx != -1 {
			k = k[idx+1:]
		}
		messages[k] = e
	}
	return jsonResponse{
		Code:    http.StatusUnprocessableEntity,
		Message: "请求参数错误",
		Errors:  &messages,
	}
}
