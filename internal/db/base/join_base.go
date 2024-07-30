package base

import (
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type JoinBaseBuilder struct {
	join        *structs.Joins
	columnNames *[]string
}

func NewJoinBaseBuilder(j *structs.Joins) *JoinBaseBuilder {
	return &JoinBaseBuilder{
		join:        j,
		columnNames: &[]string{},
	}
}

// Join builds the JOIN query.
func (jb *JoinBaseBuilder) Join(sb *strings.Builder, joins *structs.Joins) []interface{} {
	if jb.columnNames == nil {
		jb.columnNames = &[]string{}
	}

	values := jb.buildJoinStatement(sb, joins)

	return values
}

// buildJoinStatement builds the JOIN statement.
func (jb *JoinBaseBuilder) buildJoinStatement(sb *strings.Builder, joins *structs.Joins) []interface{} {
	if joins == nil {
		return []interface{}{}
	}
	if joins.JoinClause != nil {
		values := make([]interface{}, 0, len(*joins.JoinClause.Conditions))

		j := &structs.Join{
			TargetNameMap: joins.JoinClause.TargetNameMap,
			Name:          joins.JoinClause.Name,
		}

		joinType, targetName := jb.processJoin(j)

		if joins.JoinClause.Query != nil {
			b := &BaseQueryBuilder{}
			sqQuery, sqValues := b.Build("", joins.JoinClause.Query, 0, nil)
			targetName = "(" + sqQuery + ")" + " AS " + targetName
			values = append(values, sqValues...)
		}
		sb.WriteString(" ")
		sb.WriteString(joinType)
		sb.WriteString(" JOIN ")
		sb.WriteString(targetName)
		sb.WriteString(" ON ")

		op := ""
		for i, on := range *joins.JoinClause.On {
			if i > 0 {
				if on.Operator == consts.LogicalOperator_OR {
					op = " OR "
				} else {
					op = " AND "
				}
			}

			sb.WriteString(op)
			sb.WriteString(on.Column)
			sb.WriteString(" ")
			sb.WriteString(on.Condition)
			if on.Value != "" {
				sb.WriteString(" ")
				sb.WriteString(on.Value.(string)) // TODO: check if this is correct
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
			sb.WriteString(op)
			sb.WriteString(condition.Column)
			sb.WriteString(" ")
			sb.WriteString(condition.Condition)
			sb.WriteString(" ?") // TODO: check if this is correct
			values = append(values, condition.Value...)
		}

		return values
	}

	if joins.Joins == nil {
		return []interface{}{}
	}

	values := make([]interface{}, 0) //  length is unknown

	// sort by lateral joins first
	//	sortedJoins := make([]*structs.Join, 0, len(*joins.Joins))
	lateralJoins := make([]*structs.Join, 0, len(*joins.Joins))
	otherJoins := make([]*structs.Join, 0, len(*joins.Joins))

	for _, join := range *joins.Joins {
		if _, ok := join.TargetNameMap[consts.Join_LATERAL]; ok {
			lateralJoins = append(lateralJoins, &join)
		} else if _, ok := join.TargetNameMap[consts.Join_LEFT_LATERAL]; ok {
			lateralJoins = append(lateralJoins, &join)
		} else {
			otherJoins = append(otherJoins, &join)
		}
	}

	//sortedJoins = append(lateralJoins, otherJoins...)

	for _, join := range append(lateralJoins, otherJoins...) {
		joinType, targetName := jb.processJoin(join)
		if targetName == "" {
			continue
		}

		if joinType == "" {
			continue
		}

		if join.Query != nil {
			b := &BaseQueryBuilder{}
			sqQuery, sqValues := b.Build("", join.Query, 0, nil)
			targetName = "(" + sqQuery + ")" + " AS " + targetName
			values = append(values, sqValues...)
		}

		if _, ok := join.TargetNameMap[consts.Join_LATERAL]; ok {
			sb.WriteString(" ,")
			sb.WriteString(joinType)
			sb.WriteString(targetName)
		} else if _, ok := join.TargetNameMap[consts.Join_LEFT_LATERAL]; ok {
			sb.WriteString(" ,")
			sb.WriteString(joinType)
			sb.WriteString(targetName)
		} else if _, ok := join.TargetNameMap[consts.Join_CROSS]; ok {
			sb.WriteString(" ")
			sb.WriteString(joinType)
			sb.WriteString(" JOIN ")
			sb.WriteString(targetName)
		} else {
			sb.WriteString(" ")
			sb.WriteString(joinType)
			sb.WriteString(" JOIN ")
			sb.WriteString(targetName)
			sb.WriteString(" ON ")
			sb.WriteString(join.SearchColumn)
			sb.WriteString(" ")
			sb.WriteString(join.SearchCondition)
			sb.WriteString(" ")
			sb.WriteString(join.SearchTargetColumn)
		}
	}

	return values
}

func (j *JoinBaseBuilder) processJoin(join *structs.Join) (string, string) {
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
	if _, ok := join.TargetNameMap[consts.Join_LATERAL]; ok {
		targetName = join.TargetNameMap[consts.Join_LATERAL]
		joinType = consts.Join_Type_LATERAL
	}
	if _, ok := join.TargetNameMap[consts.Join_LEFT_LATERAL]; ok {
		targetName = join.TargetNameMap[consts.Join_LEFT_LATERAL]
		joinType = consts.Join_Type_LEFT_LATERAL
	}

	if targetName == "" {
		return "", ""
	}

	return joinType, targetName
}
