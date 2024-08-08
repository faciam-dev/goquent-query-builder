package base

import (
	"sort"

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

// Insert builds the INSERT query.
func (m InsertBaseBuilder) Insert(q *structs.InsertQuery) (string, []interface{}, error) {
	//sb := &strings.Builder{}
	//sb.Grow(consts.StringBuffer_Middle_Query_Grow) // todo: check if this is necessary
	sb := make([]byte, 0, consts.StringBuffer_Middle_Query_Grow)

	// INSERT INTO
	sb = append(sb, "INSERT INTO "...)
	sb = m.u.EscapeIdentifier2(sb, q.Table)
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
		sb = m.u.EscapeIdentifier2(sb, column)
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
	//query := sb.String()
	//sb.Reset()

	return query, values, nil
}

// InsertBatch builds the INSERT query for batch insert.
func (m InsertBaseBuilder) InsertBatch(q *structs.InsertQuery) (string, []interface{}, error) {
	//sb := &strings.Builder{}
	//sb.Grow(consts.StringBuffer_Long_Query_Grow) // todo: check if this is necessary

	sb := make([]byte, 0, consts.StringBuffer_Long_Query_Grow)

	// INSERT INTO

	sb = append(sb, "INSERT INTO "...)
	sb = m.u.EscapeIdentifier2(sb, q.Table)
	sb = append(sb, " "...)

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
	sb = append(sb, "("...)
	for i, column := range columns {
		if i > 0 {
			sb = append(sb, ", "...)
		}
		sb = m.u.EscapeIdentifier2(sb, column)
	}
	sb = append(sb, ") VALUES "...)

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

	//query := sb.String()
	//sb.Reset()

	return query, allValues, nil
}

func (m *InsertBaseBuilder) InsertUsing(q *structs.InsertQuery) (string, []interface{}, error) {
	//sb := &strings.Builder{}
	//sb.Grow(consts.StringBuffer_Middle_Query_Grow) // todo: check if this is necessary
	sb := make([]byte, 0, consts.StringBuffer_Middle_Query_Grow)

	// INSERT INTO
	sb = append(sb, "INSERT INTO "...)
	sb = m.u.EscapeIdentifier2(sb, q.Table)

	// COLUMNS
	columns := make([]string, 0, len(q.Columns))
	columns = append(columns, q.Columns...)
	sb = append(sb, " ("...)
	for i, column := range columns {
		if i > 0 {
			sb = append(sb, ", "...)
		}
		sb = m.u.EscapeIdentifier2(sb, column)
	}
	sb = append(sb, ") "...)

	// SELECT
	//b := m.u.GetQueryBuilderStrategy()
	selectValues := []interface{}{}
	//selectValues := b.Build(sb, q.Query, 0, nil)
	//sb = append(sb, selectQuery...)

	query := string(sb)
	//query := sb.String()
	//sb.Reset()

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
