package stringutils

import "strings"

// EscapeString escapes a string for use in a SQL query
func EscapeString(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, "'", "''"), "\\", "\\\\")
}

func ManualSplit(s, sep string) []string {
	var res []string
	start := 0
	for i := 0; i+len(sep) <= len(s); i++ {
		if s[i:i+len(sep)] == sep {
			res = append(res, s[start:i])
			i += len(sep) - 1
			start = i + 1
		}
	}
	if start < len(s) {
		res = append(res, s[start:])
	}
	return res
}
