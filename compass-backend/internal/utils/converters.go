package utils

// BoolToString converts a boolean value to "Yes" or "No" string
func BoolToString(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}