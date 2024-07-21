package db

import (
	"strings"

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
func (BaseQueryBuilder) Join(tableName string, joins *structs.Joins) (*[]structs.Column, string, []interface{}) {
	join := ""

	joinedTablesForSelect, joinStrings, values := buildJoinStatement(tableName, joins)
	for _, joinString := range joinStrings {
		join += " " + joinString
	}

	return joinedTablesForSelect, join, values
}

// buildJoinStatement builds the JOIN statement.
func buildJoinStatement(tableName string, joins *structs.Joins) (*[]structs.Column, []string, []interface{}) {
	if joins.JoinClause != nil {
		joinedTablesForSelect := make([]structs.Column, 0, len(*joins.JoinClause.On))
		joinStrings := make([]string, 0, len(*joins.JoinClause.On))
		values := make([]interface{}, 0, len(*joins.JoinClause.Conditions))

		j := structs.Join{
			TargetNameMap: joins.JoinClause.TargetNameMap,
			Name:          joins.JoinClause.Name,
		}

		joinType, targetName := processJoin(j, tableName, &joinedTablesForSelect)

		if joins.JoinClause.Query != nil {
			b := &BaseQueryBuilder{}
			sqQuery, sqValues := b.Build(joins.JoinClause.Query)
			targetName = "(" + sqQuery + ")" + " AS " + targetName
			values = append(values, sqValues...)
		}

		var joinBuilder strings.Builder
		joinBuilder.WriteString(joinType)
		joinBuilder.WriteString(" JOIN ")
		joinBuilder.WriteString(targetName)
		joinBuilder.WriteString(" ON ")

		op := ""
		for i, on := range *joins.JoinClause.On {
			if i > 0 {
				if on.Operator == consts.LogicalOperator_OR {
					op = " OR "
				} else {
					op = " AND "
				}
			}

			joinBuilder.WriteString(op)
			joinBuilder.WriteString(on.Column)
			joinBuilder.WriteString(" ")
			joinBuilder.WriteString(on.Condition)
			if on.Value != "" {
				joinBuilder.WriteString(on.Value.(string)) // TODO: check if this is correct
			}
		}

		op = ""
		for i, condition := range *joins.JoinClause.Conditions {
			if i > 0 || len(*joins.JoinClause.On) > 0 {
				if condition.Operator == consts.LogicalOperator_OR {
					op = " OR "
				} else {
					op = " AND "
				}
			}
			joinBuilder.WriteString(op)
			joinBuilder.WriteString(condition.Column)
			joinBuilder.WriteString(" ")
			joinBuilder.WriteString(condition.Condition)
			joinBuilder.WriteString(" ?") // TODO: check if this is correct
			values = append(values, condition.Value...)
		}

		joinStrings = append(joinStrings, joinBuilder.String())

		return &joinedTablesForSelect, joinStrings, values
	}

	joinedTablesForSelect := make([]structs.Column, 0, len(*joins.Joins))
	joinStrings := make([]string, 0, len(*joins.Joins))
	values := make([]interface{}, 0) //  length is unknown
	for _, join := range *joins.Joins {
		joinType, targetName := processJoin(join, tableName, &joinedTablesForSelect)
		if joinType == "" && targetName == "" {
			continue
		}

		if join.Query != nil {
			b := &BaseQueryBuilder{}
			sqQuery, sqValues := b.Build(join.Query)
			targetName = "(" + sqQuery + ")" + " AS " + targetName
			values = append(values, sqValues...)
		}

		var joinBuilder strings.Builder
		joinBuilder.WriteString(joinType)
		joinBuilder.WriteString(" JOIN ")
		joinBuilder.WriteString(targetName)
		joinBuilder.WriteString(" ON ")
		joinBuilder.WriteString(join.SearchColumn)
		joinBuilder.WriteString(" ")
		joinBuilder.WriteString(join.SearchCondition)
		joinBuilder.WriteString(" ")
		joinBuilder.WriteString(join.SearchTargetColumn)

		if _, ok := join.TargetNameMap[consts.Join_CROSS]; ok {
			joinBuilder.WriteString(joinType)
			joinBuilder.WriteString(" JOIN ")
			joinBuilder.WriteString(targetName)
		}

		joinStrings = append(joinStrings, joinBuilder.String())
	}

	return &joinedTablesForSelect, joinStrings, values
}

func processJoin(join structs.Join, tableName string, joinedTablesForSelect *[]structs.Column) (string, string) {
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
		return "", ""
	}

	name := tableName
	if join.Name != "" {
		name = join.Name
	}

	targetNameForSelect := targetName + ".*"
	if !sliceutils.Contains[string](*getNowColNames(joinedTablesForSelect), targetNameForSelect) {
		*joinedTablesForSelect = append(*joinedTablesForSelect, structs.Column{
			Name: targetNameForSelect,
		})
	}
	nameForSelect := name + ".*"
	if !sliceutils.Contains[string](*getNowColNames(joinedTablesForSelect), nameForSelect) {
		*joinedTablesForSelect = append(*joinedTablesForSelect, structs.Column{
			Name: nameForSelect,
		})
	}

	return joinType, targetName //, joinedTablesForSelect
}

// getNowColNames returns the names of the columns in the slice.
func getNowColNames(joinedTablesForSelect *[]structs.Column) *[]string {
	nowColNames := make([]string, len(*joinedTablesForSelect))
	for _, joinedTable := range *joinedTablesForSelect {
		nowColNames = append(nowColNames, joinedTable.Name)
	}
	return &nowColNames
}
