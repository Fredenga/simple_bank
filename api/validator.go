package api

import (
	"github.com/fredrick/simplebank/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		//check if currency is supported, if ok is false, field is not a String
		return util.IsSupportedCurrency(currency)
	}
	return false
}