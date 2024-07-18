package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
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

func (qb *QueryBuilder) SelectRaw(raw string, value ...interface{}) *QueryBuilder {
	qb.builder.SelectRaw(raw, value)
	return qb
}

func (qb *QueryBuilder) Count(columns ...string) *QueryBuilder {
	qb.builder.Count(columns...)
	return qb
}

func (qb *QueryBuilder) Max(column string) *QueryBuilder {
	qb.builder.Max(column)
	return qb
}

func (qb *QueryBuilder) Min(column string) *QueryBuilder {
	qb.builder.Min(column)
	return qb
}

func (qb *QueryBuilder) Sum(column string) *QueryBuilder {
	qb.builder.Sum(column)
	return qb
}

func (qb *QueryBuilder) Avg(column string) *QueryBuilder {
	qb.builder.Avg(column)
	return qb
}

func (qb *QueryBuilder) Where(column string, condition string, value interface{}) *QueryBuilder {
	switch v := value.(type) {
	case QueryBuilder:
		qb.builder.WhereQuery(column, condition, v.builder)
	case []interface{}:
		qb.builder.Where(column, condition, v...)
	}
	return qb
}

func (qb *QueryBuilder) OrWhere(column string, condition string, value interface{}) *QueryBuilder {
	switch v := value.(type) {
	case QueryBuilder:
		qb.builder.OrWhereQuery(column, condition, v.builder)
	case []interface{}:
		qb.builder.OrWhere(column, condition, v...)
	}

	return qb
}

func (qb *QueryBuilder) WhereQuery(column string, condition string, q *QueryBuilder) *QueryBuilder {
	qb.builder.WhereQuery(column, condition, q.builder)
	return qb
}

func (qb *QueryBuilder) OrWhereQuery(column string, condition string, q *QueryBuilder) *QueryBuilder {
	qb.builder.OrWhereQuery(column, condition, q.builder)
	return qb
}

// WhereGroup
func (qb *QueryBuilder) WhereGroup(fn func(qb *query.Builder) *query.Builder) *QueryBuilder {
	qb.builder.WhereGroup(func(b *query.Builder) *query.Builder {
		return fn(b)
	})
	return qb
}

func (qb *QueryBuilder) OrWhereGroup(fn func(qb *query.Builder) *query.Builder) *QueryBuilder {
	qb.builder.OrWhereGroup(func(b *query.Builder) *query.Builder {
		return fn(b)
	})
	return qb
}

func (qb *QueryBuilder) Join(table, my, condition, target string) *QueryBuilder {
	qb.builder.Join(table, my, condition, target)
	return qb
}

func (qb *QueryBuilder) LeftJoin(table, my, condition, target string) *QueryBuilder {
	qb.builder.LeftJoin(table, my, condition, target)
	return qb
}

func (qb *QueryBuilder) RightJoin(table, my, condition, target string) *QueryBuilder {
	qb.builder.RightJoin(table, my, condition, target)
	return qb
}

func (qb *QueryBuilder) CrossJoin(table, my, condition, target string) *QueryBuilder {
	qb.builder.CrossJoin(table)
	return qb
}

func (qb *QueryBuilder) OrderBy(column, ascDesc string) *QueryBuilder {
	qb.builder.OrderBy(column, ascDesc)
	return qb
}

func (qb *QueryBuilder) OrderByRaw(raw string) *QueryBuilder {
	qb.builder.OrderByRaw(raw)
	return qb
}

func (qb *QueryBuilder) ReOrder() *QueryBuilder {
	qb.builder.ReOrder()
	return qb
}

func (qb *QueryBuilder) GroupBy(columns ...string) *QueryBuilder {
	qb.builder.GroupBy(columns...)
	return qb
}

func (qb *QueryBuilder) Having(column, condition string, value interface{}) *QueryBuilder {
	qb.builder.Having(column, condition, value)
	return qb
}

func (qb *QueryBuilder) HavingRaw(raw string) *QueryBuilder {
	qb.builder.HavingRaw(raw)
	return qb
}

func (qb *QueryBuilder) OrHaving(column, condition string, value interface{}) *QueryBuilder {
	qb.builder.OrHaving(column, condition, value)
	return qb
}

func (qb *QueryBuilder) OrHavingRaw(raw string) *QueryBuilder {
	qb.builder.OrHavingRaw(raw)
	return qb
}

func (qb *QueryBuilder) Limit(limit int64) *QueryBuilder {
	qb.builder.Limit(limit)
	return qb
}

func (qb *QueryBuilder) Take(limit int64) *QueryBuilder {
	qb.builder.Limit(limit)
	return qb
}

func (qb *QueryBuilder) Offset(offset int64) *QueryBuilder {
	qb.builder.Offset(offset)
	return qb
}

func (qb *QueryBuilder) Skip(offset int64) *QueryBuilder {
	qb.builder.Offset(offset)
	return qb
}

func (qb *QueryBuilder) SharedLock() *QueryBuilder {
	qb.builder.SharedLock()
	return qb
}

func (qb *QueryBuilder) LockForUpdate() *QueryBuilder {
	qb.builder.LockForUpdate()
	return qb
}

func (qb *QueryBuilder) Build() (string, []interface{}) {
	return qb.builder.Build()
}

func (qb *QueryBuilder) GetQuery() *structs.Query {
	return qb.builder.GetQuery()
}
