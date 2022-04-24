package validator

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"regexp"
)

func ValidateMobile(f validator.FieldLevel) bool {
	mobile := f.Field().String()
	zap.S().Infof("mobile = ", mobile)
	// 使用正则表达式
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	zap.S().Infof("ok = ", ok)
	if !ok {
		return false
	}
	return true
}
