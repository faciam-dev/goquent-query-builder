package db

import (
	"fmt"
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/sliceutils"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type BaseQueryBuilder struct{}

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
				colValues = append(colValues, column.Values)
			}
			colNames = append(colNames, column.Raw)
		} else if column.Name != "" {
			colNames = append(colNames, column.Name)
		}
	}

	return colNames, colValues
}

func (BaseQueryBuilder) From(table string) string {
	return "FROM " + table
}

func (BaseQueryBuilder) Where(wg *[]structs.WhereGroup) (string, []interface{}) {
	return buildWhereClause(wg)
}

func (BaseQueryBuilder) Join(tableName string, joins *[]structs.Join) (*[]structs.Column, string) {
	join := ""

	joinedTablesForSelect, joinStrings := buildJoinStatement(tableName, joins)
	for _, joinString := range joinStrings {
		join += " " + joinString
	}

	return joinedTablesForSelect, join
}

func (BaseQueryBuilder) OrderBy(order *[]structs.Order) string {
	if len(*order) == 0 {
		return ""
	}

	orderBy := ""
	rawOrderQuerys := make([]string, 0, len(*order))
	orders := make([]string, 0, len(*order))
	for _, order := range *order {
		if order.Raw != "" {
			rawOrderQuerys = append(rawOrderQuerys, order.Raw)
			continue
		}
		if order.Column == "" {
			continue
		}
		desc := "DESC"
		if order.IsAsc {
			desc = "ASC"
		}
		orders = append(orders, order.Column+" "+desc)
	}
	orderByHeader := " ORDER BY "
	if len(rawOrderQuerys) > 0 {
		orderBy += orderByHeader + strings.Join(rawOrderQuerys, ", ")
		orderByHeader = ", "
	}
	if len(orders) > 0 {
		orderBy += orderByHeader + strings.Join(orders, ", ")
	}

	return orderBy

}

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

	return query, values
}

// buildJoinStatement builds the JOIN statement.
func buildJoinStatement(tableName string, joins *[]structs.Join) (*[]structs.Column, []string) {
	joinedTablesForSelect := make([]structs.Column, 0, len(*joins))
	joinStrings := make([]string, 0, len(*joins))
	for _, join := range *joins {
		joinType := ""
		targetName := ""

		if _, ok := join.TargetNameMap[consts.Join_CROSS]; ok {
			targetName = join.TargetNameMap[consts.Join_CROSS]
			joinType = consts.Join_Type_CROSS
		}
		if _, ok := join.TargetNameMap[consts.Join_RIGHT]; ok {
			targetName = join.TargetNameMap[consts.Join_RIGHT]
			joinType = consts.Join_Type_RIGHT
		}
		if _, ok := join.TargetNameMap[consts.Join_LEFT]; ok {
			targetName = join.TargetNameMap[consts.Join_LEFT]
			joinType = consts.Join_Type_LEFT
		}
		if _, ok := join.TargetNameMap[consts.Join_INNER]; ok {
			targetName = join.TargetNameMap[consts.Join_INNER]
			joinType = consts.Join_Type_INNER
		}

		if targetName == "" {
			continue
		}

		name := tableName
		if join.Name != "" {
			name = join.Name
		}

		targetNameForSelect := targetName + ".*"
		if !sliceutils.Contains[string](*getNowColNames(&joinedTablesForSelect), targetNameForSelect) {
			joinedTablesForSelect = append(joinedTablesForSelect, structs.Column{
				Name: targetNameForSelect,
			})
		}
		nameForSelect := name + ".*"
		if !sliceutils.Contains[string](*getNowColNames(&joinedTablesForSelect), nameForSelect) {
			joinedTablesForSelect = append(joinedTablesForSelect, structs.Column{
				Name: nameForSelect,
			})
		}

		joinQuery := joinType + " JOIN " + targetName + " ON " + join.SearchColumn + " " + join.SearchCondition + " " + join.SearchTargetColumn

		joinStrings = append(joinStrings, joinQuery)

	}

	return &joinedTablesForSelect, joinStrings
}

// getNowColNames returns the names of the columns in the slice.
func getNowColNames(joinedTablesForSelect *[]structs.Column) *[]string {
	nowColNames := make([]string, len(*joinedTablesForSelect))
	for _, joinedTable := range *joinedTablesForSelect {
		nowColNames = append(nowColNames, joinedTable.Name)
	}
	return &nowColNames
}

func buildWhereClause(conditionGroups *[]structs.WhereGroup) (string, []interface{}) {
	// WHERE
	where := ""
	values := []interface{}{}

	//log.Default().Printf("wherewhere: %v", wherewhere)
	sep := ""
	for i, cg := range *conditionGroups {
		if len(cg.Conditions) == 0 {
			continue
		}

		// AND, OR by ConditionGroup
		if cg.IsDummyGroup {
			sep = ""
		} else if cg.Operator == consts.LogicalOperator_AND {
			sep = " AND "
		} else if cg.Operator == consts.LogicalOperator_OR {
			sep = " OR "
		}

		if where != "" {
			where += sep
		}

		parenthesesOpen := ""
		parenthesesClose := ""
		op := ""

		if !cg.IsDummyGroup {
			parenthesesOpen = "("
			parenthesesClose = ")"
		}

		where += parenthesesOpen

		for j, c := range cg.Conditions {
			convertedColumn := c.Column
			if i > 0 && j == 0 && cg.IsDummyGroup {
				if c.Operator == consts.LogicalOperator_AND {
					op = " AND "
				} else if c.Operator == consts.LogicalOperator_OR {
					op = " OR "
				}
			}
			/*
				convertedColumn := c.Colmun
				convertedSelectColumns := []structs.Column{}
				if c.Query.Columns != nil {
					for _, column := range *c.Query.Columns {
						convertedSelectColumn := column
						if column.Raw != "" {
							convertedSelectColumn.Raw = column.Raw
						}
						convertedSelectColumns = append(convertedSelectColumns, convertedSelectColumn)
					}
				}
			*/
			if c.Query != nil { // && c.Query.ConditionGroups != nil && len(*c.Query.ConditionGroups) > 0) || (c.Query != nil && c.Query.SubQuery != nil && len(*c.Query.SubQuery) > 0) {
				condQuery := convertedColumn + " " + c.Condition

				// create subquery
				b := &BaseQueryBuilder{}
				//log.Default().Printf("c.Query.Pro: %v", c.Query.Processed)
				sqQuery, sqValues := b.Build(c.Query)
				if c.Operator == consts.LogicalOperator_AND {
					if op != "" {
						op = " AND "
					}
					where += op + condQuery + " (" + sqQuery + ")"
					if op == "" {
						op = " AND "
					}
				} else if c.Operator == consts.LogicalOperator_OR {
					if op != "" {
						op = " OR "
					}
					where += op + condQuery + " (" + sqQuery + ")"
					if op == "" {
						op = " OR "
					}
				}
				values = append(values, sqValues...)
			} else if len(c.Value) == 1 {
				condQuery := convertedColumn + " " + c.Condition + " ?"
				value := c.Value[0]
				if c.Operator == consts.LogicalOperator_AND {
					if op != "" {
						op = " AND "
					}
					where += op + condQuery
					values = append(values, value)
					if op == "" {
						op = " AND "
					}
				} else if c.Operator == consts.LogicalOperator_OR {
					if op != "" {
						op = " OR "
					}
					where += op + condQuery
					values = append(values, value)
					if op == "" {
						op = " OR "
					}
				}
			} else {
				condQuery := convertedColumn + " " + c.Condition + " (?)"
				value := c.Value
				if c.Operator == consts.LogicalOperator_AND {
					if op != "" {
						op = " AND "
					}
					where += op + condQuery
					values = append(values, value)
					if op == "" {
						op = " AND "
					}
				} else if c.Operator == consts.LogicalOperator_OR {
					if op != "" {
						op = " OR "
					}
					where += op + condQuery
					values = append(values, value)
					if op == "" {
						op = " OR "
					}
				}
			}
		}
		where += parenthesesClose
	}

	if where != "" {
		where = " WHERE " + where
	}

	return where, values
}
