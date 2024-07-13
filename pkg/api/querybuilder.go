package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type QueryBuilder struct {
	builder *query.Builder
}

func NewQueryBuilder(strategy db.QueryBuilderStrategy, cache *cache.AsyncQueryCache) *QueryBuilder {
	return &QueryBuilder{
		builder: query.NewBuilder(strategy, cache),
	}
}

func (qb *QueryBuilder) Table(table string) *QueryBuilder {
	qb.builder.Table(table)
	return qb
}

func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.builder.Select(columns...)
	return qb
}

func (qb *QueryBuilder) Where(column string, condition string, value interface{}) *QueryBuilder {
	qb.builder.Where(column, condition, value)
	return qb
}

func (qb *QueryBuilder) Join(table, my, condition, target string) *QueryBuilder {
	qb.builder.Join(table, my, condition, target)
	return qb
}

func (qb *QueryBuilder) OrderBy(column, ascDesc string) *QueryBuilder {
	qb.builder.OrderBy(column, ascDesc)
	return qb
}

func (qb *QueryBuilder) Build() (string, []interface{}) {
	return qb.builder.Build()
}
