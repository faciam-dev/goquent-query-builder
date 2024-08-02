package base

import (
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/sliceutils"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type SelectBaseBuilder struct {
	columnNames *[]string
}

func NewSelectBaseBuilder(columnNames *[]string) *SelectBaseBuilder {
	return &SelectBaseBuilder{
		columnNames: columnNames,
	}
}

func (b *SelectBaseBuilder) Select(sb *strings.Builder, columns *[]structs.Column, tableName string, joins *structs.Joins) []interface{} {
	if columns == nil {
		sb.WriteString(" * ")
		return []interface{}{}
	}
	//colNames := make([]string, 0, len(*columns))

	// if there are no columns to select, select all columns
	if len(*columns) == 0 && joins.Joins != nil {
		for i, join := range *joins.Joins {
			b.processJoin(sb, &join, tableName, i)
		}

		if joins.JoinClause != nil {
			for _, joinClause := range *joins.JoinClause {
				join := structs.Join{
					TargetNameMap: joinClause.TargetNameMap,
					Name:          joinClause.Name,
				}
				b.processJoin(sb, &join, tableName, 0)
			}
		}

		return []interface{}{}
	}

	colValues := make([]interface{}, 0, len(*columns))
	firstDistinct := false

	// if there are columns to select
	for i, column := range *columns {
		if column.Distinct && !column.Count && !firstDistinct {
			sb.WriteString("DISTINCT ")
			firstDistinct = true
		}

		if column.Count {
			sb.WriteString("COUNT(")
			if column.Distinct {
				sb.WriteString("DISTINCT ")
			}
			if column.Name != "" {
				sb.WriteString(column.Name)
			} else {
				sb.WriteString("*")
			}
			sb.WriteString(")")
			if i < len(*columns)-1 {
				sb.WriteString(", ")
			}

			continue
		}

		if column.Raw != "" {
			if column.Values != nil && len(column.Values) > 0 {
				colValues = append(colValues, column.Values...)
			}
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(column.Raw) // or colNames = column.Raw
			//colNames = append(colNames, column.Raw)
		} else if column.Name != "" {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(column.Name)
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

	targetNameForSelect := targetName + ".*"

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

	nameForSelect := name + ".*"
	if !sliceutils.Contains(*j.columnNames, nameForSelect) {
		if idx > 0 || outputed {
			sb.WriteString(", ")
		}
		sb.WriteString(nameForSelect)
		*j.columnNames = append(*j.columnNames, nameForSelect)
	}

}
