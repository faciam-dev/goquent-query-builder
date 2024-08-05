package base

import (
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

type JoinBaseBuilder struct {
	u           interfaces.SQLUtils
	join        *structs.Joins
	columnNames *[]string
}

func NewJoinBaseBuilder(util interfaces.SQLUtils, j *structs.Joins) *JoinBaseBuilder {
	return &JoinBaseBuilder{
		u:           util,
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
	values := make([]interface{}, 0) //  length is unknown
	if joins == nil {
		return []interface{}{}
	}
	if joins.JoinClause != nil {
		for _, joinClause := range *joins.JoinClause {
			vals := make([]interface{}, 0, len(*joinClause.Conditions))

			j := &structs.Join{
				TargetNameMap: joinClause.TargetNameMap,
				Name:          joinClause.Name,
			}

			joinType, targetName := jb.processJoin(j)

			/*
				if joinClause.Query != nil {
					b := jb.u.GetQueryBuilderStrategy()
					sqQuery, sqValues := b.Build("", joinClause.Query, 0, nil)
					targetName = "(" + sqQuery + ")" + " as " + jb.u.EscapeIdentifier(sb, targetName)
					values = append(values, sqValues...)
				}
			*/
			sb.WriteString(" ")
			sb.WriteString(joinType)
			sb.WriteString(" JOIN ")
			if joinClause.Query != nil {
				sb.WriteString("(")
				//sb.WriteString(sqQuery)
				b := jb.u.GetQueryBuilderStrategy()
				_, sqValues := b.Build(sb, "", joinClause.Query, 0, nil)
				//targetName = "(" + sqQuery + ")" + " as " + jb.u.EscapeIdentifier(sb, targetName)
				values = append(values, sqValues...)
				//sb.WriteString(targetName)
				sb.WriteString(") as ")
				sb.WriteString(jb.u.EscapeIdentifier(sb, targetName))
			} else {
				sb.WriteString(jb.u.EscapeIdentifier(sb, targetName))
			}
			sb.WriteString(" ON ")

			op := ""
			for i, on := range *joinClause.On {
				if i > 0 {
					if on.Operator == consts.LogicalOperator_OR {
						op = " OR "
					} else {
						op = " AND "
					}
				}

				sb.WriteString(op)
				sb.WriteString(jb.u.EscapeIdentifier(sb, on.Column))
				sb.WriteString(" ")
				sb.WriteString(on.Condition)
				if on.Value != "" {
					sb.WriteString(" ")
					sb.WriteString(jb.u.EscapeIdentifier(sb, on.Value.(string))) // TODO: check if this is correct
				}
			}

			op = ""
			for i, condition := range *joinClause.Conditions {
				if i > 0 || len(*joinClause.On) > 0 {
					if condition.Operator == consts.LogicalOperator_OR {
						op = " OR "
					} else {
						op = " AND "
					}
				}
				sb.WriteString(op)
				sb.WriteString(jb.u.EscapeIdentifier(sb, condition.Column))
				sb.WriteString(" ")
				sb.WriteString(condition.Condition)
				sb.WriteString(" " + jb.u.GetPlaceholder()) // TODO: check if this is correct
				vals = append(vals, condition.Value...)
			}
			values = append(values, vals...)
		}
	}

	if joins.Joins == nil {
		return values
	}

	// sort by lateral joins first
	//	sortedJoins := make([]*structs.Join, 0, len(*joins.Joins))
	lateralJoins := make([]*structs.Join, 0, len(*joins.Joins))
	otherJoins := make([]*structs.Join, 0, len(*joins.Joins))

	for i := range *joins.Joins {
		if _, ok := (*joins.Joins)[i].TargetNameMap[consts.Join_LATERAL]; ok {
			lateralJoins = append(lateralJoins, &(*joins.Joins)[i])
		} else if _, ok := (*joins.Joins)[i].TargetNameMap[consts.Join_LEFT_LATERAL]; ok {
			lateralJoins = append(lateralJoins, &(*joins.Joins)[i])
		} else {
			otherJoins = append(otherJoins, &(*joins.Joins)[i])
		}
	}

	sortedJoins := append(lateralJoins, otherJoins...)
	for i := range sortedJoins {
		joinType, targetName := jb.processJoin(sortedJoins[i])
		if targetName == "" {
			continue
		}

		if joinType == "" {
			continue
		}

		/*
			if sortedJoins[i].Query != nil {
				b := jb.u.GetQueryBuilderStrategy()
				sqQuery, sqValues := b.Build("", sortedJoins[i].Query, 0, nil)
				targetName = "(" + sqQuery + ")" + " as " + jb.u.EscapeIdentifier(sb, targetName)
				values = append(values, sqValues...)
			}
		*/

		if _, ok := sortedJoins[i].TargetNameMap[consts.Join_LATERAL]; ok {
			sb.WriteString(" ,")
			sb.WriteString(joinType)
			if sortedJoins[i].Query != nil {
				sb.WriteString("(")
				b := jb.u.GetQueryBuilderStrategy()
				_, sqValues := b.Build(sb, "", sortedJoins[i].Query, 0, nil)
				sb.WriteString(") as ")
				sb.WriteString(jb.u.EscapeIdentifier(sb, targetName))
				values = append(values, sqValues...)
				//sb.WriteString(targetName)
			} else {
				sb.WriteString(jb.u.EscapeIdentifier(sb, targetName))
			}
		} else if _, ok := sortedJoins[i].TargetNameMap[consts.Join_LEFT_LATERAL]; ok {
			sb.WriteString(" ,")
			sb.WriteString(joinType)
			if sortedJoins[i].Query != nil {
				sb.WriteString("(")
				b := jb.u.GetQueryBuilderStrategy()
				_, sqValues := b.Build(sb, "", sortedJoins[i].Query, 0, nil)
				//sb.WriteString(sqQuery)
				sb.WriteString(") as ")
				sb.WriteString(jb.u.EscapeIdentifier(sb, targetName))
				values = append(values, sqValues...)
				//sb.WriteString(targetName)
			} else {
				sb.WriteString(jb.u.EscapeIdentifier(sb, targetName))
			}
		} else if _, ok := sortedJoins[i].TargetNameMap[consts.Join_CROSS]; ok {
			sb.WriteString(" ")
			sb.WriteString(joinType)
			sb.WriteString(" JOIN ")
			sb.WriteString(jb.u.EscapeIdentifier(sb, targetName))
		} else {
			sb.WriteString(" ")
			sb.WriteString(joinType)
			sb.WriteString(" JOIN ")
			if sortedJoins[i].Query != nil {
				sb.WriteString("(")
				//sb.WriteString(targetName)
				b := jb.u.GetQueryBuilderStrategy()
				_, sqValues := b.Build(sb, "", sortedJoins[i].Query, 0, nil)
				//sb.WriteString(sqQuery)
				sb.WriteString(") as ")
				sb.WriteString(jb.u.EscapeIdentifier(sb, targetName))
				values = append(values, sqValues...)
			} else {
				sb.WriteString(jb.u.EscapeIdentifier(sb, targetName))
			}
			sb.WriteString(" ON ")
			sb.WriteString(jb.u.EscapeIdentifier(sb, sortedJoins[i].SearchColumn))
			sb.WriteString(" ")
			sb.WriteString(sortedJoins[i].SearchCondition)
			sb.WriteString(" ")
			sb.WriteString(jb.u.EscapeIdentifier(sb, sortedJoins[i].SearchTargetColumn))
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
