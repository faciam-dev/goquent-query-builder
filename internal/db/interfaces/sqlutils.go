package interfaces

type SQLUtils interface {
	GetPlaceholder() string
	EscapeIdentifier(value interface{}) string
	GetQueryBuilderStrategy() QueryBuilderStrategy
}
