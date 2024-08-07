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

func (s *SQLUtils) EscapeIdentifierAliasedValue(sb *strings.Builder, value string) string {
	eoc := strings.Index(strings.ToLower(value), " as ")
	if eoc != -1 {
		eoc := strings.Index(value, " as ")
		var pa, pb string
		pa = value[:eoc]
		pb = value[eoc+4:]
		if eoc == -1 {
			eoc = strings.Index(value, " AS ")
			pa = value[:eoc]
			pb = value[eoc+4:]
		}
		//split := strings.Split(v, " as ")
		//if len(split) != 2 {
		//	split = strings.Split(v, " AS ")
		//}
		if eoc != -1 {
			sb.WriteString(s.EscapeIdentifier(sb, pa))
			sb.WriteString(" as ")
			sb.WriteString(s.EscapeIdentifier(sb, pb))
			return ""
			//return s.EscapeIdentifier(sb, split[0]) + " as " + s.EscapeIdentifier(sb, split[1])
		}
	} else {
		return s.EscapeIdentifier(sb, value)
	}

	target := regexp.MustCompile(`(?i)\s+as\s+`)
	if target.MatchString(value) {
		parts := target.Split(value, -1)
		sb.WriteString(s.EscapeIdentifier(sb, parts[0]))
		sb.WriteString(" as ")
		sb.WriteString(s.EscapeIdentifier(sb, parts[1]))
		return ""
		//return s.EscapeIdentifier(sb, parts[0]) + " as " + s.EscapeIdentifier(sb, parts[1])
	}

	return value
}

func (s *SQLUtils) EscapeIdentifier(sb *strings.Builder, v string) string {
	if v != "*" {
		if eoc := strings.Index(v, "."); eoc != -1 {
			//sb.Grow(len(v) + 4) // 事前にバッファを拡張
			sb.WriteString("`")
			if eo := strings.Index(v, "`"); eo != -1 {
				sb.WriteString(strings.ReplaceAll(v[:eo], "`", "``"))
				sb.WriteString("`.`")
				sb.WriteString(strings.ReplaceAll(v[eo+1:eoc], "`", "``"))
			} else {
				sb.WriteString(v[:eoc])
				sb.WriteString("`.`")
				sb.WriteString(v[eoc+1:])
			}
			sb.WriteString("`")
			/*
				sb.WriteString(strings.ReplaceAll(v[:eoc], "`", "``"))
				sb.WriteString("`.`")
				sb.WriteString(strings.ReplaceAll(v[eoc+1:], "`", "``"))
				sb.WriteString("`")
			*/
			return ""
		}
		//sb.Grow(len(v) + 2) // 事前にバッファを拡張
		sb.WriteString("`")
		sb.WriteString(strings.ReplaceAll(v, "`", "``"))
		sb.WriteString("`")
		return ""
	}
	return v
}

func (s *SQLUtils) GetQueryBuilderStrategy() interfaces.QueryBuilderStrategy {
	return NewMySQLQueryBuilder()
}
