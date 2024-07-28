package query

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
)

type JoinBuilder struct {
	Table *structs.Table
	Joins *structs.Joins
	WhereBuilder[JoinBuilder]
	joinValues []interface{}
}

func NewJoinBuilder(dbBuilder db.QueryBuilderStrategy, cache cache.Cache) *JoinBuilder {
	return &JoinBuilder{
		Table: &structs.Table{},
		Joins: &structs.Joins{
			Joins: &[]structs.Join{},
		},
		WhereBuilder: *NewWhereBuilder[JoinBuilder](dbBuilder, cache),
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

func (b *JoinBuilder) JoinQuery(table string, fn func(j *JoinClauseBuilder) *JoinClauseBuilder) *JoinBuilder {
	jq := fn(NewJoinClauseBuilder())

	jq.JoinClause.Name = table
	jq.JoinClause.TargetNameMap = map[string]string{
		consts.Join_INNER: table,
	}

	b.Joins.JoinClause = jq.JoinClause

	return b
}

func (b *JoinBuilder) LeftJoinQuery(table string, fn func(j *JoinClauseBuilder) *JoinClauseBuilder) *JoinBuilder {

	jq := fn(NewJoinClauseBuilder())

	jq.JoinClause.Name = table
	jq.JoinClause.TargetNameMap = map[string]string{
		consts.Join_LEFT: table,
	}

	b.Joins.JoinClause = jq.JoinClause

	return b
}

func (b *JoinBuilder) RightJoinQuery(table string, fn func(j *JoinClauseBuilder) *JoinClauseBuilder) *JoinBuilder {

	jq := fn(NewJoinClauseBuilder())

	jq.JoinClause.Name = table
	jq.JoinClause.TargetNameMap = map[string]string{
		consts.Join_RIGHT: table,
	}

	b.Joins.JoinClause = jq.JoinClause

	return b
}

func (b *JoinBuilder) JoinSub(q *Builder, alias, my, condition, target string) *JoinBuilder {
	b.joinSubCommon(consts.Join_INNER, q, alias, my, condition, target)
	return b
}

func (b *JoinBuilder) LeftJoinSub(q *Builder, alias, my, condition, target string) *JoinBuilder {
	b.joinSubCommon(consts.Join_LEFT, q, alias, my, condition, target)
	return b
}

func (b *JoinBuilder) RightJoinSub(q *Builder, alias, my, condition, target string) *JoinBuilder {
	b.joinSubCommon(consts.Join_RIGHT, q, alias, my, condition, target)
	return b
}

func (b *JoinBuilder) joinSubCommon(joinType string, q *Builder, alias, my, condition, target string) *JoinBuilder {

	*q.WhereBuilder.query.ConditionGroups = append(*q.WhereBuilder.query.ConditionGroups, structs.WhereGroup{
		Conditions:   *q.WhereBuilder.query.Conditions,
		IsDummyGroup: true,
	})

	sq := &structs.Query{
		ConditionGroups: q.WhereBuilder.query.ConditionGroups,
		Table:           structs.Table{Name: q.selectQuery.Table},
		Columns:         q.selectQuery.Columns,
		Joins:           q.joinBuilder.Joins,
		Order:           q.orderByBuilder.Order,
	}

	myTable := b.Table.Name
	args := &structs.Join{
		Name: myTable,
		TargetNameMap: map[string]string{
			joinType: alias,
		},
		SearchColumn:       my,
		SearchCondition:    condition,
		SearchTargetColumn: target,
		Query:              sq,
	}

	// todo: use cache
	_, value := b.WhereBuilder.BuildSq(sq)

	*b.Joins.Joins = append(*b.Joins.Joins, *args)
	b.joinValues = append(b.joinValues, value...)
	return b
}

func (b *JoinBuilder) JoinLateral(q *Builder, alias string) *JoinBuilder {
	b.joinLateralCommon(consts.Join_LATERAL, q, alias)
	return b
}

func (b *JoinBuilder) LeftJoinLateral(q *Builder, alias string) *JoinBuilder {
	b.joinLateralCommon(consts.Join_LEFT_LATERAL, q, alias)
	return b
}

func (b *JoinBuilder) joinLateralCommon(joinType string, q *Builder, alias string) *JoinBuilder {

	*q.WhereBuilder.query.ConditionGroups = append(*q.WhereBuilder.query.ConditionGroups, structs.WhereGroup{
		Conditions:   *q.WhereBuilder.query.Conditions,
		IsDummyGroup: true,
	})

	sq := &structs.Query{
		ConditionGroups: q.WhereBuilder.query.ConditionGroups,
		Table:           structs.Table{Name: q.selectQuery.Table},
		Columns:         q.selectQuery.Columns,
		Joins:           q.joinBuilder.Joins,
		Order:           q.orderByBuilder.Order,
	}

	myTable := b.Table.Name
	args := &structs.Join{
		Name: myTable,
		TargetNameMap: map[string]string{
			joinType: alias,
		},
		Query: sq,
	}

	// todo: use cache
	_, value := b.WhereBuilder.BuildSq(sq)

	*b.Joins.Joins = append(*b.Joins.Joins, *args)
	b.joinValues = append(b.joinValues, value...)
	return b
}
