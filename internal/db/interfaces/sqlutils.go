package interfaces

type SQLUtils interface {
	GetPlaceholder() string
	EscapeIdentifier(value interface{}) string
	EscapeIdentifierAliasedValue(value interface{}) string
	GetQueryBuilderStrategy() QueryBuilderStrategy
}
