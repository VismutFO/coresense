package validation

import "strconv"

func ValidateFieldType(value string, expectedType string) bool {
	switch expectedType {
	case "string":
		return true
	case "integer", "int":
		_, err := strconv.Atoi(value)
		return err == nil
	case "float":
		_, err := strconv.ParseFloat(value, 64)
		return err == nil
	default:
		return false
	}
}
