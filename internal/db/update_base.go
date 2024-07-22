package db

import (
	"sort"
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type UpdateBaseBuilder struct {
}

func NewUpdateBaseBuilder(iq *structs.UpdateQuery) *UpdateBaseBuilder {
	return &UpdateBaseBuilder{}
}

func (m *UpdateBaseBuilder) Update(q *structs.UpdateQuery) *UpdateBaseBuilder {
	return m
}

// UpdateBatch builds the Update query for Update.
func (m *UpdateBaseBuilder) BuildUpdate(q *structs.UpdateQuery) (string, []interface{}) {
	sb := &strings.Builder{}
	sb.Grow(consts.StringBuffer_Update_Grow)

	// UPDATE
	query := "UPDATE " + q.Table

	// JOIN
	b := &BaseQueryBuilder{}
	joinValues := b.Join(sb, q.Query.Joins)

	query += " SET "

	values := make([]interface{}, 0, len(q.Values)+len(joinValues))
	columns := make([]string, 0, len(q.Values))

	// JOIN
	values = append(values, joinValues...)

	// SET
	for column := range q.Values {
		columns = append(columns, column)
	}
	sort.Strings(columns)
	for _, column := range columns {
		query += column + " = ?, "
		values = append(values, q.Values[column])
	}
	query = query[:len(query)-2]

	// WHERE
	if len(*q.Query.ConditionGroups) > 0 {
		wb := NewWhereBaseBuilder(q.Query.ConditionGroups)
		whereValues := wb.Where(sb, q.Query.ConditionGroups)
		//query += where
		values = append(values, whereValues...)
	}

	if len(*q.Query.Order) > 0 {
		ob := NewOrderByBaseBuilder(q.Query.Order)
		ob.OrderBy(sb, q.Query.Order)
		//query += order
	}

	return query, values

}
