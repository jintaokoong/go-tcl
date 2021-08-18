package utils

func GetFirstN(s string, n int) string {
	if len(s) >= n {
		return s[0:n]
	} else {
		return s
	}
}
