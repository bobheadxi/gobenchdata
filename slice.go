package main

import "strings"

func popleft(s []string) (string, []string) {
	return strings.Trim(s[0], " "), s[1:]
}
