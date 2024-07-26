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
	if hasCondition(*wg) {
		sb.WriteString(" WHERE ")
	}

	values := make([]interface{}, 0)

	for i, cg := range *wg {
		if len(cg.Conditions) == 0 {
			continue
		}

		if i > 0 {
			sb.WriteString(getConditionGroupSeparator(cg, i))
		}

		sb.WriteString(getNotSeparator(cg))
		sb.WriteString(getParenthesesOpen(cg))

		for j, c := range cg.Conditions {
			if j > 0 || (i > 0 && j == 0 && cg.IsDummyGroup) {
				sb.WriteString(getConditionOperator(c))
			}

			switch {
			case c.Query != nil:
				values = append(values, processSubQuery(sb, c)...)
			case c.Exists != nil:
				values = append(values, processExistsQuery(sb, c)...)
			case c.Between != nil:
				values = append(values, processBetweenCondition(sb, c)...)
			default:
				values = append(values, processRawCondition(sb, c)...)
			}
		}
		sb.WriteString(getParenthesesClose(cg))
	}

	return values
}

func hasCondition(wg []structs.WhereGroup) bool {
	for _, cg := range wg {
		if len(cg.Conditions) > 0 {
			return true
		}
	}
	return false
}

func getConditionGroupSeparator(cg structs.WhereGroup, index int) string {
	if cg.IsDummyGroup {
		return ""
	}
	if index == 0 {
		return ""
	}
	switch cg.Operator {
	case consts.LogicalOperator_AND:
		return " AND "
	case consts.LogicalOperator_OR:
		return " OR "
	}
	return ""
}

func getNotSeparator(cg structs.WhereGroup) string {
	if cg.IsNot {
		return "NOT "
	}
	return ""
}

func getParenthesesOpen(cg structs.WhereGroup) string {
	if cg.IsDummyGroup {
		return ""
	}
	return "("
}

func getParenthesesClose(cg structs.WhereGroup) string {
	if cg.IsDummyGroup {
		return ""
	}
	return ")"
}

func getConditionOperator(c structs.Where) string {
	switch c.Operator {
	case consts.LogicalOperator_AND:
		return " AND "
	case consts.LogicalOperator_OR:
		return " OR "
	}
	return ""
}

func processSubQuery(sb *strings.Builder, c structs.Where) []interface{} {
	condQuery := c.Column + " " + c.Condition
	b := &BaseQueryBuilder{}
	sqQuery, sqValues := b.Build("", c.Query)

	sb.WriteString(condQuery + " (" + sqQuery + ")")
	return sqValues
}

func processExistsQuery(sb *strings.Builder, c structs.Where) []interface{} {
	condQuery := c.Condition
	b := &BaseQueryBuilder{}
	sqQuery, sqValues := b.Build("", c.Exists.Query)

	sb.WriteString(condQuery + " (" + sqQuery + ")")
	return sqValues
}

func processBetweenCondition(sb *strings.Builder, c structs.Where) []interface{} {
	wsb := strings.Builder{}
	wsb.Grow(consts.StringBuffer_Where_Grow)
	values := make([]interface{}, 0, 2)
	if c.Between.IsColumn {
		wsb.WriteString(c.Column + " " + c.Condition + " " + c.Between.From.(string) + " AND " + c.Between.To.(string))
	} else {
		wsb.WriteString(c.Column + " " + c.Condition + " ? AND ?")
		values = []interface{}{c.Between.From, c.Between.To}
	}

	condQuery := wsb.String()

	sb.WriteString(condQuery)
	return values
}

func processRawCondition(sb *strings.Builder, c structs.Where) []interface{} {
	wsb := strings.Builder{}
	wsb.Grow(consts.StringBuffer_Where_Grow)

	if c.Raw != "" {
		wsb.WriteString(c.Raw)
	} else {
		wsb.WriteString(c.Column + " " + c.Condition)
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
	values := c.Value

	sb.WriteString(condQuery)
	return values
}
