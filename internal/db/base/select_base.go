package base

import (
	"strings"

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

func (b *SelectBaseBuilder) Select(sb *strings.Builder, columns *[]structs.Column, tableName string, joins *structs.Joins) []interface{} {
	if columns == nil {
		sb.WriteString("*")
		return []interface{}{}
	}
	//colNames := make([]string, 0, len(*columns))

	outputed := false
	// if there are no columns to select, select all columns
	if len(*columns) == 0 && (joins.Joins != nil || joins.LateralJoins != nil) {
		//joined := append(*joins.LateralJoins, *joins.Joins)
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
		sb.WriteString("*")
		return []interface{}{}
	}

	colValues := make([]interface{}, 0, len(*columns))
	firstDistinct := false

	// if there are columns to select
	for i := 0; i < len(*columns); i++ {
		//	for i := range *columns {
		if (*columns)[i].Distinct && !(*columns)[i].Count && !firstDistinct {
			sb.WriteString("DISTINCT ")
			firstDistinct = true
		}

		if (*columns)[i].Count {
			sb.WriteString("COUNT(")
			if (*columns)[i].Distinct {
				sb.WriteString("DISTINCT ")
			}
			if (*columns)[i].Name != "" {
				b.u.EscapeIdentifierAliasedValue(sb, (*columns)[i].Name)
			} else {
				sb.WriteString("*")
			}
			sb.WriteString(")")
			if i < len(*columns)-1 {
				sb.WriteString(", ")
			}

			continue
		}

		if (*columns)[i].Function != "" {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString((*columns)[i].Function)
			sb.WriteString("(")
			if (*columns)[i].Distinct {
				sb.WriteString("DISTINCT ")
			}
			if (*columns)[i].Name != "" {
				b.u.EscapeIdentifierAliasedValue(sb, (*columns)[i].Name)
			} else {
				sb.WriteString("*")
			}
			sb.WriteString(")")
		} else if (*columns)[i].Raw != "" {
			if (*columns)[i].Values != nil && len((*columns)[i].Values) > 0 {
				colValues = append(colValues, (*columns)[i].Values...)
			}
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString((*columns)[i].Raw) // or colNames = column.Raw
			//colNames = append(colNames, column.Raw)
		} else if (*columns)[i].Name != "" {
			if i > 0 {
				sb.WriteString(", ")
			}

			b.u.EscapeIdentifierAliasedValue(sb, (*columns)[i].Name)
			//colNames = append(colNames, column.Name)
		}
	}

	return colValues
}

func (j *SelectBaseBuilder) processJoin(sb *strings.Builder, join *structs.Join, tableName string, idx int) {
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

	wsb := strings.Builder{}
	wsb.Grow(consts.StringBuffer_Short_Query_Grow)
	j.u.EscapeIdentifier(&wsb, targetName)
	wsb.WriteString(".*")
	targetNameForSelect := wsb.String()
	wsb.Reset()

	//targetNameForSelect := j.u.EscapeIdentifier(&wsb, targetName) + ".*"

	//targetNameForSelect := j.u.EscapeIdentifier(&wsb, targetName) + ".*"

	//sb.Grow(consts.StringBuffer_Select_Grow)

	outputed := false
	if !sliceutils.Contains(*j.columnNames, targetNameForSelect) {
		if idx > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(targetNameForSelect)
		*j.columnNames = append(*j.columnNames, targetNameForSelect)
		outputed = true
	}

	wsb.Grow(consts.StringBuffer_Short_Query_Grow)
	j.u.EscapeIdentifier(&wsb, name)
	wsb.WriteString(".*")
	nameForSelect := wsb.String()
	wsb.Reset()

	//nameForSelect := j.u.EscapeIdentifier(sb, name) + ".*"
	if !sliceutils.Contains(*j.columnNames, nameForSelect) {
		if idx > 0 || outputed {
			sb.WriteString(", ")
		}
		sb.WriteString(nameForSelect)
		*j.columnNames = append(*j.columnNames, nameForSelect)
	}

}
