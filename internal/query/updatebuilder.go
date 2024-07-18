package query

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
)

type UpdateBuilder struct {
	dbBuilder    db.QueryBuilderStrategy
	cache        *cache.AsyncQueryCache
	query        *structs.UpdateQuery
	whereBuilder WhereBuilder
	joinBuilder  JoinBuilder
}

func NewUpdateBuilder(strategy db.QueryBuilderStrategy, cache *cache.AsyncQueryCache) *UpdateBuilder {
	return &UpdateBuilder{
		dbBuilder: strategy,
		cache:     cache,
		query: &structs.UpdateQuery{
			SelectQuery: &structs.Query{},
		},
		whereBuilder: WhereBuilder{
			dbBuilder: strategy,
			cache:     cache,
			query: &structs.Query{
				ConditionGroups: &[]structs.WhereGroup{},
				Conditions:      &[]structs.Where{},
			},
		},
		joinBuilder: JoinBuilder{
			Table: &structs.Table{},
			Joins: &[]structs.Join{},
		},
	}
}

func (b *UpdateBuilder) Table(table string) *UpdateBuilder {
	b.query.Table = table
	b.joinBuilder.Table.Name = table
	return b
}

func (b *UpdateBuilder) Where(column string, condition string, value ...interface{}) *UpdateBuilder {
	b.whereBuilder.Where(column, condition, value...)
	return b
}

func (b *UpdateBuilder) OrWhere(column string, condition string, value ...interface{}) *UpdateBuilder {
	b.whereBuilder.OrWhere(column, condition, value...)
	return b
}

func (b *UpdateBuilder) WhereQuery(column string, condition string, q *Builder) *UpdateBuilder {
	b.whereBuilder.WhereQuery(column, condition, q)

	return b
}

func (b *UpdateBuilder) OrWhereQuery(column string, condition string, q *Builder) *UpdateBuilder {
	b.whereBuilder.OrWhereQuery(column, condition, q)

	return b
}

func (b *UpdateBuilder) WhereGroup(fn func(b *WhereBuilder) *WhereBuilder) *UpdateBuilder {
	b.whereBuilder.WhereGroup(fn)

	return b
}

func (b *UpdateBuilder) OrWhereGroup(fn func(b *WhereBuilder) *WhereBuilder) *UpdateBuilder {
	b.whereBuilder.OrWhereGroup(fn)

	return b
}

func (b *UpdateBuilder) Update(data map[string]interface{}) *UpdateBuilder {
	b.query.Values = data

	// If there are conditions, add them to the query
	if len(*b.whereBuilder.query.Conditions) > 0 {
		*b.whereBuilder.query.ConditionGroups = append(*b.whereBuilder.query.ConditionGroups, structs.WhereGroup{
			Conditions:   *b.whereBuilder.query.Conditions,
			Operator:     consts.LogicalOperator_AND,
			IsDummyGroup: true,
		})
		b.whereBuilder.query.Conditions = &[]structs.Where{}
	}

	b.query.SelectQuery.Conditions = b.whereBuilder.query.Conditions
	b.query.SelectQuery.ConditionGroups = b.whereBuilder.query.ConditionGroups
	b.query.SelectQuery.Joins = b.joinBuilder.Joins

	return b
}

func (u *UpdateBuilder) Build() (string, []interface{}) {
	query, values := u.dbBuilder.BuildUpdate(u.query)
	return query, values
}

func (b *UpdateBuilder) Join(table, my, condition, target string) *UpdateBuilder {
	b.joinBuilder.Join(table, my, condition, target)
	return b
}

func (b *UpdateBuilder) LeftJoin(table, my, condition, target string) *UpdateBuilder {
	b.joinBuilder.LeftJoin(table, my, condition, target)
	return b
}

func (b *UpdateBuilder) RightJoin(table, my, condition, target string) *UpdateBuilder {
	b.joinBuilder.RightJoin(table, my, condition, target)
	return b
}

func (b *UpdateBuilder) CrossJoin(table string) *UpdateBuilder {
	b.joinBuilder.CrossJoin(table)
	return b
}
