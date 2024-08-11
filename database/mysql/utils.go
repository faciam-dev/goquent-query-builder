package mysql

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

func (s *SQLUtils) EscapeIdentifierAliasedValue(sb []byte, value string) []byte {
	eoc := strings.Index(value, " as ")
	if eoc != -1 {
		sb = s.EscapeIdentifier(sb, value[:eoc])
		sb = append(sb, " as "...)
		sb = s.EscapeIdentifier(sb, value[eoc+4:])
		return sb
	} else {
		eoc = strings.Index(value, " AS ")
		if eoc != -1 {
			sb = s.EscapeIdentifier(sb, value[:eoc])
			sb = append(sb, " as "...)
			sb = s.EscapeIdentifier(sb, value[eoc+4:])
			return sb
		} else {
			sb = s.EscapeIdentifier(sb, value)
			return sb
		}
	}
}

func (s *SQLUtils) GetQueryBuilderStrategy() interfaces.QueryBuilderStrategy {
	return NewMySQLQueryBuilder()
}

func (s *SQLUtils) EscapeIdentifier(sb []byte, v string) []byte {
	if v != "*" {
		if eoc := strings.IndexByte(v, '.'); eoc != -1 {
			sb = append(sb, "`"...)
			if eo := strings.IndexByte(v, '`'); eo != -1 {
				sb = append(sb, strings.ReplaceAll(v[:eo], "`", "``")...)
				sb = append(sb, "`.`"...)
				sb = append(sb, strings.ReplaceAll(v[eo+1:eoc], "`", "``")...)
			} else {
				sb = append(sb, v[:eoc]...)
				sb = append(sb, "`.`"...)
				sb = append(sb, v[eoc+1:]...)
			}
			return append(sb, "`"...)
		} else {
			sb = append(sb, "`"...)
			if eo := strings.IndexByte(v, '`'); eo != -1 {
				sb = append(sb, strings.ReplaceAll(v[:eo], "`", "``")...)
				sb = append(sb, "`.`"...)
				sb = append(sb, strings.ReplaceAll(v[eo+1:], "`", "``")...)
			} else {
				sb = append(sb, v...)
			}
			return append(sb, "`"...)
		}
	}
	sb = append(sb, v...)
	return sb
}
