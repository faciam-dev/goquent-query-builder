package base

import (
	"sort"

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
	ptr := poolBytes.Get().(*[]byte)
	sb := *ptr
	if len(sb) > 0 {
		sb = sb[:0]
	}

	vPtr := poolValues.Get().(*[]interface{})
	values := *vPtr
	if len(values) > 0 {
		values = values[0:0]
	}

	// UPDATE
	sb = append(sb, "UPDATE "...)
	sb = m.u.EscapeIdentifier2(sb, q.Table)

	// JOIN
	b := NewJoinBaseBuilder(m.u, q.Query.Joins)
	joinValues := b.Join(&sb, q.Query.Joins)
	values = append(values, joinValues...)

	// SET
	sb = append(sb, " SET "...)
	columns := make([]string, 0, len(q.Values))
	for column := range q.Values {
		columns = append(columns, column)
	}
	sort.Strings(columns)
	for i, column := range columns {
		sb = m.u.EscapeIdentifier2(sb, column)
		sb = append(sb, " = "+m.u.GetPlaceholder()...)
		if i < len(columns)-1 {
			sb = append(sb, ", "...)
		}
		values = append(values, q.Values[column])
	}

	// WHERE
	if len(q.Query.ConditionGroups) > 0 {
		wb := NewWhereBaseBuilder(m.u, q.Query.ConditionGroups)
		whereValues := wb.Where(&sb, q.Query.ConditionGroups)
		values = append(values, whereValues...)
	}

	if len(*q.Query.Order) > 0 {
		ob := NewOrderByBaseBuilder(m.u, q.Query.Order)
		ob.OrderBy(&sb, q.Query.Order)
	}

	query := string(sb)

	*ptr = sb
	poolBytes.Put(ptr)

	*vPtr = values
	poolValues.Put(vPtr)

	return query, values, nil
}
