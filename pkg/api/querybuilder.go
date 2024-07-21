package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type QueryBuilder struct {
	builder             *query.Builder
	whereQueryBuilder   *WhereQueryBuilder
	joinQueryBuilder    *JoinQueryBuilder
	orderByQueryBuilder *OrderByQueryBuilder
}

func NewQueryBuilder(strategy db.QueryBuilderStrategy, cache cache.Cache) *QueryBuilder {
	return &QueryBuilder{
		builder: query.NewBuilder(strategy, cache),
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

// Where
func (qb *QueryBuilder) Where(column string, condition string, value interface{}) *QueryBuilder {
	qb.whereQueryBuilder.Where(column, condition, value)

	return qb
}

// OrWhere
func (qb *QueryBuilder) OrWhere(column string, condition string, value interface{}) *QueryBuilder {
	qb.whereQueryBuilder.OrWhere(column, condition, value)

	return qb
}

func (qb *QueryBuilder) WhereRaw(raw string, value ...interface{}) *QueryBuilder {
	qb.whereQueryBuilder.WhereRaw(raw, value)

	return qb
}

func (qb *QueryBuilder) OrWhereRaw(raw string, value ...interface{}) *QueryBuilder {
	qb.whereQueryBuilder.OrWhereRaw(raw, value)

	return qb
}

// WhereQuery
func (qb *QueryBuilder) WhereQuery(column string, condition string, q *QueryBuilder) *QueryBuilder {
	qb.whereQueryBuilder.WhereQuery(column, condition, q)

	return qb
}

// OrWhereQuery
func (qb *QueryBuilder) OrWhereQuery(column string, condition string, q *QueryBuilder) *QueryBuilder {
	qb.whereQueryBuilder.OrWhereQuery(column, condition, q)

	return qb
}

// WhereGroup
func (qb *QueryBuilder) WhereGroup(fn func(wb *query.WhereBuilder) *query.WhereBuilder) *QueryBuilder {
	qb.whereQueryBuilder.WhereGroup(fn)

	return qb
}

// OrWhereGroup
func (qb *QueryBuilder) OrWhereGroup(fn func(qb *query.WhereBuilder) *query.WhereBuilder) *QueryBuilder {
	qb.whereQueryBuilder.OrWhereGroup(fn)

	return qb
}

func (qb *QueryBuilder) Join(table, my, condition, target string) *QueryBuilder {
	qb.joinQueryBuilder.Join(table, my, condition, target)

	return qb
}

func (qb *QueryBuilder) LeftJoin(table, my, condition, target string) *QueryBuilder {
	qb.joinQueryBuilder.LeftJoin(table, my, condition, target)

	return qb
}

func (qb *QueryBuilder) RightJoin(table, my, condition, target string) *QueryBuilder {
	qb.joinQueryBuilder.RightJoin(table, my, condition, target)

	return qb
}

func (qb *QueryBuilder) CrossJoin(table, my, condition, target string) *QueryBuilder {
	qb.joinQueryBuilder.CrossJoin(table)

	return qb
}

func (qb *QueryBuilder) JoinQuery(table string, fn func(b *query.JoinClauseBuilder) *query.JoinClauseBuilder) *QueryBuilder {
	qb.joinQueryBuilder.JoinQuery(table, fn)

	return qb
}

func (qb *QueryBuilder) JoinSub(q *QueryBuilder, alias, my, condition, target string) *QueryBuilder {
	qb.joinQueryBuilder.JoinSub(q, alias, my, condition, target)
	return qb
}

func (qb *QueryBuilder) LeftJoinSub(q *QueryBuilder, alias, my, condition, target string) *QueryBuilder {
	qb.joinQueryBuilder.LeftJoinSub(q, alias, my, condition, target)
	return qb
}

func (qb *QueryBuilder) RightJoinSub(q *QueryBuilder, alias, my, condition, target string) *QueryBuilder {
	qb.joinQueryBuilder.RightJoinSub(q, alias, my, condition, target)
	return qb
}

func (qb *QueryBuilder) OrderBy(column, ascDesc string) *QueryBuilder {
	qb.orderByQueryBuilder.OrderBy(column, ascDesc)
	return qb
}

func (qb *QueryBuilder) OrderByRaw(raw string) *QueryBuilder {
	qb.orderByQueryBuilder.OrderByRaw(raw)
	return qb
}

func (qb *QueryBuilder) ReOrder() *QueryBuilder {
	qb.orderByQueryBuilder.ReOrder()
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
	qb.builder.SetWhereBuilder(qb.whereQueryBuilder.builder)
	qb.builder.SetJoinBuilder(qb.joinQueryBuilder.builder)
	qb.builder.SetOrderByBuilder(qb.orderByQueryBuilder.builder)
	return qb.builder.Build()
}

func (qb *QueryBuilder) GetQuery() *structs.Query {
	return qb.builder.GetQuery()
}
