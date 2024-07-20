package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type DeleteQueryBuilder struct {
	builder             *query.DeleteBuilder
	whereQueryBuilder   *WhereQueryBuilder
	joinQueryBuilder    *JoinQueryBuilder
	orderByQueryBuilder *OrderByQueryBuilder
}

func NewDeleteQueryBuilder(strategy db.QueryBuilderStrategy, cache *cache.AsyncQueryCache) *DeleteQueryBuilder {
	return &DeleteQueryBuilder{
		builder: query.NewDeleteBuilder(strategy, cache),
		whereQueryBuilder: &WhereQueryBuilder{
			builder: query.NewWhereBuilder(strategy, cache),
		},
		joinQueryBuilder: &JoinQueryBuilder{
			builder: query.NewJoinBuilder(&structs.Joins{
				Joins: &[]structs.Join{},
			}),
		},
		orderByQueryBuilder: &OrderByQueryBuilder{
			builder: query.NewOrderByBuilder(&[]structs.Order{}),
		},
	}
}

func (qb *DeleteQueryBuilder) Delete() *DeleteQueryBuilder {
	qb.builder.Delete()

	return qb
}

// Using
/*
func (ub *UpdateQueryBuilder) Using(qb *QueryBuilder) *UpdateQueryBuilder {
	ub.builder.Using(qb)

	return ub
}
*/

func (qb *DeleteQueryBuilder) Table(table string) *DeleteQueryBuilder {
	qb.builder.Table(table)
	return qb
}

// Where
func (ub *DeleteQueryBuilder) Where(column string, condition string, value interface{}) *DeleteQueryBuilder {
	ub.whereQueryBuilder.Where(column, condition, value)

	return ub
}

// OrWhere
func (ub *DeleteQueryBuilder) OrWhere(column string, condition string, value interface{}) *DeleteQueryBuilder {
	ub.whereQueryBuilder.OrWhere(column, condition, value)

	return ub
}

// WhereQuery
func (ub *DeleteQueryBuilder) WhereQuery(column string, condition string, q *QueryBuilder) *DeleteQueryBuilder {
	ub.whereQueryBuilder.WhereQuery(column, condition, q)

	return ub
}

// OrWhereQuery
func (ub *DeleteQueryBuilder) OrWhereQuery(column string, condition string, q *QueryBuilder) *DeleteQueryBuilder {
	ub.whereQueryBuilder.OrWhereQuery(column, condition, q)

	return ub
}

// WhereGroup
func (ub *DeleteQueryBuilder) WhereGroup(fn func(wb *query.WhereBuilder) *query.WhereBuilder) *DeleteQueryBuilder {
	ub.whereQueryBuilder.WhereGroup(fn)

	return ub
}

// OrWhereGroup
func (ub *DeleteQueryBuilder) OrWhereGroup(fn func(qb *query.WhereBuilder) *query.WhereBuilder) *DeleteQueryBuilder {
	ub.whereQueryBuilder.OrWhereGroup(fn)

	return ub
}

// Join
func (qb *DeleteQueryBuilder) Join(table, my, condition, target string) *DeleteQueryBuilder {
	qb.joinQueryBuilder.Join(table, my, condition, target)
	return qb
}

func (qb *DeleteQueryBuilder) LeftJoin(table, my, condition, target string) *DeleteQueryBuilder {
	qb.joinQueryBuilder.LeftJoin(table, my, condition, target)
	return qb
}

func (qb *DeleteQueryBuilder) RightJoin(table, my, condition, target string) *DeleteQueryBuilder {
	qb.joinQueryBuilder.RightJoin(table, my, condition, target)
	return qb
}

func (qb *DeleteQueryBuilder) CrossJoin(table, my, condition, target string) *DeleteQueryBuilder {
	qb.joinQueryBuilder.CrossJoin(table)
	return qb
}

// OrderBy

func (qb *DeleteQueryBuilder) OrderBy(column, ascDesc string) *DeleteQueryBuilder {
	qb.orderByQueryBuilder.OrderBy(column, ascDesc)
	return qb
}

func (qb *DeleteQueryBuilder) OrderByRaw(raw string) *DeleteQueryBuilder {
	qb.orderByQueryBuilder.OrderByRaw(raw)
	return qb
}

func (qb *DeleteQueryBuilder) ReOrder() *DeleteQueryBuilder {
	qb.orderByQueryBuilder.ReOrder()
	return qb
}

func (qb *DeleteQueryBuilder) Build() (string, []interface{}) {
	qb.builder.SetWhereBuilder(qb.whereQueryBuilder.builder)
	qb.builder.SetJoinBuilder(qb.joinQueryBuilder.builder)
	qb.builder.SetOrderByBuilder(qb.orderByQueryBuilder.builder)

	return qb.builder.Build()
}
