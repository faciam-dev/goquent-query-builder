package interfaces

import "strings"

type SQLUtils interface {
	GetPlaceholder() string
	EscapeIdentifier2(sb []byte, value string) []byte
	EscapeIdentifier(sb *strings.Builder, value string)
	EscapeIdentifierAliasedValue(sb *strings.Builder, value string)
	GetQueryBuilderStrategy() QueryBuilderStrategy
}
