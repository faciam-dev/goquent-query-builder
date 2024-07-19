package query

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
)

type DeleteBuilder struct {
	dbBuilder      db.QueryBuilderStrategy
	cache          *cache.AsyncQueryCache
	query          *structs.DeleteQuery
	whereBuilder   *WhereBuilder
	joinBuilder    *JoinBuilder
	orderByBuilder *OrderByBuilder
}

func NewDeleteBuilder(strategy db.QueryBuilderStrategy, cache *cache.AsyncQueryCache) *DeleteBuilder {
	return &DeleteBuilder{
		dbBuilder: strategy,
		cache:     cache,
		query: &structs.DeleteQuery{
			Query: &structs.Query{},
		},
		whereBuilder: &WhereBuilder{
			dbBuilder: strategy,
			cache:     cache,
			query: &structs.Query{
				ConditionGroups: &[]structs.WhereGroup{},
				Conditions:      &[]structs.Where{},
			},
		},
		joinBuilder: &JoinBuilder{
			Table: &structs.Table{},
			Joins: &[]structs.Join{},
		},
		orderByBuilder: &OrderByBuilder{
			Order: &[]structs.Order{},
		},
	}
}

func (b *DeleteBuilder) Table(table string) *DeleteBuilder {
	b.query.Table = table
	b.joinBuilder.Table.Name = table
	return b
}

func (b *DeleteBuilder) Where(column string, condition string, value ...interface{}) *DeleteBuilder {
	b.whereBuilder.Where(column, condition, value...)
	return b
}

func (b *DeleteBuilder) OrWhere(column string, condition string, value ...interface{}) *DeleteBuilder {
	b.whereBuilder.OrWhere(column, condition, value...)
	return b
}

func (b *DeleteBuilder) WhereQuery(column string, condition string, q *Builder) *DeleteBuilder {
	b.whereBuilder.WhereQuery(column, condition, q)

	return b
}

func (b *DeleteBuilder) OrWhereQuery(column string, condition string, q *Builder) *DeleteBuilder {
	b.whereBuilder.OrWhereQuery(column, condition, q)

	return b
}

func (b *DeleteBuilder) WhereGroup(fn func(b *WhereBuilder) *WhereBuilder) *DeleteBuilder {
	b.whereBuilder.WhereGroup(fn)

	return b
}

func (b *DeleteBuilder) OrWhereGroup(fn func(b *WhereBuilder) *WhereBuilder) *DeleteBuilder {
	b.whereBuilder.OrWhereGroup(fn)

	return b
}

func (b *DeleteBuilder) Delete() *DeleteBuilder {
	// If there are conditions, add them to the query
	if len(*b.whereBuilder.query.Conditions) > 0 {
		*b.whereBuilder.query.ConditionGroups = append(*b.whereBuilder.query.ConditionGroups, structs.WhereGroup{
			Conditions:   *b.whereBuilder.query.Conditions,
			Operator:     consts.LogicalOperator_AND,
			IsDummyGroup: true,
		})
		b.whereBuilder.query.Conditions = &[]structs.Where{}
	}

	b.query.Query.Conditions = b.whereBuilder.query.Conditions
	b.query.Query.ConditionGroups = b.whereBuilder.query.ConditionGroups
	b.query.Query.Joins = b.joinBuilder.Joins
	b.query.Query.Order = b.orderByBuilder.Order

	return b
}

func (u *DeleteBuilder) Build() (string, []interface{}) {
	query, values := u.dbBuilder.BuildDelete(u.query)
	return query, values
}

func (b *DeleteBuilder) Join(table, my, condition, target string) *DeleteBuilder {
	b.joinBuilder.Join(table, my, condition, target)
	return b
}

func (b *DeleteBuilder) LeftJoin(table, my, condition, target string) *DeleteBuilder {
	b.joinBuilder.LeftJoin(table, my, condition, target)
	return b
}

func (b *DeleteBuilder) RightJoin(table, my, condition, target string) *DeleteBuilder {
	b.joinBuilder.RightJoin(table, my, condition, target)
	return b
}

func (b *DeleteBuilder) CrossJoin(table string) *DeleteBuilder {
	b.joinBuilder.CrossJoin(table)
	return b
}

func (b *DeleteBuilder) OrderBy(column string, direction string) *DeleteBuilder {
	b.orderByBuilder.OrderBy(column, direction)
	return b
}

func (b *DeleteBuilder) OrderByRaw(raw string) *DeleteBuilder {
	b.orderByBuilder.OrderByRaw(raw)
	return b
}

func (b *DeleteBuilder) ReOrder() *DeleteBuilder {
	b.orderByBuilder.ReOrder()
	return b
}

func (b *DeleteBuilder) GetQuery() *structs.DeleteQuery {
	return b.query
}
