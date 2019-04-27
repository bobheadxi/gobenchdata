package internal

import "strings"

// Popleft pops the first string off a slice
func Popleft(s []string) (string, []string) {
	return strings.Trim(s[0], " "), s[1:]
}
