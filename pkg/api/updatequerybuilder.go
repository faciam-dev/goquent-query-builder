package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type UpdateQueryBuilder struct {
	builder             *query.UpdateBuilder
	whereQueryBuilder   *WhereQueryBuilder
	joinQueryBuilder    *JoinQueryBuilder
	orderByQueryBuilder *OrderByQueryBuilder
}

func NewUpdateQueryBuilder(strategy db.QueryBuilderStrategy, cache *cache.AsyncQueryCache) *UpdateQueryBuilder {
	return &UpdateQueryBuilder{
		builder: query.NewUpdateBuilder(strategy, cache),
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

// Update
func (ub *UpdateQueryBuilder) Update(data map[string]interface{}) *UpdateQueryBuilder {
	ub.builder.Update(data)

	return ub
}

// Table
func (ub *UpdateQueryBuilder) Table(table string) *UpdateQueryBuilder {
	ub.builder.Table(table)
	return ub
}

// Where
func (ub *UpdateQueryBuilder) Where(column string, condition string, value interface{}) *UpdateQueryBuilder {
	ub.whereQueryBuilder.Where(column, condition, value)

	return ub
}

// OrWhere
func (ub *UpdateQueryBuilder) OrWhere(column string, condition string, value interface{}) *UpdateQueryBuilder {
	ub.whereQueryBuilder.OrWhere(column, condition, value)

	return ub
}

// WhereQuery
func (ub *UpdateQueryBuilder) WhereQuery(column string, condition string, q *QueryBuilder) *UpdateQueryBuilder {
	ub.whereQueryBuilder.WhereQuery(column, condition, q)

	return ub
}

// OrWhereQuery
func (ub *UpdateQueryBuilder) OrWhereQuery(column string, condition string, q *QueryBuilder) *UpdateQueryBuilder {
	ub.whereQueryBuilder.OrWhereQuery(column, condition, q)

	return ub
}

// WhereGroup
func (ub *UpdateQueryBuilder) WhereGroup(fn func(wb *query.WhereBuilder) *query.WhereBuilder) *UpdateQueryBuilder {
	ub.whereQueryBuilder.WhereGroup(fn)

	return ub
}

// OrWhereGroup
func (ub *UpdateQueryBuilder) OrWhereGroup(fn func(qb *query.WhereBuilder) *query.WhereBuilder) *UpdateQueryBuilder {
	ub.whereQueryBuilder.OrWhereGroup(fn)

	return ub
}

// Join
func (qb *UpdateQueryBuilder) Join(table, my, condition, target string) *UpdateQueryBuilder {
	qb.joinQueryBuilder.Join(table, my, condition, target)
	return qb
}

func (qb *UpdateQueryBuilder) LeftJoin(table, my, condition, target string) *UpdateQueryBuilder {
	qb.joinQueryBuilder.LeftJoin(table, my, condition, target)
	return qb
}

func (qb *UpdateQueryBuilder) RightJoin(table, my, condition, target string) *UpdateQueryBuilder {
	qb.joinQueryBuilder.RightJoin(table, my, condition, target)
	return qb
}

func (qb *UpdateQueryBuilder) CrossJoin(table, my, condition, target string) *UpdateQueryBuilder {
	qb.joinQueryBuilder.CrossJoin(table)
	return qb
}

// OrderBy

func (qb *UpdateQueryBuilder) OrderBy(column, ascDesc string) *UpdateQueryBuilder {
	qb.orderByQueryBuilder.OrderBy(column, ascDesc)
	return qb
}

func (qb *UpdateQueryBuilder) OrderByRaw(raw string) *UpdateQueryBuilder {
	qb.orderByQueryBuilder.OrderByRaw(raw)
	return qb
}

func (qb *UpdateQueryBuilder) ReOrder() *UpdateQueryBuilder {
	qb.orderByQueryBuilder.ReOrder()
	return qb
}

// Build
func (ub *UpdateQueryBuilder) Build() (string, []interface{}) {
	ub.builder.SetWhereBuilder(ub.whereQueryBuilder.builder)
	ub.builder.SetJoinBuilder(ub.joinQueryBuilder.builder)
	ub.builder.SetOrderByBuilder(ub.orderByQueryBuilder.builder)

	return ub.builder.Build()
}
