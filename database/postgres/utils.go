package postgres

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

type SQLUtils struct {
	placeholderNumber int
}

func NewSQLUtils() *SQLUtils {
	return &SQLUtils{
		placeholderNumber: 0,
	}
}

func (s *SQLUtils) GetPlaceholder() string {
	s.placeholderNumber++
	phn := strconv.Itoa(s.placeholderNumber)
	return strings.Join([]string{"$", phn}, "")
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
			sb = s.EscapeIdentifier(sb, pa)
			sb = append(sb, " as "...)
			sb = s.EscapeIdentifier(sb, pb)
			return sb
			//return s.EscapeIdentifier(sb, split[0]) + " as " + s.EscapeIdentifier(sb, split[1])
		}
	} else {
		sb = s.EscapeIdentifier(sb, value)
		return sb
	}

	target := regexp.MustCompile(`(?i)\s+as\s+`)
	if target.MatchString(value) {
		parts := target.Split(value, -1)
		sb = s.EscapeIdentifier(sb, parts[0])
		sb = append(sb, " as "...)
		sb = s.EscapeIdentifier(sb, parts[1])
		return sb
		//return s.EscapeIdentifier(sb, parts[0]) + " as " + s.EscapeIdentifier(sb, parts[1])
	}

	return append(sb, value...)
}

func (s *SQLUtils) GetQueryBuilderStrategy() interfaces.QueryBuilderStrategy {
	return NewPostgreSQLQueryBuilder()
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
