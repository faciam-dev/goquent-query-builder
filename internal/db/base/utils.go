package base

import (
	"regexp"
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

var asRegexp = regexp.MustCompile(`(?i)\s+as\s+`)

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
		if eoc != -1 {
			sb = s.EscapeIdentifier(sb, pa)
			sb = append(sb, " as "...)
			sb = s.EscapeIdentifier(sb, pb)
			return sb
		}
	} else {
		sb = s.EscapeIdentifier(sb, value)
		return sb
	}

	if asRegexp.MatchString(value) {
		parts := asRegexp.Split(value, -1)
		sb = s.EscapeIdentifier(sb, parts[0])
		sb = append(sb, " as "...)
		sb = s.EscapeIdentifier(sb, parts[1])
		return sb
	}

	return append(sb, value...)
}

func (s *SQLUtils) GetQueryBuilderStrategy() interfaces.QueryBuilderStrategy {
	return NewBaseQueryBuilder()
}

func (s *SQLUtils) Dialect() string {
	return consts.DialectBase
}

func (s *SQLUtils) EscapeIdentifier(sb []byte, v string) []byte {
	if v != "*" {
		if eoc := strings.Index(v, "."); eoc != -1 {
			sb = append(sb, `"`...)
			if eo := strings.Index(v, `"`); eo != -1 {
				sb = append(sb, strings.ReplaceAll(v[:eo], `"`, `""`)...)
				sb = append(sb, `"."`...)
				sb = append(sb, strings.ReplaceAll(v[eo+1:eoc], `"`, `""`)...)
			} else {
				sb = append(sb, v[:eoc]...)
				sb = append(sb, `"."`...)
				sb = append(sb, v[eoc+1:]...)
			}
			sb = append(sb, `"`...)
			return sb
		} else {
			sb = append(sb, `"`...)
			if eo := strings.Index(v, `"`); eo != -1 {
				sb = append(sb, strings.ReplaceAll(v[:eo], `"`, `""`)...)
				sb = append(sb, `"."`...)
				sb = append(sb, strings.ReplaceAll(v[eo+1:], `"`, `""`)...)
			} else {
				sb = append(sb, v...)
			}
			sb = append(sb, `"`...)
			return sb
		}
	}
	sb = append(sb, v...)
	return sb
}

func (s *SQLUtils) GetAlias(value string) string {
	eoc := strings.Index(value, " as ")
	if eoc != -1 {
		return value[eoc+4:]
	} else {
		eoc = strings.Index(value, " AS ")
		if eoc != -1 {
			return value[eoc+4:]
		}
	}

	if asRegexp.MatchString(value) {
		parts := asRegexp.Split(value, -1)
		return parts[1]
	}

	return value
}
