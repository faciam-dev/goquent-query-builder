package interfaces

import "strings"

type SQLUtils interface {
	GetPlaceholder() string
	EscapeIdentifier(sb *strings.Builder, value interface{}) string
	EscapeIdentifierAliasedValue(sb *strings.Builder, value interface{}) string
	GetQueryBuilderStrategy() QueryBuilderStrategy
}
