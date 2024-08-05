package base

import (
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

type DeleteBaseBuilder struct {
	u interfaces.SQLUtils
}

func NewDeleteBaseBuilder(util interfaces.SQLUtils, iq *structs.DeleteQuery) *DeleteBaseBuilder {
	return &DeleteBaseBuilder{
		u: util,
	}
}

func (m *DeleteBaseBuilder) Delete(q *structs.DeleteQuery) *DeleteBaseBuilder {
	return m
}

// DeleteBatch builds the Delete query for Delete.
func (m *DeleteBaseBuilder) BuildDelete(q *structs.DeleteQuery) (string, []interface{}, error) {
	values := make([]interface{}, 0)
	sb := &strings.Builder{}
	sb.Grow(consts.StringBuffer_Delete_Grow)

	// DELETE
	sb.WriteString("DELETE")
	if q.Query.Joins != nil &&
		q.Query.Joins.Joins != nil &&
		(len(*q.Query.Joins.Joins) > 0 || (q.Query.Joins.JoinClause != nil && len(*q.Query.Joins.JoinClause) > 0)) {
		sb.WriteString(" ")
		sb.WriteString(m.u.EscapeIdentifier(sb, q.Table))
	}

	// FROM
	sb.WriteString(" FROM ")
	sb.WriteString(m.u.EscapeIdentifier(sb, q.Table))

	// JOIN
	jb := NewJoinBaseBuilder(m.u, q.Query.Joins)
	jb.Join(sb, q.Query.Joins)

	// WHERE
	if len(*q.Query.ConditionGroups) > 0 {
		wb := NewWhereBaseBuilder(m.u, q.Query.ConditionGroups)
		whereValues := wb.Where(sb, q.Query.ConditionGroups)
		values = append(values, whereValues...)
	}

	// ORDER BY
	if len(*q.Query.Order) > 0 {
		ob := NewOrderByBaseBuilder(m.u, q.Query.Order)
		ob.OrderBy(sb, q.Query.Order)
	}

	query := sb.String()
	sb.Reset()

	return query, values, nil
}
