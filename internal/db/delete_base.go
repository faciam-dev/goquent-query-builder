package db

import (
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

	// JOIN
	b := &BaseQueryBuilder{}
	_, join, _ := b.Join(q.Table, q.Query.Joins)

	// DELETE
	query := "DELETE"
	if join != "" {
		query += " " + q.Table
	}

	// FROM
	query += " FROM " + q.Table

	// JOIN
	query += join

	// WHERE
	if len(*q.Query.ConditionGroups) > 0 {
		wb := NewWhereBaseBuilder(q.Query.ConditionGroups)
		where, whereValues := wb.Where(q.Query.ConditionGroups)
		query += where
		values = append(values, whereValues...)
	}

	// ORDER BY
	if len(*q.Query.Order) > 0 {
		ob := NewOrderByBaseBuilder(q.Query.Order)
		order := ob.OrderBy(q.Query.Order)
		query += order
	}

	return query, values
}
