package base

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
func (m *UpdateBaseBuilder) BuildUpdate(q *structs.UpdateQuery) (string, []interface{}, error) {
	sb := &strings.Builder{}
	sb.Grow(consts.StringBuffer_Middle_Query_Grow) // todo: check if this is necessary

	// UPDATE
	sb.WriteString("UPDATE ")
	sb.WriteString(q.Table)

	// JOIN
	b := &BaseQueryBuilder{}
	joinValues := b.Join(sb, q.Query.Joins)
	values := make([]interface{}, 0, len(q.Values)+len(joinValues))

	values = append(values, joinValues...)

	// SET
	sb.WriteString(" SET ")
	columns := make([]string, 0, len(q.Values))
	for column := range q.Values {
		columns = append(columns, column)
	}
	sort.Strings(columns)
	for i, column := range columns {
		sb.WriteString(column)
		sb.WriteString(" = ?")
		if i < len(columns)-1 {
			sb.WriteString(", ")
		}
		values = append(values, q.Values[column])
	}

	// WHERE
	if len(*q.Query.ConditionGroups) > 0 {
		wb := NewWhereBaseBuilder(q.Query.ConditionGroups)
		whereValues := wb.Where(sb, q.Query.ConditionGroups)
		values = append(values, whereValues...)
	}

	if len(*q.Query.Order) > 0 {
		ob := NewOrderByBaseBuilder(q.Query.Order)
		ob.OrderBy(sb, q.Query.Order)
	}

	query := sb.String()
	sb.Reset()

	return query, values, nil
}
