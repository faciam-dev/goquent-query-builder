package base

import (
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/sliceutils"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

type SelectBaseBuilder struct {
	columnNames *[]string
	u           interfaces.SQLUtils
}

func NewSelectBaseBuilder(u interfaces.SQLUtils, columnNames *[]string) *SelectBaseBuilder {
	return &SelectBaseBuilder{
		columnNames: columnNames,
		u:           u,
	}
}

func (b *SelectBaseBuilder) Select(sb *[]byte, columns *[]structs.Column, tableName string, joins *structs.Joins) []interface{} {
	if columns == nil {
		*sb = append(*sb, "*"...)
		return []interface{}{}
	}

	outputed := false
	// if there are no columns to select, select all columns
	if len(*columns) == 0 && (joins.Joins != nil || joins.LateralJoins != nil) {
		for i, join := range append(*joins.LateralJoins, *joins.Joins...) {
			b.processJoin(sb, &join, tableName, i)
			outputed = true
		}

		if joins.JoinClauses != nil {
			for _, joinClause := range *joins.JoinClauses {
				join := structs.Join{
					TargetNameMap: joinClause.TargetNameMap,
					Name:          joinClause.Name,
				}
				b.processJoin(sb, &join, tableName, 0)
				outputed = true
			}
		}

	}

	if len(*columns) == 0 && !outputed {
		*sb = append(*sb, "*"...)
		return []interface{}{}
	}

	// if there are columns has values
	var colValues []interface{}
	hasValues := false
	for i := 0; i < len(*columns); i++ {
		if (*columns)[i].Values != nil && len((*columns)[i].Values) > 0 {
			hasValues = true
			break
		}
	}
	if hasValues {
		colValues = make([]interface{}, 0, len(*columns))
	}

	// if there are columns to select
	firstDistinct := false
	for i := 0; i < len(*columns); i++ {
		if (*columns)[i].Distinct && !(*columns)[i].Count && !firstDistinct {
			*sb = append(*sb, "DISTINCT "...)
			firstDistinct = true
		}

		if (*columns)[i].Count {
			*sb = append(*sb, "COUNT("...)
			if (*columns)[i].Distinct {
				*sb = append(*sb, "DISTINCT "...)
			}
			if (*columns)[i].Name != "" {
				*sb = b.u.EscapeIdentifierAliasedValue(*sb, (*columns)[i].Name)
			} else {
				*sb = append(*sb, "*"...)
			}
			*sb = append(*sb, ")"...)
			if i < len(*columns)-1 {
				*sb = append(*sb, ", "...)
			}

			continue
		}

		if (*columns)[i].Function != "" {
			if i > 0 {
				*sb = append(*sb, ", "...)
			}
			*sb = append(*sb, (*columns)[i].Function...)
			*sb = append(*sb, "("...)
			if (*columns)[i].Distinct {
				*sb = append(*sb, "DISTINCT "...)
			}
			if (*columns)[i].Name != "" {
				*sb = b.u.EscapeIdentifierAliasedValue(*sb, (*columns)[i].Name)
			} else {
				*sb = append(*sb, "*"...)
			}
			*sb = append(*sb, ")"...)
		} else if (*columns)[i].Raw != "" {
			if (*columns)[i].Values != nil && len((*columns)[i].Values) > 0 {
				colValues = append(colValues, (*columns)[i].Values...)
			}
			if i > 0 {
				*sb = append(*sb, ", "...)
			}
			*sb = append(*sb, (*columns)[i].Raw...) // or colNames = column.Raw
		} else if (*columns)[i].Name != "" {
			if i > 0 {
				*sb = append(*sb, ", "...)
			}

			*sb = b.u.EscapeIdentifierAliasedValue(*sb, (*columns)[i].Name)
		}
	}

	return colValues
}

func (j *SelectBaseBuilder) processJoin(sb *[]byte, join *structs.Join, tableName string, idx int) {
	targetName := ""
	//joinedTablesForSelect := ""

	if _, ok := join.TargetNameMap[consts.Join_CROSS]; ok {
		targetName = join.TargetNameMap[consts.Join_CROSS]
	}
	if _, ok := join.TargetNameMap[consts.Join_RIGHT]; ok {
		targetName = join.TargetNameMap[consts.Join_RIGHT]
	}
	if _, ok := join.TargetNameMap[consts.Join_LEFT]; ok {
		targetName = join.TargetNameMap[consts.Join_LEFT]
	}
	if _, ok := join.TargetNameMap[consts.Join_INNER]; ok {
		targetName = join.TargetNameMap[consts.Join_INNER]
	}
	if _, ok := join.TargetNameMap[consts.Join_LATERAL]; ok {
		targetName = join.TargetNameMap[consts.Join_LATERAL]
	}
	if _, ok := join.TargetNameMap[consts.Join_LEFT_LATERAL]; ok {
		targetName = join.TargetNameMap[consts.Join_LEFT_LATERAL]
	}

	if targetName == "" {
		return
	}

	name := tableName
	if join.Name != "" {
		name = join.Name
	}

	wsb := make([]byte, 0, consts.StringBuffer_Short_Query_Grow)
	wsb = j.u.EscapeIdentifier2(wsb, targetName)
	wsb = append(wsb, ".*"...)
	targetNameForSelect := string(wsb)
	wsb = wsb[:0]

	outputed := false
	if !sliceutils.Contains(*j.columnNames, targetNameForSelect) {
		if idx > 0 {
			*sb = append(*sb, ", "...)
		}
		*sb = append(*sb, targetNameForSelect...)
		*j.columnNames = append(*j.columnNames, targetNameForSelect)
		outputed = true
	}

	wsb = j.u.EscapeIdentifier2(wsb, name)
	wsb = append(wsb, ".*"...)
	nameForSelect := string(wsb)

	if !sliceutils.Contains(*j.columnNames, nameForSelect) {
		if idx > 0 || outputed {
			*sb = append(*sb, ", "...)
		}
		*sb = append(*sb, nameForSelect...)
		*j.columnNames = append(*j.columnNames, nameForSelect)
	}

}
