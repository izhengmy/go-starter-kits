package middleware

import (
	"app/internal/errorx"
	"errors"
	"fmt"
	"net/http/httputil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		end := time.Now()
		latency := end.Sub(start)
		request, _ := httputil.DumpRequest(ctx.Request, false)

		fields := []zapcore.Field{
			zap.String("request", string(request)),
			zap.String("clientIP", ctx.ClientIP()),
			zap.Int("statusCode", ctx.Writer.Status()),
			zap.Duration("latency", latency),
		}

		var err error
		var serviceError *errorx.ServiceError
		var validationErrors validator.ValidationErrors

		if len(ctx.Errors) > 0 {
			err = ctx.Errors.Last().Err
		}

		switch {
		case errors.As(err, &serviceError):
			err = nil
		case errors.As(err, &validationErrors):
			err = nil
		}

		if err != nil {
			fields = append(fields, zap.String("error", fmt.Sprintf("%+v", err)))
			logger.Error("server error", fields...)
		} else {
			logger.Debug("http request", fields...)
		}
	}
}
