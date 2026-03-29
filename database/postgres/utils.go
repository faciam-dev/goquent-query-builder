package postgres

import (
	"strconv"
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/sqlutils"
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

func (s *SQLUtils) ResetPlaceholderCounter() {
	s.placeholderNumber = 0
}

func (s *SQLUtils) EscapeRelation(sb []byte, value string) []byte {
	return sqlutils.AppendEscapedRelation(sb, value, '"')
}

func (s *SQLUtils) EscapeReference(sb []byte, value string) []byte {
	return sqlutils.AppendEscapedReference(sb, value, '"')
}

func (s *SQLUtils) EscapeAliasedValue(sb []byte, value string) []byte {
	return sqlutils.AppendEscapedAliasedValue(sb, value, '"')
}

func (s *SQLUtils) GetQueryBuilderStrategy() interfaces.QueryBuilderStrategy {
	return newPostgreSQLQueryBuilderWithUtil(s)
}

func (s *SQLUtils) Dialect() string {
	return consts.DialectPostgreSQL
}
