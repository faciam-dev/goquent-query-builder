package mysql

import (
	"regexp"
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

func (s *SQLUtils) EscapeIdentifierAliasedValue(value interface{}) string {
	if value == nil {
		return "NULL"
	}

	target := regexp.MustCompile(`(?i)\s+as\s+`)
	if target.MatchString(value.(string)) {
		parts := target.Split(value.(string), -1)
		return s.EscapeIdentifier(parts[0]) + " as " + s.EscapeIdentifier(parts[1])
	}

	return value.(string)
}

func (s *SQLUtils) EscapeIdentifier(value interface{}) string {
	if value == nil {
		return "NULL"
	}

	if v, ok := value.(string); ok {
		if strings.Contains(strings.ToLower(v), " as ") {
			return s.EscapeIdentifierAliasedValue(value)
		}

		if v != "*" {
			if strings.Contains(v, ".") {
				parts := strings.Split(v, ".")
				return "`" + strings.ReplaceAll(parts[0], "`", "``") + "`.`" + strings.ReplaceAll(parts[1], "`", "``") + "`"
			}
			return "`" + strings.ReplaceAll(v, "`", "``") + "`"
		}
	}

	return value.(string)
}

func (s *SQLUtils) GetQueryBuilderStrategy() interfaces.QueryBuilderStrategy {
	return NewMySQLQueryBuilder()
}
