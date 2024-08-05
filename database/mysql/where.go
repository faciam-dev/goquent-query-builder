package mysql

import (
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db/base"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

type WhereMySQLBuilder struct {
	base.WhereBaseBuilder
	whereBaseBuilder *base.WhereBaseBuilder
	u                interfaces.SQLUtils
}

func NewWhereMySQLBuilder(util interfaces.SQLUtils, wg *[]structs.WhereGroup) *WhereMySQLBuilder {
	return &WhereMySQLBuilder{
		whereBaseBuilder: base.NewWhereBaseBuilder(util, wg),
		u:                util,
	}
}

func (wb *WhereMySQLBuilder) Where(sb *strings.Builder, wg *[]structs.WhereGroup) []interface{} {
	if wg == nil || len(*wg) == 0 {
		return []interface{}{}
	}

	// WHERE
	if wb.whereBaseBuilder.HasCondition(*wg) {
		sb.WriteString(" WHERE ")
	}

	values := make([]interface{}, 0)

	for i, cg := range *wg {
		if len(cg.Conditions) == 0 {
			continue
		}

		if i > 0 {
			sb.WriteString(wb.WhereBaseBuilder.GetConditionGroupSeparator(cg, i))
		}

		sb.WriteString(wb.whereBaseBuilder.GetNotSeparator(cg))
		sb.WriteString(wb.whereBaseBuilder.GetParenthesesOpen(cg))

		for j, c := range cg.Conditions {
			if j > 0 || (i > 0 && j == 0 && cg.IsDummyGroup) {
				sb.WriteString(wb.whereBaseBuilder.GetConditionOperator(c))
			}

			switch {
			case c.Query != nil:
				values = append(values, wb.whereBaseBuilder.ProcessSubQuery(sb, c)...)
			case c.Exists != nil:
				values = append(values, wb.whereBaseBuilder.ProcessExistsQuery(sb, c)...)
			case c.Between != nil:
				values = append(values, wb.whereBaseBuilder.ProcessBetweenCondition(sb, c)...)
			case c.FullText != nil:
				values = append(values, wb.ProcessFullText(sb, c)...)
			case c.Function != "":
				values = append(values, wb.whereBaseBuilder.ProcessFunction(sb, c)...)
			default:
				values = append(values, wb.whereBaseBuilder.ProcessRawCondition(sb, c)...)
			}
		}
		sb.WriteString(wb.whereBaseBuilder.GetParenthesesClose(cg))
	}

	return values
}

func (wb *WhereMySQLBuilder) ProcessFullText(sb *strings.Builder, c structs.Where) []interface{} {
	// parse options
	mode := "IN NATURAL LANGUAGE MODE"
	expand := ""
	if c.FullText.Options != nil {
		if mmode, ok := c.FullText.Options["mode"]; ok {
			if mmode.(string) == "boolean" {
				mode = "IN BOOLEAN MODE"
			}
		}
		if with, ok := c.FullText.Options["expanded"]; ok {
			if with.(bool) {
				expand = " WITH QUERY EXPANSION"
			}
		}
	}

	sb.WriteString("MATCH (")
	for i, column := range c.FullText.Columns {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(wb.u.EscapeIdentifier(sb, column))
	}
	sb.WriteString(") AGAINST (" + wb.u.GetPlaceholder() + " " + mode + expand + ")")
	values := []interface{}{c.Value}

	return values
}
