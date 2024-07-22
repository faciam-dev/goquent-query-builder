package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type DeleteBuilder struct {
	builder             *query.DeleteBuilder
	whereQueryBuilder   *WhereQueryBuilder
	joinQueryBuilder    *JoinQueryBuilder
	orderByQueryBuilder *OrderByQueryBuilder
}

func NewDeleteBuilder(strategy db.QueryBuilderStrategy, cache cache.Cache) *DeleteBuilder {
	return &DeleteBuilder{
		builder: query.NewDeleteBuilder(strategy, cache),
		whereQueryBuilder: &WhereQueryBuilder{
			builder: query.NewWhereBuilder(strategy, cache),
		},
		joinQueryBuilder: &JoinQueryBuilder{
			builder: query.NewJoinBuilder(strategy, cache),
		},
		orderByQueryBuilder: &OrderByQueryBuilder{
			builder: query.NewOrderByBuilder(&[]structs.Order{}),
		},
	}
}

func (qb *DeleteBuilder) Delete() *DeleteBuilder {
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

func (qb *DeleteBuilder) Table(table string) *DeleteBuilder {
	qb.builder.Table(table)
	return qb
}

// Where
func (ub *DeleteBuilder) Where(column string, condition string, value interface{}) *DeleteBuilder {
	ub.whereQueryBuilder.Where(column, condition, value)

	return ub
}

// OrWhere
func (ub *DeleteBuilder) OrWhere(column string, condition string, value interface{}) *DeleteBuilder {
	ub.whereQueryBuilder.OrWhere(column, condition, value)

	return ub
}

// WhereQuery
func (ub *DeleteBuilder) WhereQuery(column string, condition string, q *QueryBuilder) *DeleteBuilder {
	ub.whereQueryBuilder.WhereQuery(column, condition, q)

	return ub
}

// OrWhereQuery
func (ub *DeleteBuilder) OrWhereQuery(column string, condition string, q *QueryBuilder) *DeleteBuilder {
	ub.whereQueryBuilder.OrWhereQuery(column, condition, q)

	return ub
}

// WhereGroup
func (ub *DeleteBuilder) WhereGroup(fn func(wb *query.WhereBuilder) *query.WhereBuilder) *DeleteBuilder {
	ub.whereQueryBuilder.WhereGroup(fn)

	return ub
}

// OrWhereGroup
func (ub *DeleteBuilder) OrWhereGroup(fn func(qb *query.WhereBuilder) *query.WhereBuilder) *DeleteBuilder {
	ub.whereQueryBuilder.OrWhereGroup(fn)

	return ub
}

// Join
func (qb *DeleteBuilder) Join(table, my, condition, target string) *DeleteBuilder {
	qb.joinQueryBuilder.Join(table, my, condition, target)
	return qb
}

func (qb *DeleteBuilder) LeftJoin(table, my, condition, target string) *DeleteBuilder {
	qb.joinQueryBuilder.LeftJoin(table, my, condition, target)
	return qb
}

func (qb *DeleteBuilder) RightJoin(table, my, condition, target string) *DeleteBuilder {
	qb.joinQueryBuilder.RightJoin(table, my, condition, target)
	return qb
}

func (qb *DeleteBuilder) CrossJoin(table, my, condition, target string) *DeleteBuilder {
	qb.joinQueryBuilder.CrossJoin(table)
	return qb
}

// OrderBy

func (qb *DeleteBuilder) OrderBy(column, ascDesc string) *DeleteBuilder {
	qb.orderByQueryBuilder.OrderBy(column, ascDesc)
	return qb
}

func (qb *DeleteBuilder) OrderByRaw(raw string) *DeleteBuilder {
	qb.orderByQueryBuilder.OrderByRaw(raw)
	return qb
}

func (qb *DeleteBuilder) ReOrder() *DeleteBuilder {
	qb.orderByQueryBuilder.ReOrder()
	return qb
}

func (qb *DeleteBuilder) Build() (string, []interface{}) {
	qb.builder.SetWhereBuilder(qb.whereQueryBuilder.builder)
	qb.builder.SetJoinBuilder(qb.joinQueryBuilder.builder)
	qb.builder.SetOrderByBuilder(qb.orderByQueryBuilder.builder)

	return qb.builder.Build()
}
