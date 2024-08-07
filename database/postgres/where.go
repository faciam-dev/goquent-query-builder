package postgres

import (
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db/base"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

type WherePostgreSQLBuilder struct {
	base.WhereBaseBuilder
	whereBaseBuilder *base.WhereBaseBuilder
	u                interfaces.SQLUtils
}

func NewWherePostgreSQLBuilder(util interfaces.SQLUtils, wg []structs.WhereGroup) *WherePostgreSQLBuilder {
	return &WherePostgreSQLBuilder{
		whereBaseBuilder: base.NewWhereBaseBuilder(util, wg),
		u:                util,
	}
}

func (wb *WherePostgreSQLBuilder) Where(sb *strings.Builder, wg []structs.WhereGroup) []interface{} {
	if len(wg) == 0 {
		return []interface{}{}
	}

	// WHERE
	if wb.whereBaseBuilder.HasCondition(wg) {
		sb.WriteString(" WHERE ")
	}

	values := make([]interface{}, 0)

	for i, cg := range wg {
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

func (wb *WherePostgreSQLBuilder) ProcessFullText(sb *strings.Builder, c structs.Where) []interface{} {
	values := make([]interface{}, 0)

	// parse options
	language := "english"
	if c.FullText.Options != nil {
		if lang, ok := c.FullText.Options["language"]; ok {
			language = lang.(string)
		}
	}

	mode := "plainto_tsquery"
	if c.FullText.Options != nil {
		if mmode, ok := c.FullText.Options["mode"]; ok {
			if mmode.(string) == "phrase" {
				mode = "phraseto_tsquery"
			}
			if mmode.(string) == "websearch" {
				mode = "websearch_to_tsquery"
			}
		}
	}

	sb.WriteString("(")
	for i, column := range c.FullText.Columns {
		if i > 0 {
			sb.WriteString(" || ")
		}
		sb.WriteString("to_tsvector(")
		sb.WriteString(wb.u.GetPlaceholder())
		sb.WriteString(", ")
		wb.u.EscapeIdentifier(sb, column)
		sb.WriteString(")")
		values = append(values, language)
	}
	sb.WriteString(") @@ " + mode + "(" + wb.u.GetPlaceholder() + ", " + wb.u.GetPlaceholder() + ")")
	values = append(values, language, c.FullText.Search)

	return values
}
