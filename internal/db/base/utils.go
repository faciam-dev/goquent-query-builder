package base

import (
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

type SQLUtils struct {
}

func NewSQLUtils() *SQLUtils {
	return &SQLUtils{}
}

func (s *SQLUtils) GetPlaceholder() string {
	return "?"
}

func (s *SQLUtils) EscapeIdentifier(value interface{}) string {
	if value == nil {
		return "NULL"
	}

	if v, ok := value.(string); ok {
		if v != "*" {
			return `"` + strings.ReplaceAll(v, `"`, `""`) + `"`
		}
	}

	return value.(string)
}

func (s *SQLUtils) GetQueryBuilderStrategy() interfaces.QueryBuilderStrategy {
	return NewBaseQueryBuilder()
}
