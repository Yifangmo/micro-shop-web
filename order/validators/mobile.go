package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var mobileRe = regexp.MustCompile(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`)

func Mobile(fl validator.FieldLevel) bool {
	return mobileRe.MatchString(fl.Field().String())
}
