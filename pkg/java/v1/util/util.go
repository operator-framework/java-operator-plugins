package util

import "strings"

func ReverseDomain(domain string) string {
	s := strings.Split(domain, ".")

	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return strings.Join(s, ".")
}
