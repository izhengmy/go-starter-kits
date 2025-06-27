package provider

import (
	"app/pkg/ginx"
	"app/pkg/gormx"
	"app/pkg/server"
	"app/pkg/translator"
	"app/pkg/zapx"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewGin(logger *zap.Logger, routes ginx.Routes) *gin.Engine {
	return ginx.New(logger, routes)
}

func NewServer(engine *gin.Engine) *server.Server {
	var config server.Config
	if err := viper.UnmarshalKey("server", &config); err != nil {
		panic(err)
	}
	return server.New(config, engine)
}

func NewZapLogger() (*zap.Logger, func()) {
	var config zapx.Config
	if err := viper.UnmarshalKey("zap", &config); err != nil {
		panic(err)
	}
	return zapx.NewLogger(config)
}

func NewGORMDB(zapLogger *zap.Logger) (*gorm.DB, func()) {
	var config gormx.Config
	if err := viper.UnmarshalKey("gorm", &config); err != nil {
		panic(err)
	}
	db, cleanup, err := gormx.NewDB(zapLogger, config)
	if err != nil {
		panic(err)
	}
	return db, cleanup
}

func NewTranslator() ut.Translator {
	locale := viper.GetString("app.locale")
	return translator.New(locale)
}
