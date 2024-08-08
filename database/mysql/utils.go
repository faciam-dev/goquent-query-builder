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

func (s *SQLUtils) EscapeIdentifierAliasedValue(sb []byte, value string) []byte {
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
			sb = s.EscapeIdentifier2(sb, pa)
			sb = append(sb, " as "...)
			sb = s.EscapeIdentifier2(sb, pb)
			return sb
			//return s.EscapeIdentifier(sb, split[0]) + " as " + s.EscapeIdentifier(sb, split[1])
		}
	} else {
		sb = s.EscapeIdentifier2(sb, value)
		return sb
	}

	target := regexp.MustCompile(`(?i)\s+as\s+`)
	if target.MatchString(value) {
		parts := target.Split(value, -1)
		sb = s.EscapeIdentifier2(sb, parts[0])
		sb = append(sb, " as "...)
		sb = s.EscapeIdentifier2(sb, parts[1])
		return sb
		//return s.EscapeIdentifier(sb, parts[0]) + " as " + s.EscapeIdentifier(sb, parts[1])
	}

	return append(sb, value...)
}

func (s *SQLUtils) EscapeIdentifier(sb *strings.Builder, v string) {
	if v != "*" {
		if eoc := strings.Index(v, "."); eoc != -1 {
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
			return
		} else {
			sb.WriteString("`")
			if eo := strings.Index(v, "`"); eo != -1 {
				sb.WriteString(strings.ReplaceAll(v[:eo], "`", "``"))
				sb.WriteString("`.`")
				sb.WriteString(strings.ReplaceAll(v[eo+1:], "`", "``"))
			} else {
				sb.WriteString(v)
			}
			sb.WriteString("`")
			return
		}
	}
	sb.WriteString(v)
}

func (s *SQLUtils) GetQueryBuilderStrategy() interfaces.QueryBuilderStrategy {
	return NewMySQLQueryBuilder()
}

func (s *SQLUtils) EscapeIdentifier2(sb []byte, v string) []byte {
	if v != "*" {
		if eoc := strings.Index(v, "."); eoc != -1 {
			sb = append(sb, "`"...)
			if eo := strings.Index(v, "`"); eo != -1 {
				sb = append(sb, strings.ReplaceAll(v[:eo], "`", "``")...)
				sb = append(sb, "`.`"...)
				sb = append(sb, strings.ReplaceAll(v[eo+1:eoc], "`", "``")...)
			} else {
				sb = append(sb, v[:eoc]...)
				sb = append(sb, "`.`"...)
				sb = append(sb, v[eoc+1:]...)
			}
			sb = append(sb, "`"...)
			/*
				sb = append(sb, strings.ReplaceAll(v[:eoc], "`", "``")...)
				sb = append(sb, "`.`"...)
				sb = append(sb, strings.ReplaceAll(v[eoc+1:], "`", "``")...)
				sb = append(sb, "`"...)
			*/
			return sb
		} else {
			sb = append(sb, "`"...)
			if eo := strings.Index(v, "`"); eo != -1 {
				sb = append(sb, strings.ReplaceAll(v[:eo], "`", "``")...)
				sb = append(sb, "`.`"...)
				sb = append(sb, strings.ReplaceAll(v[eo+1:], "`", "``")...)
			} else {
				sb = append(sb, v...)
			}
			sb = append(sb, "`"...)
			return sb
		}
	}
	sb = append(sb, v...)
	return sb
}
