package api

import "github.com/faciam-dev/goquent-query-builder/internal/query"

type OrderByQueryBuilder struct {
	builder *query.OrderByBuilder
}

func (qb *OrderByQueryBuilder) OrderBy(column, ascDesc string) *OrderByQueryBuilder {
	qb.builder.OrderBy(column, ascDesc)
	return qb
}

func (qb *OrderByQueryBuilder) OrderByRaw(raw string) *OrderByQueryBuilder {
	qb.builder.OrderByRaw(raw)
	return qb
}

func (qb *OrderByQueryBuilder) ReOrder() *OrderByQueryBuilder {
	qb.builder.ReOrder()
	return qb
}
