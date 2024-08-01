package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type UpdateBuilder struct {
	WhereQueryBuilder[UpdateBuilder, query.UpdateBuilder]
	JoinQueryBuilder[UpdateBuilder, query.UpdateBuilder]
	builder             *query.UpdateBuilder
	orderByQueryBuilder *OrderByQueryBuilder
}

func NewUpdateBuilder(strategy db.QueryBuilderStrategy, cache cache.Cache) *UpdateBuilder {
	ub := &UpdateBuilder{
		builder: query.NewUpdateBuilder(strategy, cache),
		orderByQueryBuilder: &OrderByQueryBuilder{
			builder: query.NewOrderByBuilder(&[]structs.Order{}),
		},
	}

	whereQueryBuilder := NewWhereQueryBuilder[UpdateBuilder, query.UpdateBuilder](strategy, cache)
	whereQueryBuilder.SetParent(ub)
	ub.WhereQueryBuilder = *whereQueryBuilder

	joinQueryBuilder := NewJoinQueryBuilder[UpdateBuilder, query.UpdateBuilder](strategy, cache)
	joinQueryBuilder.SetParent(ub)
	ub.JoinQueryBuilder = *joinQueryBuilder

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
func (ub *UpdateBuilder) Build() (string, []interface{}, error) {
	ub.builder.SetWhereBuilder(*ub.WhereQueryBuilder.builder)
	ub.builder.SetJoinBuilder(*ub.JoinQueryBuilder.builder)
	ub.builder.SetOrderByBuilder(ub.orderByQueryBuilder.builder)

	return ub.builder.Build()
}

func (ub *UpdateBuilder) Dump() (string, []interface{}, error) {
	ub.builder.SetWhereBuilder(*ub.WhereQueryBuilder.builder)
	ub.builder.SetJoinBuilder(*ub.JoinQueryBuilder.builder)
	ub.builder.SetOrderByBuilder(ub.orderByQueryBuilder.builder)

	b := query.NewDebugBuilder[*query.UpdateBuilder, UpdateBuilder](ub.builder)

	return b.Dump()
}

func (ub *UpdateBuilder) RawSql() (string, error) {
	ub.builder.SetWhereBuilder(*ub.WhereQueryBuilder.builder)
	ub.builder.SetJoinBuilder(*ub.JoinQueryBuilder.builder)
	ub.builder.SetOrderByBuilder(ub.orderByQueryBuilder.builder)

	b := query.NewDebugBuilder[*query.UpdateBuilder, UpdateBuilder](ub.builder)

	return b.RawSql()
}
