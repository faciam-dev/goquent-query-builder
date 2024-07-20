package query

import (
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type JoinBuilder struct {
	Table        *structs.Table
	Joins        *structs.Joins
	whereBuilder *WhereBuilder
}

func NewJoinBuilder(j *structs.Joins) *JoinBuilder {
	return &JoinBuilder{
		Table: &structs.Table{},
		Joins: j,
	}
}

// Join adds a JOIN clause.
func (b *JoinBuilder) Join(table string, my string, condition string, target string) *JoinBuilder {
	return b.joinCommon(consts.Join_INNER, table, my, condition, target)
}

// LeftJoin adds a LEFT JOIN clause.
func (b *JoinBuilder) LeftJoin(table string, my string, condition string, target string) *JoinBuilder {
	return b.joinCommon(consts.Join_LEFT, table, my, condition, target)
}

// RightJoin adds a RIGHT JOIN clause.
func (b *JoinBuilder) RightJoin(table string, my string, condition string, target string) *JoinBuilder {
	return b.joinCommon(consts.Join_RIGHT, table, my, condition, target)
}

// joinCommon is a helper function for JOIN, LEFT JOIN, and RIGHT JOIN.
func (b *JoinBuilder) joinCommon(joinType string, table string, my string, condition string, target string) *JoinBuilder {
	myTable := b.Table.Name
	// If a previous JOIN exists, retrieve the table name of that JOIN.
	if b.Joins.Joins != nil && len(*b.Joins.Joins) > 0 {
		myTable = (*b.Joins.Joins)[len(*b.Joins.Joins)-1].Name
	}
	*b.Joins.Joins = append(*b.Joins.Joins, structs.Join{
		Name: myTable,
		TargetNameMap: map[string]string{
			joinType: table,
		},
		SearchColumn:       my,
		SearchCondition:    condition,
		SearchTargetColumn: target,
	})
	return b
}

// CrossJoin adds a CROSS JOIN clause.
func (b *JoinBuilder) CrossJoin(table string) *JoinBuilder {
	myTable := b.Table.Name
	// If a previous JOIN exists, retrieve the table name of that JOIN.
	if b.Joins != nil && len(*b.Joins.Joins) > 0 {
		myTable = (*b.Joins.Joins)[len(*b.Joins.Joins)-1].Name
	}
	*b.Joins.Joins = append(*b.Joins.Joins, structs.Join{
		Name: myTable,
		TargetNameMap: map[string]string{
			consts.Join_CROSS: table,
		},
	})
	return b
}
