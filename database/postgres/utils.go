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

func (s *SQLUtils) EscapeIdentifierAliasedValue(sb *strings.Builder, value string) {
	if strings.Contains(strings.ToLower(value), " as ") {
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
			s.EscapeIdentifier(sb, pa)
			sb.WriteString(" as ")
			s.EscapeIdentifier(sb, pb)
			return
			//return s.EscapeIdentifier(sb, split[0]) + " as " + s.EscapeIdentifier(sb, split[1])
		}
	} else {
		s.EscapeIdentifier(sb, value)
		return
	}

	target := regexp.MustCompile(`(?i)\s+as\s+`)
	if target.MatchString(value) {
		parts := target.Split(value, -1)
		s.EscapeIdentifier(sb, parts[0])
		sb.WriteString(" as ")
		s.EscapeIdentifier(sb, parts[1])
		return
		//return s.EscapeIdentifier(sb, parts[0]) + " as " + s.EscapeIdentifier(sb, parts[1])
	}
	sb.WriteString(value)
}

func (s *SQLUtils) EscapeIdentifier(sb *strings.Builder, v string) {
	//if v, ok := value.(string); ok {
	if v != "*" {
		if strings.Contains(v, ".") {
			eoc := strings.Index(v, ".")
			var pa, pb string
			pa = v[:eoc]
			pb = v[eoc+1:]
			//parts := strings.Split(v, ".")
			sb.WriteString(`"`)
			sb.WriteString(strings.ReplaceAll(pa, `"`, `""`))
			sb.WriteString(`"."`)
			sb.WriteString(strings.ReplaceAll(pb, "`", "``"))
			sb.WriteString(`"`)
			return
			//return "`" + strings.ReplaceAll(parts[0], "`", "``") + "`.`" + strings.ReplaceAll(parts[1], "`", "``") + "`"
		}
		sb.WriteString(`"`)
		sb.WriteString(strings.ReplaceAll(v, `"`, `""`))
		sb.WriteString(`"`)
		return
		//return "`" + strings.ReplaceAll(v, "`", "``") + "`"
	}
	//}
	sb.WriteString(v)
}

func (s *SQLUtils) GetQueryBuilderStrategy() interfaces.QueryBuilderStrategy {
	return NewPostgreSQLQueryBuilder()
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
