package db

import (
	"log"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type WhereBaseBuilder struct {
	whereGroups *[]structs.WhereGroup
}

func NewWhereBaseBuilder(wg *[]structs.WhereGroup) *WhereBaseBuilder {
	return &WhereBaseBuilder{
		whereGroups: wg,
	}
}

func (wb *WhereBaseBuilder) Where(wg *[]structs.WhereGroup) (string, []interface{}) {
	// WHERE
	where := ""
	values := []interface{}{}

	//log.Default().Printf("wherewhere: %v", wherewhere)
	sep := ""
	for i, cg := range *wg {
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
				log.Default().Printf("c.Query: %v", *c.Query.ConditionGroups)
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
