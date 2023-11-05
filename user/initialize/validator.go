package initialize

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Yifangmo/micro-shop-web/user/global"
	"github.com/Yifangmo/micro-shop-web/user/validators"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

func InitValidator() {
	v := binding.Validator.Engine().(*validator.Validate)
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	zhT := zh.New()
	enT := en.New()
	uni := ut.New(enT, zhT, enT)
	var ok bool
	global.Translators, ok = uni.GetTranslator(global.Locale)
	if !ok {
		panic(fmt.Errorf("uni.GetTranslator(%s)", global.Locale))
	}

	switch global.Locale {
	case "en":
		en_translations.RegisterDefaultTranslations(v, global.Translators)
	case "zh":
		zh_translations.RegisterDefaultTranslations(v, global.Translators)
	default:
		en_translations.RegisterDefaultTranslations(v, global.Translators)
	}
	_ = v.RegisterValidation("mobile", validators.Mobile)
	_ = v.RegisterTranslation("mobile", global.Translators, func(ut ut.Translator) error {
		return ut.Add("mobile", "{0} 非法的手机号码!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("mobile", fe.Field())
		return t
	})
}
