package api

import (
	"github.com/faciam-dev/goquent-query-builder/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type OrderByQueryBuilder[T QueryBuilderStrategy[T, C], C any] struct {
	builder *query.OrderByBuilder[C]
	parent  *T
}

func NewOrderByQueryBuilder[T QueryBuilderStrategy[T, C], C any](strategy interfaces.QueryBuilderStrategy, cache cache.Cache) *OrderByQueryBuilder[T, C] {
	return &OrderByQueryBuilder[T, C]{
		builder: query.NewOrderByBuilder[C](strategy, cache),
	}
}

func (qb *OrderByQueryBuilder[T, C]) SetParent(parent *T) *T {
	qb.parent = parent

	return qb.parent
}

func (qb *OrderByQueryBuilder[T, C]) OrderBy(column, ascDesc string) T {
	(*qb.parent).GetOrderByBuilder().OrderBy(column, ascDesc)
	return (*qb.parent).GetQueryBuilder()
}

func (qb *OrderByQueryBuilder[T, C]) OrderByRaw(raw string) T {
	(*qb.parent).GetOrderByBuilder().OrderByRaw(raw)
	return (*qb.parent).GetQueryBuilder()
}

func (qb *OrderByQueryBuilder[T, C]) ReOrder() T {
	(*qb.parent).GetOrderByBuilder().ReOrder()
	return (*qb.parent).GetQueryBuilder()
}
