package internal

import "strings"

// Popleft pops the first string off a slice
func Popleft(s []string) (string, []string) {
	if len(s) == 0 {
		return "", s
	}
	return strings.Trim(s[0], " "), s[1:]
}
