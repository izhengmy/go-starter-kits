package translator

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
)

func New(locale string) ut.Translator {
	engine, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil
	}

	engine.RegisterTagNameFunc(func(field reflect.StructField) string {
		if text := field.Tag.Get("json"); text != "" {
			name := strings.SplitN(text, ",", 2)[0]
			if name == "-" || name == "" {
				return field.Name
			} else {
				return name
			}
		}
		return field.Name
	})

	zhT := zh.New()
	enT := en.New()
	uni := ut.New(enT, zhT, enT)

	translator, ok := uni.GetTranslator(locale)
	if !ok {
		return nil
	}

	switch locale {
	case "en":
		_ = enTrans.RegisterDefaultTranslations(engine, translator)
	case "zh":
		_ = zhTrans.RegisterDefaultTranslations(engine, translator)
	default:
		_ = enTrans.RegisterDefaultTranslations(engine, translator)
	}

	return translator
}
