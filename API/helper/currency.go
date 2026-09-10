package helper

func Currency(Currency string) string {
	switch Currency {
	case "USD":
		return "ដុល្លា"
	case "KHR":
		return "រៀល"
	default:
		return ""
	}
}
