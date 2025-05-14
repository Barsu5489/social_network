package utils

// Convert a boolean to SQLite integer (1 for true, 0 for false)
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Returns nil for nil values or Unix timestamp for populated values
func NilOrNullInt(t *int64) interface{} {
	if t == nil {
		return nil
	}
	return *t
}
