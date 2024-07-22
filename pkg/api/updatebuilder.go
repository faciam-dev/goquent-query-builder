package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type UpdateBuilder struct {
	builder             *query.UpdateBuilder
	whereQueryBuilder   *WhereQueryBuilder
	joinQueryBuilder    *JoinQueryBuilder
	orderByQueryBuilder *OrderByQueryBuilder
}

func NewUpdateBuilder(strategy db.QueryBuilderStrategy, cache cache.Cache) *UpdateBuilder {
	return &UpdateBuilder{
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
func (ub *UpdateBuilder) Update(data map[string]interface{}) *UpdateBuilder {
	ub.builder.Update(data)

	return ub
}

// Table
func (ub *UpdateBuilder) Table(table string) *UpdateBuilder {
	ub.builder.Table(table)
	return ub
}

// Where
func (ub *UpdateBuilder) Where(column string, condition string, value interface{}) *UpdateBuilder {
	ub.whereQueryBuilder.Where(column, condition, value)

	return ub
}

// OrWhere
func (ub *UpdateBuilder) OrWhere(column string, condition string, value interface{}) *UpdateBuilder {
	ub.whereQueryBuilder.OrWhere(column, condition, value)

	return ub
}

// WhereQuery
func (ub *UpdateBuilder) WhereQuery(column string, condition string, q *SelectBuilder) *UpdateBuilder {
	ub.whereQueryBuilder.WhereQuery(column, condition, q)

	return ub
}

// OrWhereQuery
func (ub *UpdateBuilder) OrWhereQuery(column string, condition string, q *SelectBuilder) *UpdateBuilder {
	ub.whereQueryBuilder.OrWhereQuery(column, condition, q)

	return ub
}

// WhereGroup
func (ub *UpdateBuilder) WhereGroup(fn func(wb *query.WhereBuilder) *query.WhereBuilder) *UpdateBuilder {
	ub.whereQueryBuilder.WhereGroup(fn)

	return ub
}

// OrWhereGroup
func (ub *UpdateBuilder) OrWhereGroup(fn func(qb *query.WhereBuilder) *query.WhereBuilder) *UpdateBuilder {
	ub.whereQueryBuilder.OrWhereGroup(fn)

	return ub
}

// Join
func (qb *UpdateBuilder) Join(table, my, condition, target string) *UpdateBuilder {
	qb.joinQueryBuilder.Join(table, my, condition, target)
	return qb
}

func (qb *UpdateBuilder) LeftJoin(table, my, condition, target string) *UpdateBuilder {
	qb.joinQueryBuilder.LeftJoin(table, my, condition, target)
	return qb
}

func (qb *UpdateBuilder) RightJoin(table, my, condition, target string) *UpdateBuilder {
	qb.joinQueryBuilder.RightJoin(table, my, condition, target)
	return qb
}

func (qb *UpdateBuilder) CrossJoin(table, my, condition, target string) *UpdateBuilder {
	qb.joinQueryBuilder.CrossJoin(table)
	return qb
}

// OrderBy

func (qb *UpdateBuilder) OrderBy(column, ascDesc string) *UpdateBuilder {
	qb.orderByQueryBuilder.OrderBy(column, ascDesc)
	return qb
}

func (qb *UpdateBuilder) OrderByRaw(raw string) *UpdateBuilder {
	qb.orderByQueryBuilder.OrderByRaw(raw)
	return qb
}

func (qb *UpdateBuilder) ReOrder() *UpdateBuilder {
	qb.orderByQueryBuilder.ReOrder()
	return qb
}

// Build
func (ub *UpdateBuilder) Build() (string, []interface{}) {
	ub.builder.SetWhereBuilder(ub.whereQueryBuilder.builder)
	ub.builder.SetJoinBuilder(ub.joinQueryBuilder.builder)
	ub.builder.SetOrderByBuilder(ub.orderByQueryBuilder.builder)

	return ub.builder.Build()
}
