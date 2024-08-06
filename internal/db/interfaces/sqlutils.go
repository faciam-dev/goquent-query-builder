package interfaces

import "strings"

type SQLUtils interface {
	GetPlaceholder() string
	EscapeIdentifier(sb *strings.Builder, value string) string
	EscapeIdentifierAliasedValue(sb *strings.Builder, value string) string
	GetQueryBuilderStrategy() QueryBuilderStrategy
}
