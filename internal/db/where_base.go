package db

import (
	"strings"

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

func (wb *WhereBaseBuilder) Where(sb *strings.Builder, wg *[]structs.WhereGroup) []interface{} {
	if wg == nil || len(*wg) == 0 {
		return []interface{}{}
	}

	// WHERE
	hasCondition := false
	for _, cg := range *wg {
		if len(cg.Conditions) > 0 {
			hasCondition = true
			break
		}
	}

	if hasCondition {
		sb.WriteString(" WHERE ")
	}

	values := make([]interface{}, 0)

	sep := ""
	notSep := ""
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

		if i > 0 {
			sb.WriteString(sep)
		}

		parenthesesOpen := ""
		parenthesesClose := ""
		op := ""

		if !cg.IsDummyGroup {
			if cg.IsNot {
				notSep = "NOT "
			} else {
				notSep = ""
			}
			parenthesesOpen = "("
			parenthesesClose = ")"
		}

		sb.WriteString(notSep)
		sb.WriteString(parenthesesOpen)

		for j, c := range cg.Conditions {
			convertedColumn := c.Column
			if i > 0 && j == 0 && cg.IsDummyGroup {
				if c.Operator == consts.LogicalOperator_AND {
					op = " AND "
				} else if c.Operator == consts.LogicalOperator_OR {
					op = " OR "
				}
			}

			if c.Query != nil {
				condQuery := convertedColumn + " " + c.Condition

				// create subquery
				b := &BaseQueryBuilder{}
				sqQuery, sqValues := b.Build("", c.Query)

				if c.Operator == consts.LogicalOperator_AND {
					if op != "" {
						op = " AND "
					}
					sb.WriteString(op + condQuery + " (" + sqQuery + ")")
					if op == "" {
						op = " AND "
					}
				} else if c.Operator == consts.LogicalOperator_OR {
					if op != "" {
						op = " OR "
					}
					sb.WriteString(op + condQuery + " (" + sqQuery + ")")
					if op == "" {
						op = " OR "
					}
				}

				values = append(values, sqValues...)
			} else if c.Exists != nil {
				condQuery := c.Condition

				// create subquery
				b := &BaseQueryBuilder{}
				sqQuery, sqValues := b.Build("", c.Exists.Query)

				if c.Operator == consts.LogicalOperator_AND {
					if op != "" {
						op = " AND "
					}
					sb.WriteString(op + condQuery + " (" + sqQuery + ")")
					if op == "" {
						op = " AND "
					}
				} else if c.Operator == consts.LogicalOperator_OR {
					if op != "" {
						op = " OR "
					}
					sb.WriteString(op + condQuery + " (" + sqQuery + ")")
					if op == "" {
						op = " OR "
					}
				}

				values = append(values, sqValues...)

			} else if c.Between != nil {
				wsb := strings.Builder{}
				wsb.Grow(consts.StringBuffer_Where_Grow)
				if c.Between.IsColumn {
					wsb.WriteString(convertedColumn + " " + c.Condition + " " + c.Between.From.(string) + " AND " + c.Between.To.(string))
				} else {
					wsb.WriteString(convertedColumn + " " + c.Condition + " ? AND ?")
					values = append(values, c.Between.From, c.Between.To)
				}
				condQuery := wsb.String()
				wsb.Reset()

				if c.Operator == consts.LogicalOperator_AND {
					if op != "" {
						op = " AND "
					}
					sb.WriteString(op + condQuery)
					if op == "" {
						op = " AND "
					}
				} else if c.Operator == consts.LogicalOperator_OR {
					if op != "" {
						op = " OR "
					}
					sb.WriteString(op + condQuery)
					if op == "" {
						op = " OR "
					}
				}
			} else {
				raw := c.Raw
				wsb := strings.Builder{}
				wsb.Grow(consts.StringBuffer_Where_Grow)
				if raw != "" {
					wsb.WriteString(raw)
				} else {
					wsb.WriteString(convertedColumn + " " + c.Condition)
					if c.ValueColumn != "" {
						wsb.WriteString(" " + c.ValueColumn)
					} else if c.Value != nil {
						if len(c.Value) > 1 {
							wsb.WriteString(" (")
							for k := 0; k < len(c.Value); k++ {
								if k > 0 {
									wsb.WriteString(", ")
								}
								wsb.WriteString("?")
							}
							wsb.WriteString(")")
						} else {
							wsb.WriteString(" ?")
						}
					}
				}
				condQuery := wsb.String()
				wsb.Reset()

				if c.Operator == consts.LogicalOperator_AND {
					if op != "" {
						op = " AND "
					}
					sb.WriteString(op + condQuery)
					if len(c.Value) > 0 {
						values = append(values, c.Value...)
					}
					if op == "" {
						op = " AND "
					}
				} else if c.Operator == consts.LogicalOperator_OR {
					if op != "" {
						op = " OR "
					}
					sb.WriteString(op)
					sb.WriteString(condQuery)
					if len(c.Value) > 0 {
						values = append(values, c.Value...)
					}
					if op == "" {
						op = " OR "
					}
				}
			}
		}
		sb.WriteString(parenthesesClose)
	}

	return values
}
