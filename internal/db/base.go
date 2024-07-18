package db

import (
	"fmt"
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type BaseQueryBuilder struct {
	WhereBaseBuilder
	OrderByBaseBuilder
	JoinBaseBuilder
	InsertBaseBuilder
}

func (BaseQueryBuilder) Select(columns *[]structs.Column, joinedTablesForSelect *[]structs.Column) ([]string, []interface{}) {
	colNames := make([]string, 0, len(*columns))
	colValues := make([]interface{}, 0, len(*columns))
	selectColumns := columns
	if len(*selectColumns) == 0 && len(*joinedTablesForSelect) > 0 {
		selectColumns = joinedTablesForSelect
	}
	for _, column := range *selectColumns {
		if column.Raw != "" {
			if column.Values != nil && len(column.Values) > 0 {
				colValues = append(colValues, column.Values...)
			}
			colNames = append(colNames, column.Raw) // or colNames = column.Raw
		} else if column.Name != "" {
			colNames = append(colNames, column.Name)
		}
	}

	return colNames, colValues
}

func (BaseQueryBuilder) From(table string) string {
	return "FROM " + table
}

func (BaseQueryBuilder) GroupBy(groupBy *structs.GroupBy) (string, []interface{}) {
	if groupBy == nil || len(groupBy.Columns) == 0 {
		return "", []interface{}{}
	}

	query := " GROUP BY "
	values := []interface{}{}
	groupByColumns := groupBy.Columns
	if len(groupByColumns) > 0 {
		query += strings.Join(groupByColumns, ", ")
	}

	if len(*groupBy.Having) > 0 {
		havingRaw := ""
		havingValues := []interface{}{}
		for n, having := range *groupBy.Having {
			op := "AND"
			if having.Operator == consts.LogicalOperator_AND {
				op = "AND"
			} else if having.Operator == consts.LogicalOperator_OR {
				op = "OR"
			}

			if having.Raw != "" {
				if n > 0 {
					havingRaw += " " + op + " "
				}
				havingRaw += having.Raw
				continue
			}
			if having.Column == "" {
				continue
			}
			if having.Condition == "" {
				continue
			}
			if having.Value == "" {
				continue
			}
			havingValues = append(havingValues, having.Value)

			if n > 0 {
				havingRaw += " " + op + " "
			}
			havingRaw += having.Column + " " + having.Condition + " ?"

		}

		if havingRaw != "" {
			query += " HAVING " + havingRaw

			if len(havingValues) > 0 {
				values = append(values, havingValues...)
			}
		}
	}

	return query, values
}

func (BaseQueryBuilder) Limit(limit *structs.Limit) string {
	if limit == nil || limit.Limit == 0 {
		return ""
	}

	return " LIMIT " + fmt.Sprint(limit.Limit)
}

func (BaseQueryBuilder) Offset(offset *structs.Offset) string {
	if offset == nil || offset.Offset == 0 {
		return ""
	}

	return " OFFSET " + fmt.Sprint(offset.Offset)
}

// Lock returns the lock statement.
func (BaseQueryBuilder) Lock(lock *structs.Lock) string {
	if lock == nil || lock.LockType == "" {
		return ""
	}

	return " " + lock.LockType
}

// Build builds the query.
func (m BaseQueryBuilder) Build(q *structs.Query) (string, []interface{}) {
	// JOIN
	joinedTablesForSelect, join := m.Join(q.Table.Name, q.Joins)

	// WHERE
	where, whereValues := m.Where(q.ConditionGroups)

	// ORDER BY
	orderBy := m.OrderBy(q.Order)

	// SELECT
	columns, colValues := m.Select(q.Columns, joinedTablesForSelect)

	// assemble the query
	// SELECT AND FROM
	query := fmt.Sprintf("SELECT %s %s", strings.Join(columns, ", "), m.From(q.Table.Name))
	values := colValues

	// JOIN
	query += join

	// WHERE
	query += where
	values = append(values, whereValues...)

	// ORDER BY
	query += orderBy

	// GROUP BY / HAVING
	groupBy, groupByValues := m.GroupBy(q.Group)
	query += groupBy
	values = append(values, groupByValues...)

	// LIMIT
	limit := m.Limit(q.Limit)
	query += limit

	// OFFSET
	offset := m.Offset(q.Offset)
	query += offset

	// LOCK
	lock := m.Lock(q.Lock)
	query += lock

	return query, values
}
