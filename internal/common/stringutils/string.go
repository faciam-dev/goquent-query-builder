package stringutils

import "strings"

// EscapeString escapes a string for use in a SQL query
func EscapeString(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, "'", "''"), "\\", "\\\\")
}
