package utils

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	ARS = "ARS"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, ARS:
		return true
	}
	return false
}