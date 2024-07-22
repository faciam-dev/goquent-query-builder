package db

import (
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type DeleteBaseBuilder struct {
}

func NewDeleteBaseBuilder(iq *structs.DeleteQuery) *DeleteBaseBuilder {
	return &DeleteBaseBuilder{}
}

func (m *DeleteBaseBuilder) Delete(q *structs.DeleteQuery) *DeleteBaseBuilder {
	return m
}

// DeleteBatch builds the Delete query for Delete.
func (m *DeleteBaseBuilder) BuildDelete(q *structs.DeleteQuery) (string, []interface{}) {
	values := make([]interface{}, 0)
	sb := &strings.Builder{}
	sb.Grow(consts.StringBuffer_Delete_Grow)

	// DELETE
	sb.WriteString("DELETE")
	if q.Query.Joins != nil &&
		q.Query.Joins.Joins != nil &&
		(len(*q.Query.Joins.Joins) > 0 || q.Query.Joins.JoinClause != nil) {
		sb.WriteString(" ")
		sb.WriteString(q.Table)
	}

	// FROM
	sb.WriteString(" FROM " + q.Table)

	// JOIN
	b := &BaseQueryBuilder{}
	b.Join(sb, q.Query.Joins)

	// WHERE
	if len(*q.Query.ConditionGroups) > 0 {
		wb := NewWhereBaseBuilder(q.Query.ConditionGroups)
		whereValues := wb.Where(sb, q.Query.ConditionGroups)
		values = append(values, whereValues...)
	}

	// ORDER BY
	if len(*q.Query.Order) > 0 {
		ob := NewOrderByBaseBuilder(q.Query.Order)
		ob.OrderBy(sb, q.Query.Order)
	}

	query := sb.String()
	sb.Reset()

	return query, values
}
