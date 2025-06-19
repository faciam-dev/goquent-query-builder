package base

import (
	"sort"
	"strings"
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

func (m InsertBaseBuilder) InsertIgnore(q *structs.InsertQuery) (string, []interface{}, error) {
	query, values, err := m.Insert(q)
	if err != nil {
		return "", nil, err
	}

	if m.u.Dialect() == consts.DialectMySQL {
		query = strings.Replace(query, "INSERT INTO", "INSERT IGNORE INTO", 1)
	} else if m.u.Dialect() == consts.DialectPostgreSQL {
		query += " ON CONFLICT DO NOTHING"
	}

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
	selectValues, err := b.Build(&sb, q.Query, 0, nil)
	if err != nil {
		return "", nil, err
	}

	query := string(sb)

	*ptr = sb
	poolBytes.Put(ptr)

	return query, selectValues, nil
}

func (m InsertBaseBuilder) Upsert(q *structs.InsertQuery) (string, []interface{}, error) {
	// ensure ValuesBatch is set
	if len(q.ValuesBatch) == 0 && len(q.Values) > 0 {
		q.ValuesBatch = []map[string]interface{}{q.Values}
	}

	baseQuery, values, err := m.InsertBatch(q)
	if err != nil {
		return "", nil, err
	}

	sb := []byte(baseQuery)

	if m.u.Dialect() == consts.DialectMySQL {
		sb = append(sb, []byte(" ON DUPLICATE KEY UPDATE ")...)
		for i, col := range q.Upsert.UpdateColumns {
			if i > 0 {
				sb = append(sb, []byte(", ")...)
			}
			sb = m.u.EscapeIdentifier(sb, col)
			sb = append(sb, []byte(" = VALUES(")...)
			sb = m.u.EscapeIdentifier(sb, col)
			sb = append(sb, []byte(")")...)
		}
	} else if m.u.Dialect() == consts.DialectPostgreSQL {
		sb = append(sb, []byte(" ON CONFLICT (")...)
		for i, col := range q.Upsert.UniqueColumns {
			if i > 0 {
				sb = append(sb, []byte(", ")...)
			}
			sb = m.u.EscapeIdentifier(sb, col)
		}
		sb = append(sb, []byte(") DO UPDATE SET ")...)
		for i, col := range q.Upsert.UpdateColumns {
			if i > 0 {
				sb = append(sb, []byte(", ")...)
			}
			sb = m.u.EscapeIdentifier(sb, col)
			sb = append(sb, []byte(" = EXCLUDED.")...)
			sb = m.u.EscapeIdentifier(sb, col)
		}
	}

	return string(sb), values, nil
}

// BuildInsert builds the INSERT query.
func (m InsertBaseBuilder) BuildInsert(q *structs.InsertQuery) (string, []interface{}, error) {
	if q.Upsert != nil {
		return m.Upsert(q)
	}

	if q.Ignore {
		if len(q.ValuesBatch) > 0 {
			// treat as single batch but with ignore
			query, values, err := m.InsertBatch(q)
			if err != nil {
				return "", nil, err
			}
			if m.u.Dialect() == consts.DialectMySQL {
				query = strings.Replace(query, "INSERT INTO", "INSERT IGNORE INTO", 1)
			} else if m.u.Dialect() == consts.DialectPostgreSQL {
				query += " ON CONFLICT DO NOTHING"
			}
			return query, values, nil
		}
		return m.InsertIgnore(q)
	}

	if q.Query != nil {
		return m.InsertUsing(q)
	}

	if len(q.Values) > 0 && len(q.ValuesBatch) == 0 {
		return m.Insert(q)
	}

	return m.InsertBatch(q)
}
