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
			sb.WriteString(g.u.EscapeIdentifier(column))
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
			sb.WriteString(g.u.EscapeIdentifier(having.Column))
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
