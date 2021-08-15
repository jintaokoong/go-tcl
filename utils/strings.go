package utils

func GetFirstN(s string, n int) string {
	if len(s) >= n {
		return s[0:n]
	} else {
		return s
	}
}

func DeleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
