package photon

// btoi converts a boolean into an integer.
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

// itob converts an integer into a boolean.
func itob(i int) bool {
	return i > 0
}
