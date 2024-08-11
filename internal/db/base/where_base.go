package base

import (
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

func (wb *WhereBaseBuilder) Where(sb *[]byte, wg []structs.WhereGroup) []interface{} {
	if len(wg) == 0 {
		return []interface{}{}
	}

	// WHERE
	if wb.HasCondition(wg) {
		*sb = append(*sb, " WHERE "...)
	}

	// estimate the cap of values
	cap := 0
	for _, cg := range wg {
		for _, c := range cg.Conditions {
			if c.Query != nil {
				cap += 5
				continue
			}
			if c.Exists != nil {
				cap += 5
				continue
			}
			if c.Between != nil {
				cap += 2
				continue
			}
			if c.FullText != nil {
				cap += 2
				continue
			}
			if c.Function != "" {
				cap += 5
				continue
			}
			if c.Raw != "" {
				cap += 1
				continue
			}
			if c.Value != nil {
				cap += len(c.Value)
				continue
			}
		}
	}

	values := make([]interface{}, 0, cap)

	for i, cg := range wg {
		if len(cg.Conditions) == 0 {
			continue
		}

		if i > 0 {
			*sb = append(*sb, wb.GetConditionGroupSeparator(cg, i)...)
		}

		*sb = append(*sb, wb.GetNotSeparator(cg)...)
		*sb = append(*sb, wb.GetParenthesesOpen(cg)...)

		for j, c := range cg.Conditions {
			if j > 0 || (i > 0 && j == 0 && cg.IsDummyGroup) {
				*sb = append(*sb, wb.GetConditionOperator(c)...)
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
		*sb = append(*sb, wb.GetParenthesesClose(cg)...)
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

func (wb *WhereBaseBuilder) ProcessSubQuery(sb *[]byte, c structs.Where) []interface{} {
	*sb = wb.u.EscapeIdentifier(*sb, c.Column)
	*sb = append(*sb, " "...)
	*sb = append(*sb, c.Condition...)

	*sb = append(*sb, " ("...)

	b := wb.u.GetQueryBuilderStrategy()
	sqValues := b.Build(sb, c.Query, 0, nil)

	*sb = append(*sb, ")"...)
	return sqValues
}

func (wb *WhereBaseBuilder) ProcessExistsQuery(sb *[]byte, c structs.Where) []interface{} {
	*sb = append(*sb, c.Condition...)

	*sb = append(*sb, " ("...)
	b := wb.u.GetQueryBuilderStrategy()
	sqValues := b.Build(sb, c.Exists.Query, 0, nil)
	*sb = append(*sb, ")"...)

	return sqValues
}

func (wb *WhereBaseBuilder) ProcessBetweenCondition(sb *[]byte, c structs.Where) []interface{} {
	values := make([]interface{}, 0, 2)
	if c.Between.IsColumn {
		*sb = wb.u.EscapeIdentifier(*sb, c.Column)
		*sb = append(*sb, " "...)
		*sb = append(*sb, c.Condition...)
		*sb = append(*sb, " "...)
		*sb = wb.u.EscapeIdentifier(*sb, c.Between.From.(string))
		*sb = append(*sb, " AND "...)
		*sb = wb.u.EscapeIdentifier(*sb, c.Between.To.(string))
	} else {
		*sb = wb.u.EscapeIdentifier(*sb, c.Column)
		*sb = append(*sb, " "...)
		*sb = append(*sb, c.Condition...)
		*sb = append(*sb, " "...)
		*sb = append(*sb, wb.u.GetPlaceholder()...)
		*sb = append(*sb, " AND "...)
		*sb = append(*sb, wb.u.GetPlaceholder()...)
		values = []interface{}{c.Between.From, c.Between.To}
	}

	return values
}

func (wb *WhereBaseBuilder) ProcessRawCondition(sb *[]byte, c structs.Where) []interface{} {
	if c.Raw != "" {
		*sb = append(*sb, c.Raw...)
	} else {
		*sb = wb.u.EscapeIdentifier(*sb, c.Column)
		*sb = append(*sb, " "...)
		*sb = append(*sb, c.Condition...)
		if c.ValueColumn != "" {
			*sb = append(*sb, " "...)
			*sb = wb.u.EscapeIdentifier(*sb, c.ValueColumn)
		} else if c.Value != nil {
			if len(c.Value) > 1 {
				*sb = append(*sb, " ("...)
				for k := 0; k < len(c.Value); k++ {
					if k > 0 {
						*sb = append(*sb, ", "...)
					}
					*sb = append(*sb, wb.u.GetPlaceholder()...)
				}
				*sb = append(*sb, ")"...)
			} else {
				*sb = append(*sb, " "...)
				*sb = append(*sb, wb.u.GetPlaceholder()...)
			}
		}
	}

	values := c.Value

	return values
}

func (wb *WhereBaseBuilder) ProcessFullText(sb *[]byte, c structs.Where) []interface{} {
	values := []interface{}{}

	// Implement FullText

	return values
}

func (wb *WhereBaseBuilder) ProcessFunction(sb *[]byte, c structs.Where) []interface{} {
	*sb = append(*sb, c.Function...)
	*sb = append(*sb, "("...)
	*sb = wb.u.EscapeIdentifier(*sb, c.Column)
	*sb = append(*sb, ") "...)
	*sb = append(*sb, c.Condition...)
	if c.ValueColumn != "" {
		*sb = append(*sb, " "...)
		*sb = wb.u.EscapeIdentifier(*sb, c.ValueColumn)
	} else if c.Value != nil {
		if len(c.Value) > 1 {
			*sb = append(*sb, " ("...)
			for k := 0; k < len(c.Value); k++ {
				if k > 0 {
					*sb = append(*sb, ", "...)
				}
				*sb = append(*sb, wb.u.GetPlaceholder()...)
			}
			*sb = append(*sb, ")"...)
		} else {
			*sb = append(*sb, " "...)
			*sb = append(*sb, wb.u.GetPlaceholder()...)
		}
	}

	values := c.Value

	return values
}
