package mysql

import (
	"encoding/json"
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

func NewWhereMySQLBuilder(util interfaces.SQLUtils, wg []structs.WhereGroup) *WhereMySQLBuilder {
	return &WhereMySQLBuilder{
		whereBaseBuilder: base.NewWhereBaseBuilder(util, wg),
		u:                util,
	}
}

func (wb *WhereMySQLBuilder) Where(sb *[]byte, wg []structs.WhereGroup) ([]interface{}, error) {
	if len(wg) == 0 {
		return []interface{}{}, nil
	}

	// WHERE
	if wb.whereBaseBuilder.HasCondition(wg) {
		*sb = append(*sb, " WHERE "...)
	}

	values := make([]interface{}, 0)

	for i := range wg {
		if len((wg)[i].Conditions) == 0 {
			continue
		}

		if i > 0 {
			*sb = append(*sb, wb.WhereBaseBuilder.GetConditionGroupSeparator((wg)[i], i)...)
		}

		*sb = append(*sb, wb.whereBaseBuilder.GetNotSeparator((wg)[i])...)
		*sb = append(*sb, wb.whereBaseBuilder.GetParenthesesOpen((wg)[i])...)

		for j := range (wg)[i].Conditions {
			if j > 0 || (i > 0 && j == 0 && (wg)[i].IsDummyGroup) {
				*sb = append(*sb, wb.whereBaseBuilder.GetConditionOperator((wg)[i].Conditions[j])...)
			}

			switch {
			case (wg)[i].Conditions[j].Query != nil:
				values = append(values, wb.whereBaseBuilder.ProcessSubQuery(sb, (wg)[i].Conditions[j])...)
			case (wg)[i].Conditions[j].Exists != nil:
				values = append(values, wb.whereBaseBuilder.ProcessExistsQuery(sb, (wg)[i].Conditions[j])...)
			case (wg)[i].Conditions[j].Between != nil:
				values = append(values, wb.whereBaseBuilder.ProcessBetweenCondition(sb, (wg)[i].Conditions[j])...)
			case (wg)[i].Conditions[j].FullText != nil:
				values = append(values, wb.ProcessFullText(sb, (wg)[i].Conditions[j])...)
			case (wg)[i].Conditions[j].JsonContains != nil:
				values = append(values, wb.ProcessJsonContains(sb, (wg)[i].Conditions[j])...)
			case (wg)[i].Conditions[j].JsonLength != nil:
				values = append(values, wb.ProcessJsonLength(sb, (wg)[i].Conditions[j])...)
			case (wg)[i].Conditions[j].Function != "":
				values = append(values, wb.whereBaseBuilder.ProcessFunction(sb, (wg)[i].Conditions[j])...)
			default:
				values = append(values, wb.whereBaseBuilder.ProcessRawCondition(sb, (wg)[i].Conditions[j])...)
			}
		}
		*sb = append(*sb, wb.whereBaseBuilder.GetParenthesesClose((wg)[i])...)
	}

	return values, nil
}

func (wb *WhereMySQLBuilder) ProcessFullText(sb *[]byte, c structs.Where) []interface{} {
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

	*sb = append(*sb, "MATCH ("...)
	for i, column := range c.FullText.Columns {
		if i > 0 {
			*sb = append(*sb, ", "...)
		}
		*sb = wb.u.EscapeIdentifier(*sb, column)
	}
	*sb = append(*sb, ") AGAINST ("+wb.u.GetPlaceholder()+" "+mode+expand+")"...)
	values := append(c.Value, c.FullText.Search)

	return values
}

func (wb *WhereMySQLBuilder) ProcessJsonContains(sb *[]byte, c structs.Where) []interface{} {
	field, path := parseJsonFieldAndPath(c.Column)
	*sb = append(*sb, "JSON_CONTAINS("...)
	*sb = wb.u.EscapeIdentifier(*sb, field)
	*sb = append(*sb, ", "...)
	*sb = append(*sb, wb.u.GetPlaceholder()...)
	if path != "" {
		*sb = append(*sb, ", '$."+path+"')"...)
	} else {
		*sb = append(*sb, ")"...)
	}

	var jsonVal []byte
	if len(c.JsonContains.Values) == 1 {
		jsonVal, _ = json.Marshal(c.JsonContains.Values[0])
	} else {
		jsonVal, _ = json.Marshal(c.JsonContains.Values)
	}
	return []interface{}{string(jsonVal)}
}

func (wb *WhereMySQLBuilder) ProcessJsonLength(sb *[]byte, c structs.Where) []interface{} {
	field, path := parseJsonFieldAndPath(c.Column)
	*sb = append(*sb, "JSON_LENGTH("...)
	*sb = wb.u.EscapeIdentifier(*sb, field)
	if path != "" {
		*sb = append(*sb, ", '$."+path+"')"...)
	} else {
		*sb = append(*sb, ")"...)
	}
	*sb = append(*sb, " "...)
	*sb = append(*sb, c.JsonLength.Operator...)
	*sb = append(*sb, " "...)
	*sb = append(*sb, wb.u.GetPlaceholder()...)
	return []interface{}{c.JsonLength.Value}
}

func parseJsonFieldAndPath(column string) (string, string) {
	parts := strings.Split(column, "->")
	field := parts[0]
	if len(parts) > 1 {
		return field, strings.Join(parts[1:], ".")
	}
	return field, ""
}
