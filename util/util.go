package util

func IsTerminal(a byte) bool {
	if a < 'A' || a > 'Z' {
		return true
	}
	return false
}
