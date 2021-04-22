package util

import "strings"

var (
	wordMapping = map[string]string{
		"http": "HTTP",
		"url":  "URL",
		"ip":   "IP",
	}
)

func translateWord(word string, initCase bool) string {
	if val, ok := wordMapping[word]; ok {
		return val
	}
	if initCase {
		return strings.Title(word)
	}
	return word
}

func ReverseDomain(domain string) string {
	s := strings.Split(domain, ".")

	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return strings.Join(s, ".")
}

// Converts a string to CamelCase
func ToCamel(s string) string {
	// s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	bits := []string{}
	for _, v := range s {
		if v == '_' || v == ' ' || v == '-' {
			bits = append(bits, n)
			n = ""
		} else {
			n += string(v)
		}
	}
	bits = append(bits, n)

	ret := ""
	for i, substr := range bits {
		ret += translateWord(substr, i != 0)
	}
	return ret
}

func ToClassname(s string) string {
	return translateWord(ToCamel(s), true)
}
