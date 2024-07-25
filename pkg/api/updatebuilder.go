package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type UpdateBuilder struct {
	WhereQueryBuilder[UpdateBuilder, query.UpdateBuilder]
	builder             *query.UpdateBuilder
	joinQueryBuilder    *JoinQueryBuilder
	orderByQueryBuilder *OrderByQueryBuilder
}

func NewUpdateBuilder(strategy db.QueryBuilderStrategy, cache cache.Cache) *UpdateBuilder {
	ub := &UpdateBuilder{
		builder: query.NewUpdateBuilder(strategy, cache),
		joinQueryBuilder: &JoinQueryBuilder{
			builder: query.NewJoinBuilder(strategy, cache),
		},
		orderByQueryBuilder: &OrderByQueryBuilder{
			builder: query.NewOrderByBuilder(&[]structs.Order{}),
		},
	}

	whereQueryBuilder := NewWhereQueryBuilder[UpdateBuilder, query.UpdateBuilder](strategy, cache)
	whereQueryBuilder.SetParent(ub)
	ub.WhereQueryBuilder = *whereQueryBuilder

	return ub
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
	ub.builder.SetWhereBuilder(*ub.WhereQueryBuilder.builder)
	//	ub.builder.SetWhereBuilder(ub.WhereQueryBuilder.builder)
	ub.builder.SetJoinBuilder(ub.joinQueryBuilder.builder)
	ub.builder.SetOrderByBuilder(ub.orderByQueryBuilder.builder)

	return ub.builder.Build()
}
