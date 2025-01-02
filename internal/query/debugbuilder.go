package query

import (
	"fmt"
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/stringutils"
)

type DebugBuilder[T BaseBuilder, C any] struct {
	queryBuilder T
	child        *C
}

func NewDebugBuilder[T BaseBuilder, C any](queryBuilder T) *DebugBuilder[T, C] {
	return &DebugBuilder[T, C]{
		queryBuilder: queryBuilder,
	}
}

func (b *DebugBuilder[T, C]) SetChild(child *C) *C {
	b.child = child

	return b.child
}

// Dump returns the query and values.
func (b *DebugBuilder[T, C]) Dump() (string, []interface{}, error) {
	return b.queryBuilder.Build()
}

// RawSql returns the raw SQL query.
func (b *DebugBuilder[T, C]) RawSql() (string, error) {
	query, values, err := b.queryBuilder.Build()

	if err != nil {
		return "", err
	}

	return replacePlaceholders(query, values)
}

// replacePlaceholders replaces placeholders with the actual values.
func replacePlaceholders(query string, args []interface{}) (string, error) {
	placeholderCount := strings.Count(query, "?")
	if placeholderCount != len(args) {
		return "", fmt.Errorf("placeholder count does not match the number of arguments: %d != %d", placeholderCount, len(args))
	}

	for _, arg := range args {
		var replacement string
		switch v := arg.(type) {
		case string:
			replacement = fmt.Sprintf("'%s'", stringutils.EscapeString(v))
		case int, int64, float64:
			replacement = fmt.Sprintf("%v", v)
		case nil:
			replacement = "NULL"
		default:
			return "", fmt.Errorf("not supported type: %T", v)
		}

		query = strings.Replace(query, "?", replacement, 1)
	}

	return query, nil
}
