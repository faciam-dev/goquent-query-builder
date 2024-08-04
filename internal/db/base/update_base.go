package base

import (
	"sort"
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

type UpdateBaseBuilder struct {
	u interfaces.SQLUtils
}

func NewUpdateBaseBuilder(util interfaces.SQLUtils, iq *structs.UpdateQuery) *UpdateBaseBuilder {
	return &UpdateBaseBuilder{
		u: util,
	}
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
	sb.WriteString(m.u.EscapeIdentifier(q.Table))

	// JOIN
	b := NewJoinBaseBuilder(m.u, q.Query.Joins)
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
		sb.WriteString(m.u.EscapeIdentifier(column))
		sb.WriteString(" = " + m.u.GetPlaceholder())
		if i < len(columns)-1 {
			sb.WriteString(", ")
		}
		values = append(values, q.Values[column])
	}

	// WHERE
	if len(*q.Query.ConditionGroups) > 0 {
		wb := NewWhereBaseBuilder(m.u, q.Query.ConditionGroups)
		whereValues := wb.Where(sb, q.Query.ConditionGroups)
		values = append(values, whereValues...)
	}

	if len(*q.Query.Order) > 0 {
		ob := NewOrderByBaseBuilder(m.u, q.Query.Order)
		ob.OrderBy(sb, q.Query.Order)
	}

	query := sb.String()
	sb.Reset()

	return query, values, nil
}
