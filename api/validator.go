package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/nicodanke/bankTutorial/utils"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currencyString, ok := fl.Field().Interface().(string); ok {
		return utils.IsSupportedCurrency(currencyString)
	}
	return false;
}