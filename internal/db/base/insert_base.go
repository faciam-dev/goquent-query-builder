package base

import (
	"sort"
	"sync"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

type InsertBaseBuilder struct {
	u           interfaces.SQLUtils
	insertQuery *structs.InsertQuery
}

func NewInsertBaseBuilder(util interfaces.SQLUtils, iq *structs.InsertQuery) *InsertBaseBuilder {
	return &InsertBaseBuilder{
		u:           util,
		insertQuery: iq,
	}
}

var poolBytes = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 0, consts.StringBuffer_Long_Query_Grow)
		return &b
	},
}

var poolValues = sync.Pool{
	New: func() interface{} {
		i := make([]interface{}, 0)
		return &i
	},
}

// Insert builds the INSERT query.
func (m InsertBaseBuilder) Insert(q *structs.InsertQuery) (string, []interface{}, error) {
	ptr := poolBytes.Get().(*[]byte)
	sb := *ptr
	if len(sb) > 0 {
		sb = sb[:0]
	}

	// INSERT INTO
	sb = append(sb, "INSERT INTO "...)
	sb = m.u.EscapeIdentifier(sb, q.Table)
	sb = append(sb, " "...)

	columns := make([]string, 0, len(q.Values))
	for column := range q.Values {
		columns = append(columns, column)
	}
	sort.Strings(columns)

	values := make([]interface{}, 0, len(columns))
	for _, column := range columns {
		values = append(values, q.Values[column])
	}

	sb = append(sb, "("...)
	for i, column := range columns {
		if i > 0 {
			sb = append(sb, ", "...)
		}
		sb = m.u.EscapeIdentifier(sb, column)
	}
	sb = append(sb, ") "...)

	sb = append(sb, "VALUES ("...)
	for i := range columns {
		if i > 0 {
			sb = append(sb, ", "...)
		}
		sb = append(sb, m.u.GetPlaceholder()...)
	}
	sb = append(sb, ")"...)

	query := string(sb)

	*ptr = sb
	poolBytes.Put(ptr)

	return query, values, nil
}

// InsertBatch builds the INSERT query for batch insert.
func (m InsertBaseBuilder) InsertBatch(q *structs.InsertQuery) (string, []interface{}, error) {
	ptr := poolBytes.Get().(*[]byte)
	sb := *ptr
	if len(sb) > 0 {
		sb = sb[:0]
	}

	vPtr := poolValues.Get().(*[]interface{})
	allValues := *vPtr
	if len(allValues) > 0 {
		allValues = allValues[0:0]
	}

	// INSERT INTO
	sb = append(sb, "INSERT INTO "...)
	sb = m.u.EscapeIdentifier(sb, q.Table)
	sb = append(sb, " "...)

	// get all columns from all values
	columnSet := make(map[string]struct{}, len(q.ValuesBatch))
	for i := range q.ValuesBatch {
		for column := range q.ValuesBatch[i] {
			columnSet[column] = struct{}{}
		}
	}

	// sort columns
	columns := make([]string, 0, len(columnSet))
	for column := range columnSet {
		columns = append(columns, column)
	}
	sort.Strings(columns)

	// COLUMNS
	sb = append(sb, "("...)
	for i, column := range columns {
		if i > 0 {
			sb = append(sb, ", "...)
		}
		sb = m.u.EscapeIdentifier(sb, column)
	}
	sb = append(sb, ") VALUES "...)

	// VALUES
	estimatedSize := len(q.ValuesBatch) * len(columns)
	if cap(allValues) < estimatedSize {
		newAllValue := make([]interface{}, 0, estimatedSize)
		copy(newAllValue, allValues)
		allValues = newAllValue
	}
	for i, values := range q.ValuesBatch {
		rowValues := make([]interface{}, 0, len(columns))
		for j := range columns {
			if value, ok := values[columns[j]]; ok {
				rowValues = rowValues[:len(rowValues)+1]
				rowValues[j] = value
			} else {
				rowValues = rowValues[:len(rowValues)+1]
				rowValues[j] = nil
			}
		}

		sb = append(sb, "("...)
		for i := range columns {
			if i > 0 {
				sb = append(sb, ", "...)
			}
			sb = append(sb, m.u.GetPlaceholder()...)
		}
		sb = append(sb, ")"...)

		if i < len(q.ValuesBatch)-1 {
			sb = append(sb, ", "...)
		}

		allValues = append(allValues, rowValues...)
	}
	query := string(sb)

	*ptr = sb
	poolBytes.Put(ptr)

	*vPtr = allValues
	poolValues.Put(vPtr)

	return query, allValues, nil
}

func (m *InsertBaseBuilder) InsertUsing(q *structs.InsertQuery) (string, []interface{}, error) {
	ptr := poolBytes.Get().(*[]byte)
	sb := *ptr
	if len(sb) > 0 {
		sb = sb[:0]
	}

	// INSERT INTO
	sb = append(sb, "INSERT INTO "...)
	sb = m.u.EscapeIdentifier(sb, q.Table)

	// COLUMNS
	columns := make([]string, 0, len(q.Columns))
	columns = append(columns, q.Columns...)
	sb = append(sb, " ("...)
	for i, column := range columns {
		if i > 0 {
			sb = append(sb, ", "...)
		}
		sb = m.u.EscapeIdentifier(sb, column)
	}
	sb = append(sb, ") "...)

	// SELECT
	b := m.u.GetQueryBuilderStrategy()
	selectValues := b.Build(&sb, q.Query, 0, nil)

	query := string(sb)

	*ptr = sb
	poolBytes.Put(ptr)

	return query, selectValues, nil
}

// BuildInsert builds the INSERT query.
func (m InsertBaseBuilder) BuildInsert(q *structs.InsertQuery) (string, []interface{}, error) {
	if q.Query != nil {
		return m.InsertUsing(q)
	}

	if len(q.Values) > 0 {
		return m.Insert(q)
	}

	return m.InsertBatch(q)
}
