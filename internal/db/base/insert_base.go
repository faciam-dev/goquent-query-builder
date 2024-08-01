package base

import (
	"sort"
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type InsertBaseBuilder struct {
	insertQuery *structs.InsertQuery
}

func NewInsertBaseBuilder(iq *structs.InsertQuery) *InsertBaseBuilder {
	return &InsertBaseBuilder{
		insertQuery: iq,
	}
}

// Insert builds the INSERT query.
func (m InsertBaseBuilder) Insert(q *structs.InsertQuery) (string, []interface{}, error) {
	sb := &strings.Builder{}
	sb.Grow(consts.StringBuffer_Middle_Query_Grow) // todo: check if this is necessary

	// INSERT INTO
	sb.WriteString("INSERT INTO ")
	sb.WriteString(q.Table)
	sb.WriteString(" ")

	columns := make([]string, 0, len(q.Values))
	for column := range q.Values {
		columns = append(columns, column)
	}
	sort.Strings(columns)

	values := make([]interface{}, 0, len(columns))
	for _, column := range columns {
		values = append(values, q.Values[column])
	}

	sb.WriteString("(")
	for i, column := range columns {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(column)
	}
	sb.WriteString(") ")

	sb.WriteString("VALUES (")
	for i := range columns {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("?")
	}
	sb.WriteString(")")

	query := sb.String()
	sb.Reset()

	return query, values, nil
}

// InsertBatch builds the INSERT query for batch insert.
func (m InsertBaseBuilder) InsertBatch(q *structs.InsertQuery) (string, []interface{}, error) {
	sb := &strings.Builder{}
	sb.Grow(consts.StringBuffer_Long_Query_Grow) // todo: check if this is necessary

	// INSERT INTO
	sb.WriteString("INSERT INTO ")
	sb.WriteString(q.Table)
	sb.WriteString(" ")

	// get all columns from all values
	columnSet := make(map[string]struct{}, len(q.ValuesBatch))
	for _, values := range q.ValuesBatch {
		for column := range values {
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
	sb.WriteString("(")
	for i, column := range columns {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(column)
	}
	sb.WriteString(") VALUES ")

	// VALUES
	allValues := make([]interface{}, 0, len(q.ValuesBatch)*len(columns))
	for i, values := range q.ValuesBatch {
		rowValues := make([]interface{}, 0, len(columns))
		for _, column := range columns {
			if value, ok := values[column]; ok {
				rowValues = append(rowValues, value)
			} else {
				rowValues = append(rowValues, nil)
			}
		}

		sb.WriteString("(")
		for i := range columns {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString("?")
		}
		sb.WriteString(")")

		if i < len(q.ValuesBatch)-1 {
			sb.WriteString(", ")
		}

		allValues = append(allValues, rowValues...)
	}

	query := sb.String()
	sb.Reset()

	return query, allValues, nil
}

func (m InsertBaseBuilder) InsertUsing(q *structs.InsertQuery) (string, []interface{}, error) {
	sb := &strings.Builder{}
	sb.Grow(consts.StringBuffer_Middle_Query_Grow) // todo: check if this is necessary

	// INSERT INTO
	sb.WriteString("INSERT INTO ")
	sb.WriteString(q.Table)

	// COLUMNS
	columns := make([]string, 0, len(q.Columns))
	columns = append(columns, q.Columns...)
	sb.WriteString(" (")
	for i, column := range columns {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(column)
	}
	sb.WriteString(") ")

	// SELECT
	b := &BaseQueryBuilder{}
	selectQuery, selectValues := b.Build("", q.Query, 0, nil)
	sb.WriteString(selectQuery)

	query := sb.String()
	sb.Reset()

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
