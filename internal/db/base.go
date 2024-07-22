package db

import (
	"fmt"
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/sliceutils"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type BaseQueryBuilder struct {
	WhereBaseBuilder
	OrderByBaseBuilder
	JoinBaseBuilder
	InsertBaseBuilder
	UpdateBaseBuilder
	DeleteBaseBuilder
}

func (b *BaseQueryBuilder) Select(sb *strings.Builder, columns *[]structs.Column, tableName string, joins *structs.Joins) []interface{} {
	if columns == nil {
		sb.WriteString("SELECT * ")
		return []interface{}{}
	}
	//colNames := make([]string, 0, len(*columns))

	// if there are no columns to select, select all columns
	if len(*columns) == 0 && joins.Joins != nil {
		for i, join := range *joins.Joins {
			b.processJoin(sb, &join, tableName, i)
		}

		if joins.JoinClause != nil {
			join := structs.Join{
				TargetNameMap: joins.JoinClause.TargetNameMap,
				Name:          joins.JoinClause.Name,
			}
			b.processJoin(sb, &join, tableName, 0)
		}

		return []interface{}{}
	}

	colValues := make([]interface{}, 0, len(*columns))
	// if there are columns to select
	for i, column := range *columns {
		if column.Raw != "" {
			if column.Values != nil && len(column.Values) > 0 {
				colValues = append(colValues, column.Values...)
			}
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(column.Raw) // or colNames = column.Raw
			//colNames = append(colNames, column.Raw)
		} else if column.Name != "" {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(column.Name)
			//colNames = append(colNames, column.Name)
		}
	}

	return colValues
}

func (j *BaseQueryBuilder) processJoin(sb *strings.Builder, join *structs.Join, tableName string, idx int) {
	targetName := ""
	//joinedTablesForSelect := ""

	if _, ok := join.TargetNameMap[consts.Join_CROSS]; ok {
		targetName = join.TargetNameMap[consts.Join_CROSS]
	}
	if _, ok := join.TargetNameMap[consts.Join_RIGHT]; ok {
		targetName = join.TargetNameMap[consts.Join_RIGHT]
	}
	if _, ok := join.TargetNameMap[consts.Join_LEFT]; ok {
		targetName = join.TargetNameMap[consts.Join_LEFT]
	}
	if _, ok := join.TargetNameMap[consts.Join_INNER]; ok {
		targetName = join.TargetNameMap[consts.Join_INNER]
	}

	if targetName == "" {
		return
	}

	name := tableName
	if join.Name != "" {
		name = join.Name
	}

	targetNameForSelect := targetName + ".*"

	//sb.Grow(consts.StringBuffer_Select_Grow)

	outputed := false
	if !sliceutils.Contains(*j.columnNames, targetNameForSelect) {
		if idx > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(targetNameForSelect)
		*j.columnNames = append(*j.columnNames, targetNameForSelect)
		outputed = true
	}

	nameForSelect := name + ".*"
	if !sliceutils.Contains(*j.columnNames, nameForSelect) {
		if idx > 0 || outputed {
			sb.WriteString(", ")
		}
		sb.WriteString(nameForSelect)
		*j.columnNames = append(*j.columnNames, nameForSelect)
	}

}

func (BaseQueryBuilder) From(sb *strings.Builder, table string) {
	sb.WriteString("FROM ")
	sb.WriteString(table)
}

func (BaseQueryBuilder) GroupBy(sb *strings.Builder, groupBy *structs.GroupBy) []interface{} {
	if groupBy == nil || len(groupBy.Columns) == 0 {
		return []interface{}{}
	}

	groupByColumns := groupBy.Columns
	if len(groupByColumns) > 0 {
		sb.WriteString(" GROUP BY ")
		for i, column := range groupByColumns {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(column)
		}
	}

	values := make([]interface{}, 0, len(*groupBy.Having))

	if len(*groupBy.Having) > 0 {
		sb.WriteString(" HAVING ")

		//havingValues := make([]interface{}, 0, len(*groupBy.Having))
		for n, having := range *groupBy.Having {
			op := "AND"
			if having.Operator == consts.LogicalOperator_AND {
				op = "AND"
			} else if having.Operator == consts.LogicalOperator_OR {
				op = "OR"
			}

			if having.Raw != "" {
				if n > 0 {
					sb.WriteString(" ")
					sb.WriteString(op)
					sb.WriteString(" ")
				}
				sb.WriteString(having.Raw)
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
			//havingValues = append(havingValues, having.Value)
			values = append(values, having.Value)

			if n > 0 {
				sb.WriteString(" ")
				sb.WriteString(op)
				sb.WriteString(" ")
			}
			sb.WriteString(having.Column)
			sb.WriteString(" ")
			sb.WriteString(having.Condition)
			sb.WriteString(" ?")
		}

		//if len(havingValues) > 0 {
		//	values = append(values, havingValues...)
		//}
	}

	return values
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
func (m BaseQueryBuilder) Build(cacheKey string, q *structs.Query) (string, []interface{}) {
	sb := &strings.Builder{}

	// grow the string builder based on the length of the cache key
	if len(cacheKey) < consts.StringBuffer_Short_Query_Grow {
		sb.Grow(consts.StringBuffer_Short_Query_Grow)
	} else if len(cacheKey) < consts.StringBuffer_Middle_Query_Grow {
		sb.Grow(consts.StringBuffer_Middle_Query_Grow)
	} else {
		sb.Grow(consts.StringBuffer_Long_Query_Grow)
	}

	// JOIN
	//sb.Grow(consts.StringBuffer_Query_Grow)

	// assemble the query
	// SELECT AND FROM

	sb.WriteString("SELECT ")

	// SELECT
	colValues := m.Select(sb, q.Columns, q.Table.Name, q.Joins)

	sb.WriteString(" ")
	m.From(sb, q.Table.Name)
	values := colValues

	// JOIN
	joinValues := m.Join(sb, q.Joins)
	//sb.WriteString(join)
	values = append(values, joinValues...)

	// WHERE
	whereValues := m.Where(sb, q.ConditionGroups)
	//sb.WriteString(where)
	values = append(values, whereValues...)

	// ORDER BY
	m.OrderBy(sb, q.Order)
	//sb.WriteString(orderBy)

	// GROUP BY / HAVING
	groupByValues := m.GroupBy(sb, q.Group)
	//sb.WriteString(groupBy)
	values = append(values, groupByValues...)

	// LIMIT
	limit := m.Limit(q.Limit)
	sb.WriteString(limit)

	// OFFSET
	offset := m.Offset(q.Offset)
	sb.WriteString(offset)

	// LOCK
	lock := m.Lock(q.Lock)
	sb.WriteString(lock)

	query := sb.String()
	sb.Reset()

	return query, values
}
