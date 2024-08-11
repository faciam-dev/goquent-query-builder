package base

import (
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
func (jb *JoinBaseBuilder) Join(sb *[]byte, joins *structs.Joins) []interface{} {
	if jb.columnNames == nil {
		jb.columnNames = &[]string{}
	}

	values := jb.buildJoinStatement(sb, joins)

	return values
}

// buildJoinStatement builds the JOIN statement.
func (jb *JoinBaseBuilder) buildJoinStatement(sb *[]byte, joins *structs.Joins) []interface{} {
	values := make([]interface{}, 0) //  length is unknown
	if joins == nil {
		return []interface{}{}
	}
	if joins.JoinClauses != nil {
		for i := range *joins.JoinClauses {
			vals := []interface{}{}
			j := &structs.Join{
				TargetNameMap: (*joins.JoinClauses)[i].TargetNameMap,
				Name:          (*joins.JoinClauses)[i].Name,
			}

			joinType, targetName := jb.processJoin(j)

			*sb = append(*sb, " "...)
			*sb = append(*sb, joinType...)
			*sb = append(*sb, " JOIN "...)
			if (*joins.JoinClauses)[i].Query != nil {
				*sb = append(*sb, "("...)
				b := jb.u.GetQueryBuilderStrategy()
				sqValues := b.Build(sb, (*joins.JoinClauses)[i].Query, 0, nil)
				values = append(values, sqValues...)
				*sb = append(*sb, ") as "...)
				*sb = jb.u.EscapeIdentifier2(*sb, targetName)
			} else {
				*sb = jb.u.EscapeIdentifier2(*sb, targetName)
			}
			*sb = append(*sb, " ON "...)

			op := ""
			for i, on := range *(*joins.JoinClauses)[i].On {
				if i > 0 {
					if on.Operator == consts.LogicalOperator_OR {
						op = " OR "
					} else {
						op = " AND "
					}
				}

				*sb = append(*sb, op...)
				*sb = jb.u.EscapeIdentifier2(*sb, on.Column)
				*sb = append(*sb, " "...)
				*sb = append(*sb, on.Condition...)
				if on.Value != "" {
					*sb = append(*sb, " "...)
					*sb = jb.u.EscapeIdentifier2(*sb, on.Value.(string)) // TODO: check if this is correct
				}
			}

			op = ""
			for i, condition := range *(*joins.JoinClauses)[i].Conditions {
				if i > 0 || len(*(*joins.JoinClauses)[i].On) > 0 {
					if condition.Operator == consts.LogicalOperator_OR {
						op = " OR "
					} else {
						op = " AND "
					}
				}
				*sb = append(*sb, op...)
				*sb = jb.u.EscapeIdentifier2(*sb, condition.Column)
				*sb = append(*sb, " "...)
				*sb = append(*sb, condition.Condition...)
				*sb = append(*sb, " "+jb.u.GetPlaceholder()...) // TODO: check if this is correct
				vals = append(vals, condition.Value...)
			}
			values = append(values, vals...)
		}
	}

	if joins.Joins == nil {
		return values
	}

	var sortedJoins []structs.Join
	if len(*joins.LateralJoins) == 0 {
		sortedJoins = *joins.Joins
	} else {
		sortedJoins = append(*joins.LateralJoins, *joins.Joins...)
	}
	for i := range sortedJoins {
		joinType, targetName := jb.processJoin(&sortedJoins[i])
		if targetName == "" {
			continue
		}

		if joinType == "" {
			continue
		}

		if _, ok := sortedJoins[i].TargetNameMap[consts.Join_LATERAL]; ok {
			*sb = append(*sb, " ,"...)
			*sb = append(*sb, joinType...)
			if sortedJoins[i].Query != nil {
				*sb = append(*sb, "("...)
				b := jb.u.GetQueryBuilderStrategy()
				sqValues := b.Build(sb, sortedJoins[i].Query, 0, nil)
				*sb = append(*sb, ") as "...)
				*sb = jb.u.EscapeIdentifier2(*sb, targetName)
				values = append(values, sqValues...)
			} else {
				*sb = jb.u.EscapeIdentifier2(*sb, targetName)
			}
		} else if _, ok := sortedJoins[i].TargetNameMap[consts.Join_LEFT_LATERAL]; ok {
			*sb = append(*sb, " ,"...)
			*sb = append(*sb, joinType...)
			if sortedJoins[i].Query != nil {
				*sb = append(*sb, "("...)
				b := jb.u.GetQueryBuilderStrategy()
				sqValues := b.Build(sb, sortedJoins[i].Query, 0, nil)
				*sb = append(*sb, ") as "...)
				*sb = jb.u.EscapeIdentifier2(*sb, targetName)
				values = append(values, sqValues...)
			} else {
				*sb = jb.u.EscapeIdentifier2(*sb, targetName)
			}
		} else if _, ok := sortedJoins[i].TargetNameMap[consts.Join_CROSS]; ok {
			*sb = append(*sb, " "...)
			*sb = append(*sb, joinType...)
			*sb = append(*sb, " JOIN "...)
			*sb = jb.u.EscapeIdentifier2(*sb, targetName)
		} else {
			*sb = append(*sb, " "...)
			*sb = append(*sb, joinType...)
			*sb = append(*sb, " JOIN "...)
			if sortedJoins[i].Query != nil {
				*sb = append(*sb, "("...)
				b := jb.u.GetQueryBuilderStrategy()
				sqValues := b.Build(sb, sortedJoins[i].Query, 0, nil)
				*sb = append(*sb, ") as "...)
				*sb = jb.u.EscapeIdentifier2(*sb, targetName)
				values = append(values, sqValues...)
			} else {
				*sb = jb.u.EscapeIdentifier2(*sb, targetName)
			}
			*sb = append(*sb, " ON "...)
			*sb = jb.u.EscapeIdentifier2(*sb, sortedJoins[i].SearchColumn)
			*sb = append(*sb, " "...)
			*sb = append(*sb, sortedJoins[i].SearchCondition...)
			*sb = append(*sb, " "...)
			*sb = jb.u.EscapeIdentifier2(*sb, sortedJoins[i].SearchTargetColumn)
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
