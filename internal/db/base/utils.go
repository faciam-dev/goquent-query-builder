package base

import (
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/sqlutils"
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
	return newBaseQueryBuilderWithUtil(s)
}

func (s *SQLUtils) Dialect() string {
	return consts.DialectBase
}
