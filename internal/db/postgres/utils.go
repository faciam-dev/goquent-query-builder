package postgres

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
	return "$1"
}

func (s *SQLUtils) EscapeIdentifier(value interface{}) string {
	if value == nil {
		return "NULL"
	}

	if v, ok := value.(string); ok {
		if v != "*" {
			if strings.Contains(v, ".") {
				parts := strings.Split(v, ".")
				return `"` + strings.ReplaceAll(parts[0], `"`, `""`) + `"."` + strings.ReplaceAll(parts[1], `"`, `""`) + `"`
			}
			return `"` + strings.ReplaceAll(v, `"`, `""`) + `"`
		}
	}

	return value.(string)
}

func (s *SQLUtils) GetQueryBuilderStrategy() interfaces.QueryBuilderStrategy {
	return NewPostgreSQLQueryBuilder()
}
