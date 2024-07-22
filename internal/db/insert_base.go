package db

import (
	"sort"
	"strings"

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
func (m InsertBaseBuilder) Insert(q *structs.InsertQuery) (string, []interface{}) {
	// INSERT INTO
	query := "INSERT INTO " + q.Table + " "

	columns := make([]string, 0, len(q.Values))
	for column := range q.Values {
		columns = append(columns, column)
	}
	sort.Strings(columns)

	placeholders := make([]string, len(columns))
	values := make([]interface{}, 0, len(columns))
	for i, column := range columns {
		placeholders[i] = "?"
		values = append(values, q.Values[column])
	}

	query += "(" + strings.Join(columns, ", ") + ") "

	query += "VALUES (" + strings.Join(placeholders, ", ") + ")"

	return query, values
}

// InsertBatch builds the INSERT query for batch insert.
func (m InsertBaseBuilder) InsertBatch(q *structs.InsertQuery) (string, []interface{}) {
	// INSERT INTO
	query := "INSERT INTO " + q.Table + " "

	// 全てのレコードからカラム名を収集し、重複を排除
	columnSet := make(map[string]struct{})
	for _, values := range q.ValuesBatch {
		for column := range values {
			columnSet[column] = struct{}{}
		}
	}

	// カラム名をスライスに格納し、ソート
	columns := make([]string, 0, len(columnSet))
	for column := range columnSet {
		columns = append(columns, column)
	}
	sort.Strings(columns)

	// COLUMNS
	query += "(" + strings.Join(columns, ", ") + ") VALUES "

	// VALUES
	valuePlaceholders := make([]string, 0, len(q.ValuesBatch))
	var allValues []interface{}
	for _, values := range q.ValuesBatch {
		placeholders := make([]string, len(columns))
		rowValues := make([]interface{}, 0, len(columns))
		for i, column := range columns {
			placeholders[i] = "?"
			if value, ok := values[column]; ok {
				rowValues = append(rowValues, value)
			} else {
				// カラムが存在しない場合はNULLを挿入
				rowValues = append(rowValues, nil)
			}
		}
		valuePlaceholders = append(valuePlaceholders, "("+strings.Join(placeholders, ", ")+")")
		allValues = append(allValues, rowValues...)
	}

	query += strings.Join(valuePlaceholders, ", ")

	return query, allValues
}

func (m InsertBaseBuilder) InsertUsing(q *structs.InsertQuery) (string, []interface{}) {
	// INSERT INTO
	query := "INSERT INTO " + q.Table

	// COLUMNS
	columns := make([]string, 0, len(q.Columns))
	columns = append(columns, q.Columns...)
	query += " (" + strings.Join(columns, ", ") + ") "

	// SELECT
	b := &BaseQueryBuilder{}
	selectQuery, selectValues := b.Build("", q.Query)
	query += selectQuery

	return query, selectValues
}

// BuildInsert builds the INSERT query.
func (m InsertBaseBuilder) BuildInsert(q *structs.InsertQuery) (string, []interface{}) {
	if q.Query != nil {
		return m.InsertUsing(q)
	}

	if len(q.Values) > 0 {
		return m.Insert(q)
	}

	return m.InsertBatch(q)
}
