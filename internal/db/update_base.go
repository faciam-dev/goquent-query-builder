package db

import (
	"sort"

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
	// JOIN
	b := &BaseQueryBuilder{}
	_, join := b.Join(q.Table, q.Query.Joins)

	// UPDATE
	query := "UPDATE " + q.Table + join + " SET "

	values := make([]interface{}, 0, len(q.Values))
	columns := make([]string, 0, len(q.Values))

	for column := range q.Values {
		columns = append(columns, column)
	}
	sort.Strings(columns)
	for _, column := range columns {
		query += column + " = ?, "
		values = append(values, q.Values[column])
	}
	query = query[:len(query)-2]

	if len(*q.Query.ConditionGroups) > 0 {
		wb := NewWhereBaseBuilder(q.Query.ConditionGroups)
		where, whereValues := wb.Where(q.Query.ConditionGroups)
		query += where
		values = append(values, whereValues...)
	}

	if len(*q.Query.Order) > 0 {
		ob := NewOrderByBaseBuilder(q.Query.Order)
		order := ob.OrderBy(q.Query.Order)
		query += order
	}

	return query, values

}