package postgres

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
	return "$1"
}

func (s *SQLUtils) EscapeIdentifierAliasedValue(sb *strings.Builder, value interface{}) string {
	if value == nil {
		return "NULL"
	}

	target := regexp.MustCompile(`(?i)\s+as\s+`)
	if target.MatchString(value.(string)) {
		parts := target.Split(value.(string), -1)
		return s.EscapeIdentifier(sb, parts[0]) + " as " + s.EscapeIdentifier(sb, parts[1])
	}

	return value.(string)
}

func (s *SQLUtils) EscapeIdentifier(sb *strings.Builder, value interface{}) string {
	if value == nil {
		return "NULL"
	}

	if v, ok := value.(string); ok {
		if strings.Contains(strings.ToLower(v), " as ") {
			split := strings.Split(v, " as ")
			if len(split) != 2 {
				split = strings.Split(v, " AS ")
			}
			if len(split) == 2 {
				return s.EscapeIdentifier(sb, split[0]) + " as " + s.EscapeIdentifier(sb, split[1])
			}
			return s.EscapeIdentifierAliasedValue(sb, value)
		}

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
