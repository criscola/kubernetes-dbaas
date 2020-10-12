package typeutil

// StringsToInterfaceSlice is a utility function which converts a slice of strings to []interface
func StringsToInterfaceSlice(values ...string) []interface{} {
	y := make([]interface{}, len(values))
	for i, v := range values {
		y[i] = v
	}
	return y
}
