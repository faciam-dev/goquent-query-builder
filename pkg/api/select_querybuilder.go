package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type SelectQueryBuilder struct {
	WhereQueryBuilder[SelectQueryBuilder, query.Builder]
	JoinQueryBuilder[SelectQueryBuilder, query.Builder]
	builder             *query.Builder
	orderByQueryBuilder *OrderByQueryBuilder
	Queries             []*structs.Query
}

func NewSelectQueryBuilder(strategy db.QueryBuilderStrategy, cache cache.Cache) *SelectQueryBuilder {
	sb := &SelectQueryBuilder{
		//WhereQueryBuilder: *NewWhereQueryBuilder[SelectBuilder, query.Builder](strategy, cache),
		builder: query.NewBuilder(strategy, cache),
		orderByQueryBuilder: &OrderByQueryBuilder{
			builder: query.NewOrderByBuilder(&[]structs.Order{}),
		},
	}
	whereBuilder := NewWhereQueryBuilder[SelectQueryBuilder, query.Builder](strategy, cache)
	whereBuilder.SetParent(sb)
	sb.WhereQueryBuilder = *whereBuilder

	joinBuilder := NewJoinQueryBuilder[SelectQueryBuilder, query.Builder](strategy, cache)
	joinBuilder.SetParent(sb)
	sb.JoinQueryBuilder = *joinBuilder

	return sb
}

func (qb *SelectQueryBuilder) Table(table string) *SelectQueryBuilder {
	qb.builder.Table(table)
	return qb
}

func (qb *SelectQueryBuilder) Select(columns ...string) *SelectQueryBuilder {
	qb.builder.Select(columns...)
	return qb
}

func (qb *SelectQueryBuilder) SelectRaw(raw string, value ...interface{}) *SelectQueryBuilder {
	qb.builder.SelectRaw(raw, value...)
	return qb
}

func (qb *SelectQueryBuilder) Count(columns ...string) *SelectQueryBuilder {
	qb.builder.Count(columns...)
	return qb
}

func (qb *SelectQueryBuilder) Max(column string) *SelectQueryBuilder {
	qb.builder.Max(column)
	return qb
}

func (qb *SelectQueryBuilder) Min(column string) *SelectQueryBuilder {
	qb.builder.Min(column)
	return qb
}

func (qb *SelectQueryBuilder) Sum(column string) *SelectQueryBuilder {
	qb.builder.Sum(column)
	return qb
}

func (qb *SelectQueryBuilder) Avg(column string) *SelectQueryBuilder {
	qb.builder.Avg(column)
	return qb
}

func (qb *SelectQueryBuilder) Distinct(column ...string) *SelectQueryBuilder {
	qb.builder.Distinct(column...)
	return qb
}

func (qb *SelectQueryBuilder) Union(sb *SelectQueryBuilder) *SelectQueryBuilder {
	sb.builder.SetWhereBuilder(sb.WhereQueryBuilder.builder)
	sb.builder.SetJoinBuilder(sb.JoinQueryBuilder.builder)
	sb.builder.SetOrderByBuilder(sb.orderByQueryBuilder.builder)
	qb.Queries = append(qb.Queries, sb.GetQuery())
	qb.builder.Union(sb.builder)
	return qb
}

func (qb *SelectQueryBuilder) UnionAll(sb *SelectQueryBuilder) *SelectQueryBuilder {
	sb.builder.SetWhereBuilder(sb.WhereQueryBuilder.builder)
	sb.builder.SetJoinBuilder(sb.JoinQueryBuilder.builder)
	sb.builder.SetOrderByBuilder(sb.orderByQueryBuilder.builder)
	qb.Queries = append(qb.Queries, sb.GetQuery())
	qb.builder.UnionAll(sb.builder)
	return qb
}

func (qb *SelectQueryBuilder) OrderBy(column, ascDesc string) *SelectQueryBuilder {
	qb.orderByQueryBuilder.OrderBy(column, ascDesc)
	return qb
}

func (qb *SelectQueryBuilder) OrderByRaw(raw string) *SelectQueryBuilder {
	qb.orderByQueryBuilder.OrderByRaw(raw)
	return qb
}

func (qb *SelectQueryBuilder) ReOrder() *SelectQueryBuilder {
	qb.orderByQueryBuilder.ReOrder()
	return qb
}

func (qb *SelectQueryBuilder) GroupBy(columns ...string) *SelectQueryBuilder {
	qb.builder.GroupBy(columns...)
	return qb
}

func (qb *SelectQueryBuilder) Having(column, condition string, value interface{}) *SelectQueryBuilder {
	qb.builder.Having(column, condition, value)
	return qb
}

func (qb *SelectQueryBuilder) HavingRaw(raw string) *SelectQueryBuilder {
	qb.builder.HavingRaw(raw)
	return qb
}

func (qb *SelectQueryBuilder) OrHaving(column, condition string, value interface{}) *SelectQueryBuilder {
	qb.builder.OrHaving(column, condition, value)
	return qb
}

func (qb *SelectQueryBuilder) OrHavingRaw(raw string) *SelectQueryBuilder {
	qb.builder.OrHavingRaw(raw)
	return qb
}

func (qb *SelectQueryBuilder) Limit(limit int64) *SelectQueryBuilder {
	qb.builder.Limit(limit)
	return qb
}

func (qb *SelectQueryBuilder) Take(limit int64) *SelectQueryBuilder {
	qb.builder.Limit(limit)
	return qb
}

func (qb *SelectQueryBuilder) Offset(offset int64) *SelectQueryBuilder {
	qb.builder.Offset(offset)
	return qb
}

func (qb *SelectQueryBuilder) Skip(offset int64) *SelectQueryBuilder {
	qb.builder.Offset(offset)
	return qb
}

func (qb *SelectQueryBuilder) SharedLock() *SelectQueryBuilder {
	qb.builder.SharedLock()
	return qb
}

func (qb *SelectQueryBuilder) LockForUpdate() *SelectQueryBuilder {
	qb.builder.LockForUpdate()
	return qb
}

func (qb *SelectQueryBuilder) Build() (string, []interface{}, error) {
	qb.builder.SetWhereBuilder(qb.WhereQueryBuilder.builder)
	qb.builder.SetJoinBuilder(qb.JoinQueryBuilder.builder)
	qb.builder.SetOrderByBuilder(qb.orderByQueryBuilder.builder)
	return qb.builder.Build()
}

func (qb *SelectQueryBuilder) GetQuery() *structs.Query {
	return qb.builder.GetQuery()
}

func (qb *SelectQueryBuilder) Dump() (string, []interface{}, error) {
	qb.builder.SetWhereBuilder(qb.WhereQueryBuilder.builder)
	qb.builder.SetJoinBuilder(qb.JoinQueryBuilder.builder)
	qb.builder.SetOrderByBuilder(qb.orderByQueryBuilder.builder)
	b := query.NewDebugBuilder[*query.Builder, SelectQueryBuilder](qb.builder)

	return b.Dump()
}

func (qb *SelectQueryBuilder) RawSql() (string, error) {
	qb.builder.SetWhereBuilder(qb.WhereQueryBuilder.builder)
	qb.builder.SetJoinBuilder(qb.JoinQueryBuilder.builder)
	qb.builder.SetOrderByBuilder(qb.orderByQueryBuilder.builder)
	b := query.NewDebugBuilder[*query.Builder, SelectQueryBuilder](qb.builder)

	return b.RawSql()
}
