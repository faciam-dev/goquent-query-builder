package base

import (
	"sort"
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/jsonutils"
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
	sb = m.u.EscapeIdentifier(sb, q.Table)

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
		if strings.Contains(column, "->") {
			field, path := jsonutils.ParseJsonFieldAndPath(column)
			switch m.u.Dialect() {
			case consts.DialectMySQL:
				sb = m.u.EscapeIdentifier(sb, field)
				sb = append(sb, " = JSON_SET("...)
				sb = m.u.EscapeIdentifier(sb, field)
				sb = append(sb, ", '$."+strings.Join(path, ".")+"', "...)
				sb = append(sb, m.u.GetPlaceholder()...)
				sb = append(sb, ')')
			case consts.DialectPostgreSQL:
				sb = m.u.EscapeIdentifier(sb, field)
				sb = append(sb, " = jsonb_set("...)
				sb = m.u.EscapeIdentifier(sb, field)
				sb = append(sb, ", '{"+strings.Join(path, ",")+"}', "...)
				sb = append(sb, m.u.GetPlaceholder()...)
				sb = append(sb, ')')
			default:
				sb = m.u.EscapeIdentifier(sb, column)
				sb = append(sb, " = "+m.u.GetPlaceholder()...)
			}
		} else {
			sb = m.u.EscapeIdentifier(sb, column)
			sb = append(sb, " = "+m.u.GetPlaceholder()...)
		}
		if i < len(columns)-1 {
			sb = append(sb, ", "...)
		}
		values = append(values, q.Values[column])
	}

	// WHERE
	if len(q.Query.ConditionGroups) > 0 {
		wb := NewWhereBaseBuilder(m.u, q.Query.ConditionGroups)
		whereValues, err := wb.Where(&sb, q.Query.ConditionGroups)
		if err != nil {
			return "", nil, err
		}
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
