package base

import (
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

type WhereBaseBuilder struct {
	u           interfaces.SQLUtils
	whereGroups []structs.WhereGroup
}

func NewWhereBaseBuilder(util interfaces.SQLUtils, wg []structs.WhereGroup) *WhereBaseBuilder {
	return &WhereBaseBuilder{
		u:           util,
		whereGroups: wg,
	}
}

func (wb *WhereBaseBuilder) Where(sb *strings.Builder, wg []structs.WhereGroup) []interface{} {
	if len(wg) == 0 {
		return []interface{}{}
	}

	// WHERE
	if wb.HasCondition(wg) {
		sb.WriteString(" WHERE ")
	}

	values := make([]interface{}, 0)

	for i, cg := range wg {
		if len(cg.Conditions) == 0 {
			continue
		}

		if i > 0 {
			sb.WriteString(wb.GetConditionGroupSeparator(cg, i))
		}

		sb.WriteString(wb.GetNotSeparator(cg))
		sb.WriteString(wb.GetParenthesesOpen(cg))

		for j, c := range cg.Conditions {
			if j > 0 || (i > 0 && j == 0 && cg.IsDummyGroup) {
				sb.WriteString(wb.GetConditionOperator(c))
			}

			switch {
			case c.Query != nil:
				values = append(values, wb.ProcessSubQuery(sb, c)...)
			case c.Exists != nil:
				values = append(values, wb.ProcessExistsQuery(sb, c)...)
			case c.Between != nil:
				values = append(values, wb.ProcessBetweenCondition(sb, c)...)
			case c.FullText != nil:
				values = append(values, wb.ProcessFullText(sb, c)...)
			case c.Function != "":
				values = append(values, wb.ProcessFunction(sb, c)...)
			default:
				values = append(values, wb.ProcessRawCondition(sb, c)...)
			}
		}
		sb.WriteString(wb.GetParenthesesClose(cg))
	}

	return values
}

func (wb *WhereBaseBuilder) HasCondition(wg []structs.WhereGroup) bool {
	for _, cg := range wg {
		if len(cg.Conditions) > 0 {
			return true
		}
	}
	return false
}

func (wb *WhereBaseBuilder) GetConditionGroupSeparator(cg structs.WhereGroup, index int) string {
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

func (wb *WhereBaseBuilder) GetNotSeparator(cg structs.WhereGroup) string {
	if cg.IsNot {
		return "NOT "
	}
	return ""
}

func (wb *WhereBaseBuilder) GetParenthesesOpen(cg structs.WhereGroup) string {
	if cg.IsDummyGroup {
		return ""
	}
	return "("
}

func (wb *WhereBaseBuilder) GetParenthesesClose(cg structs.WhereGroup) string {
	if cg.IsDummyGroup {
		return ""
	}
	return ")"
}

func (wb *WhereBaseBuilder) GetConditionOperator(c structs.Where) string {
	switch c.Operator {
	case consts.LogicalOperator_AND:
		return " AND "
	case consts.LogicalOperator_OR:
		return " OR "
	}
	return ""
}

func (wb *WhereBaseBuilder) ProcessSubQuery(sb *strings.Builder, c structs.Where) []interface{} {
	wb.u.EscapeIdentifier(sb, c.Column)
	sb.WriteString(" ")
	sb.WriteString(c.Condition)

	sb.WriteString(" (")

	b := wb.u.GetQueryBuilderStrategy()
	sqValues := b.Build(sb, c.Query, 0, nil)

	//sb.WriteString(sqQuery)
	sb.WriteString(")")
	return sqValues
}

func (wb *WhereBaseBuilder) ProcessExistsQuery(sb *strings.Builder, c structs.Where) []interface{} {
	sb.WriteString(c.Condition)

	//condQuery := c.Condition
	sb.WriteString(" (")
	b := wb.u.GetQueryBuilderStrategy()
	sqValues := b.Build(sb, c.Exists.Query, 0, nil)
	sb.WriteString(")")

	//sb.WriteString(condQuery + " (" + sqQuery + ")")
	return sqValues
}

func (wb *WhereBaseBuilder) ProcessBetweenCondition(sb *strings.Builder, c structs.Where) []interface{} {
	//wsb := strings.Builder{}
	//wsb.Grow(consts.StringBuffer_Where_Grow)
	values := make([]interface{}, 0, 2)
	if c.Between.IsColumn {
		wb.u.EscapeIdentifier(sb, c.Column)
		sb.WriteString(" ")
		sb.WriteString(c.Condition)
		sb.WriteString(" ")
		wb.u.EscapeIdentifier(sb, c.Between.From.(string))
		sb.WriteString(" AND ")
		wb.u.EscapeIdentifier(sb, c.Between.To.(string))
	} else {
		wb.u.EscapeIdentifier(sb, c.Column)
		sb.WriteString(" ")
		sb.WriteString(c.Condition)
		sb.WriteString(" ")
		sb.WriteString(wb.u.GetPlaceholder())
		sb.WriteString(" AND ")
		sb.WriteString(wb.u.GetPlaceholder())
		values = []interface{}{c.Between.From, c.Between.To}
	}

	//condQuery := wsb.String()

	//sb.WriteString(condQuery)
	return values
}

func (wb *WhereBaseBuilder) ProcessRawCondition(sb *strings.Builder, c structs.Where) []interface{} {
	//wsb := strings.Builder{}
	//wsb.Grow(consts.StringBuffer_Where_Grow)

	if c.Raw != "" {
		sb.WriteString(c.Raw)
	} else {
		wb.u.EscapeIdentifier(sb, c.Column)
		sb.WriteString(" ")
		sb.WriteString(c.Condition)
		if c.ValueColumn != "" {
			sb.WriteString(" ")
			wb.u.EscapeIdentifier(sb, c.ValueColumn)
		} else if c.Value != nil {
			if len(c.Value) > 1 {
				sb.WriteString(" (")
				for k := 0; k < len(c.Value); k++ {
					if k > 0 {
						sb.WriteString(", ")
					}
					sb.WriteString(wb.u.GetPlaceholder())
				}
				sb.WriteString(")")
			} else {
				sb.WriteString(" ")
				sb.WriteString(wb.u.GetPlaceholder())
			}
		}
	}

	//condQuery := wsb.String()
	values := c.Value

	//sb.WriteString(condQuery)
	return values
}

func (wb *WhereBaseBuilder) ProcessFullText(sb *strings.Builder, c structs.Where) []interface{} {
	values := []interface{}{}

	// Implement FullText

	return values
}

func (wb *WhereBaseBuilder) ProcessFunction(sb *strings.Builder, c structs.Where) []interface{} {
	//wsb := strings.Builder{}
	//wsb.Grow(consts.StringBuffer_Where_Grow)
	sb.WriteString(c.Function)
	sb.WriteString("(")
	wb.u.EscapeIdentifier(sb, c.Column)
	sb.WriteString(") ")
	sb.WriteString(c.Condition)
	if c.ValueColumn != "" {
		sb.WriteString(" ")
		wb.u.EscapeIdentifier(sb, c.ValueColumn)
	} else if c.Value != nil {
		if len(c.Value) > 1 {
			sb.WriteString(" (")
			for k := 0; k < len(c.Value); k++ {
				if k > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(wb.u.GetPlaceholder())
			}
			sb.WriteString(")")
		} else {
			sb.WriteString(" ")
			sb.WriteString(wb.u.GetPlaceholder())
		}
	}

	//condQuery := wsb.String()
	values := c.Value

	//sb.WriteString(condQuery)
	return values
}
