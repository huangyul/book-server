package bind

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator

func init() {
	InitTrans("zh")
}

func InitTrans(locale string) {
	if val, ok := binding.Validator.Engine().(*validator.Validate); ok {

		val.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New()
		enT := en.New()
		uni := ut.New(enT, zhT, enT)

		var ok bool
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			panic(fmt.Errorf("uni.GetTranslator(%s) failed", locale))
		}

		var err error

		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(val, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(val, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(val, trans)
		}

		if err != nil {
			panic(fmt.Sprintf("validator.RegisterDefaultTranslations(%s) failed", locale))
		}

	}
}

// Bind 绑定和验证请求数据
func Bind(ctx *gin.Context, obj any) error {
	err := ctx.Bind(obj)
	if err != nil {
		var v validator.ValidationErrors
		if errors.As(err, &v) {
			data, er := json.Marshal(removeTopStruct(v.Translate(trans)))
			if er != nil {
				return err
			}
			return errors.New(string(data))
		}
	}
	return nil
}

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
