package db

import (
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/sliceutils"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type JoinBaseBuilder struct {
	join *structs.Joins
}

func NewJoinBaseBuilder(j *structs.Joins) *JoinBaseBuilder {
	return &JoinBaseBuilder{
		join: j,
	}
}

// Join builds the JOIN query.
func (BaseQueryBuilder) Join(tableName string, joins *structs.Joins) (*[]structs.Column, string) {
	join := ""

	joinedTablesForSelect, joinStrings := buildJoinStatement(tableName, joins)
	for _, joinString := range joinStrings {
		join += " " + joinString
	}

	return joinedTablesForSelect, join
}

// buildJoinStatement builds the JOIN statement.
func buildJoinStatement(tableName string, joins *structs.Joins) (*[]structs.Column, []string) {
	joinedTablesForSelect := make([]structs.Column, 0, len(*joins.Joins))
	joinStrings := make([]string, 0, len(*joins.Joins))
	for _, join := range *joins.Joins {
		joinType := ""
		targetName := ""

		if _, ok := join.TargetNameMap[consts.Join_CROSS]; ok {
			targetName = join.TargetNameMap[consts.Join_CROSS]
			joinType = consts.Join_Type_CROSS
		}
		if _, ok := join.TargetNameMap[consts.Join_RIGHT]; ok {
			targetName = join.TargetNameMap[consts.Join_RIGHT]
			joinType = consts.Join_Type_RIGHT
		}
		if _, ok := join.TargetNameMap[consts.Join_LEFT]; ok {
			targetName = join.TargetNameMap[consts.Join_LEFT]
			joinType = consts.Join_Type_LEFT
		}
		if _, ok := join.TargetNameMap[consts.Join_INNER]; ok {
			targetName = join.TargetNameMap[consts.Join_INNER]
			joinType = consts.Join_Type_INNER
		}

		if targetName == "" {
			continue
		}

		name := tableName
		if join.Name != "" {
			name = join.Name
		}

		targetNameForSelect := targetName + ".*"
		if !sliceutils.Contains[string](*getNowColNames(&joinedTablesForSelect), targetNameForSelect) {
			joinedTablesForSelect = append(joinedTablesForSelect, structs.Column{
				Name: targetNameForSelect,
			})
		}
		nameForSelect := name + ".*"
		if !sliceutils.Contains[string](*getNowColNames(&joinedTablesForSelect), nameForSelect) {
			joinedTablesForSelect = append(joinedTablesForSelect, structs.Column{
				Name: nameForSelect,
			})
		}

		joinQuery := joinType + " JOIN " + targetName + " ON " + join.SearchColumn + " " + join.SearchCondition + " " + join.SearchTargetColumn
		if _, ok := join.TargetNameMap[consts.Join_CROSS]; ok {
			joinQuery = joinType + " JOIN " + targetName
		}

		joinStrings = append(joinStrings, joinQuery)

	}

	return &joinedTablesForSelect, joinStrings
}

// getNowColNames returns the names of the columns in the slice.
func getNowColNames(joinedTablesForSelect *[]structs.Column) *[]string {
	nowColNames := make([]string, len(*joinedTablesForSelect))
	for _, joinedTable := range *joinedTablesForSelect {
		nowColNames = append(nowColNames, joinedTable.Name)
	}
	return &nowColNames
}
