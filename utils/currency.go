package utils

const (
	INR = "INR"
	USD = "USD"
	CAD = "CAD"
	EUR = "EUR"
)

// returns true if currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case INR, USD, EUR, CAD:
		return true
	}
	return false
}
