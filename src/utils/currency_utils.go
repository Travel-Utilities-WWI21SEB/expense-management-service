package utils

func IsValidCurrencyCode(currencyCode string) bool {
	validCurrencyCodes := []string{
		"EUR",
		"USD",
		"GBP",
	}

	// Check if currency code is valid
	for _, code := range validCurrencyCodes {
		if code == currencyCode {
			return true
		}
	}

	return false
}
