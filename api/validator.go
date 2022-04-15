package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/skamranahmed/banking-system/utils"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	currency, ok := fieldLevel.Field().Interface().(string)
	if ok {
		// check if currency is supported or not
		return utils.IsSupportedCurrency(currency)
	}

	return false
}
