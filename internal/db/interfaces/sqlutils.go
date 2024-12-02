package interfaces

type SQLUtils interface {
	GetPlaceholder() string
	EscapeIdentifier(sb []byte, value string) []byte
	EscapeIdentifierAliasedValue(sb []byte, value string) []byte
	GetAlias(value string) string
	GetQueryBuilderStrategy() QueryBuilderStrategy
}
