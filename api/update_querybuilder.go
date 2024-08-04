package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type UpdateQueryBuilder struct {
	WhereQueryBuilder[UpdateQueryBuilder, query.UpdateBuilder]
	JoinQueryBuilder[UpdateQueryBuilder, query.UpdateBuilder]
	builder             *query.UpdateBuilder
	orderByQueryBuilder *OrderByQueryBuilder
}

func NewUpdateQueryBuilder(strategy interfaces.QueryBuilderStrategy, cache cache.Cache) *UpdateQueryBuilder {
	ub := &UpdateQueryBuilder{
		builder: query.NewUpdateBuilder(strategy, cache),
		orderByQueryBuilder: &OrderByQueryBuilder{
			builder: query.NewOrderByBuilder(&[]structs.Order{}),
		},
	}

	whereQueryBuilder := NewWhereQueryBuilder[UpdateQueryBuilder, query.UpdateBuilder](strategy, cache)
	whereQueryBuilder.SetParent(ub)
	ub.WhereQueryBuilder = *whereQueryBuilder

	joinQueryBuilder := NewJoinQueryBuilder[UpdateQueryBuilder, query.UpdateBuilder](strategy, cache)
	joinQueryBuilder.SetParent(ub)
	ub.JoinQueryBuilder = *joinQueryBuilder

	return ub
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
func (ub *UpdateQueryBuilder) Build() (string, []interface{}, error) {
	ub.builder.SetWhereBuilder(*ub.WhereQueryBuilder.builder)
	ub.builder.SetJoinBuilder(*ub.JoinQueryBuilder.builder)
	ub.builder.SetOrderByBuilder(ub.orderByQueryBuilder.builder)

	return ub.builder.Build()
}

func (ub *UpdateQueryBuilder) Dump() (string, []interface{}, error) {
	ub.builder.SetWhereBuilder(*ub.WhereQueryBuilder.builder)
	ub.builder.SetJoinBuilder(*ub.JoinQueryBuilder.builder)
	ub.builder.SetOrderByBuilder(ub.orderByQueryBuilder.builder)

	b := query.NewDebugBuilder[*query.UpdateBuilder, UpdateQueryBuilder](ub.builder)

	return b.Dump()
}

func (ub *UpdateQueryBuilder) RawSql() (string, error) {
	ub.builder.SetWhereBuilder(*ub.WhereQueryBuilder.builder)
	ub.builder.SetJoinBuilder(*ub.JoinQueryBuilder.builder)
	ub.builder.SetOrderByBuilder(ub.orderByQueryBuilder.builder)

	b := query.NewDebugBuilder[*query.UpdateBuilder, UpdateQueryBuilder](ub.builder)

	return b.RawSql()
}
