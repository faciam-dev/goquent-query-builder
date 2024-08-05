package base

import (
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

type GroupByBaseBuilder struct {
	u interfaces.SQLUtils
}

func NewGroupByBaseBuilder(util interfaces.SQLUtils) *GroupByBaseBuilder {
	return &GroupByBaseBuilder{
		u: util,
	}
}

func (g GroupByBaseBuilder) GroupBy(sb *strings.Builder, groupBy *structs.GroupBy) []interface{} {
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
			sb.WriteString(g.u.EscapeIdentifier(sb, column))
		}
	}

	values := make([]interface{}, 0, len(*groupBy.Having))

	if len(*groupBy.Having) > 0 {
		sb.WriteString(" HAVING ")

		//havingValues := make([]interface{}, 0, len(*groupBy.Having))
		for n := range *groupBy.Having {
			op := "AND"
			if (*groupBy.Having)[n].Operator == consts.LogicalOperator_AND {
				op = "AND"
			} else if (*groupBy.Having)[n].Operator == consts.LogicalOperator_OR {
				op = "OR"
			}

			if (*groupBy.Having)[n].Raw != "" {
				if n > 0 {
					sb.WriteString(" ")
					sb.WriteString(op)
					sb.WriteString(" ")
				}
				sb.WriteString((*groupBy.Having)[n].Raw)
				continue
			}
			if (*groupBy.Having)[n].Column == "" {
				continue
			}
			if (*groupBy.Having)[n].Condition == "" {
				continue
			}
			if (*groupBy.Having)[n].Value == "" {
				continue
			}
			//havingValues = append(havingValues, having.Value)
			values = append(values, (*groupBy.Having)[n].Value)

			if n > 0 {
				sb.WriteString(" ")
				sb.WriteString(op)
				sb.WriteString(" ")
			}
			sb.WriteString(g.u.EscapeIdentifier(sb, (*groupBy.Having)[n].Column))
			sb.WriteString(" ")
			sb.WriteString((*groupBy.Having)[n].Condition)
			sb.WriteString(" " + g.u.GetPlaceholder())
		}

		//if len(havingValues) > 0 {
		//	values = append(values, havingValues...)
		//}
	}

	return values
}
