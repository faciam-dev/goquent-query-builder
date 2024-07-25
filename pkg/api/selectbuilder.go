package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type SelectBuilder struct {
	WhereQueryBuilder[SelectBuilder, query.Builder]
	builder             *query.Builder
	joinQueryBuilder    *JoinQueryBuilder
	orderByQueryBuilder *OrderByQueryBuilder
}

func NewSelectBuilder(strategy db.QueryBuilderStrategy, cache cache.Cache) *SelectBuilder {
	sb := &SelectBuilder{
		WhereQueryBuilder: *NewWhereQueryBuilder[SelectBuilder, query.Builder](strategy, cache),
		builder:           query.NewBuilder(strategy, cache),
		joinQueryBuilder: &JoinQueryBuilder{
			builder: query.NewJoinBuilder(strategy, cache),
		},
		orderByQueryBuilder: &OrderByQueryBuilder{
			builder: query.NewOrderByBuilder(&[]structs.Order{}),
		},
	}
	whereBuilder := NewWhereQueryBuilder[SelectBuilder, query.Builder](strategy, cache)
	whereBuilder.SetParent(sb)
	sb.WhereQueryBuilder = *whereBuilder

	return sb
}

func (qb *SelectBuilder) Table(table string) *SelectBuilder {
	qb.builder.Table(table)
	return qb
}

func (qb *SelectBuilder) Select(columns ...string) *SelectBuilder {
	qb.builder.Select(columns...)
	return qb
}

func (qb *SelectBuilder) SelectRaw(raw string, value ...interface{}) *SelectBuilder {
	qb.builder.SelectRaw(raw, value)
	return qb
}

func (qb *SelectBuilder) Count(columns ...string) *SelectBuilder {
	qb.builder.Count(columns...)
	return qb
}

func (qb *SelectBuilder) Max(column string) *SelectBuilder {
	qb.builder.Max(column)
	return qb
}

func (qb *SelectBuilder) Min(column string) *SelectBuilder {
	qb.builder.Min(column)
	return qb
}

func (qb *SelectBuilder) Sum(column string) *SelectBuilder {
	qb.builder.Sum(column)
	return qb
}

func (qb *SelectBuilder) Avg(column string) *SelectBuilder {
	qb.builder.Avg(column)
	return qb
}

// Join
func (qb *SelectBuilder) Join(table, my, condition, target string) *SelectBuilder {
	qb.joinQueryBuilder.Join(table, my, condition, target)

	return qb
}

func (qb *SelectBuilder) LeftJoin(table, my, condition, target string) *SelectBuilder {
	qb.joinQueryBuilder.LeftJoin(table, my, condition, target)

	return qb
}

func (qb *SelectBuilder) RightJoin(table, my, condition, target string) *SelectBuilder {
	qb.joinQueryBuilder.RightJoin(table, my, condition, target)

	return qb
}

func (qb *SelectBuilder) CrossJoin(table, my, condition, target string) *SelectBuilder {
	qb.joinQueryBuilder.CrossJoin(table)

	return qb
}

func (qb *SelectBuilder) JoinQuery(table string, fn func(b *query.JoinClauseBuilder) *query.JoinClauseBuilder) *SelectBuilder {
	qb.joinQueryBuilder.JoinQuery(table, fn)

	return qb
}

func (qb *SelectBuilder) JoinSub(q *SelectBuilder, alias, my, condition, target string) *SelectBuilder {
	qb.joinQueryBuilder.JoinSub(q, alias, my, condition, target)
	return qb
}

func (qb *SelectBuilder) LeftJoinSub(q *SelectBuilder, alias, my, condition, target string) *SelectBuilder {
	qb.joinQueryBuilder.LeftJoinSub(q, alias, my, condition, target)
	return qb
}

func (qb *SelectBuilder) RightJoinSub(q *SelectBuilder, alias, my, condition, target string) *SelectBuilder {
	qb.joinQueryBuilder.RightJoinSub(q, alias, my, condition, target)
	return qb
}

func (qb *SelectBuilder) OrderBy(column, ascDesc string) *SelectBuilder {
	qb.orderByQueryBuilder.OrderBy(column, ascDesc)
	return qb
}

func (qb *SelectBuilder) OrderByRaw(raw string) *SelectBuilder {
	qb.orderByQueryBuilder.OrderByRaw(raw)
	return qb
}

func (qb *SelectBuilder) ReOrder() *SelectBuilder {
	qb.orderByQueryBuilder.ReOrder()
	return qb
}

func (qb *SelectBuilder) GroupBy(columns ...string) *SelectBuilder {
	qb.builder.GroupBy(columns...)
	return qb
}

func (qb *SelectBuilder) Having(column, condition string, value interface{}) *SelectBuilder {
	qb.builder.Having(column, condition, value)
	return qb
}

func (qb *SelectBuilder) HavingRaw(raw string) *SelectBuilder {
	qb.builder.HavingRaw(raw)
	return qb
}

func (qb *SelectBuilder) OrHaving(column, condition string, value interface{}) *SelectBuilder {
	qb.builder.OrHaving(column, condition, value)
	return qb
}

func (qb *SelectBuilder) OrHavingRaw(raw string) *SelectBuilder {
	qb.builder.OrHavingRaw(raw)
	return qb
}

func (qb *SelectBuilder) Limit(limit int64) *SelectBuilder {
	qb.builder.Limit(limit)
	return qb
}

func (qb *SelectBuilder) Take(limit int64) *SelectBuilder {
	qb.builder.Limit(limit)
	return qb
}

func (qb *SelectBuilder) Offset(offset int64) *SelectBuilder {
	qb.builder.Offset(offset)
	return qb
}

func (qb *SelectBuilder) Skip(offset int64) *SelectBuilder {
	qb.builder.Offset(offset)
	return qb
}

func (qb *SelectBuilder) SharedLock() *SelectBuilder {
	qb.builder.SharedLock()
	return qb
}

func (qb *SelectBuilder) LockForUpdate() *SelectBuilder {
	qb.builder.LockForUpdate()
	return qb
}

func (qb *SelectBuilder) Build() (string, []interface{}) {
	//	qb.builder.SetWhereBuilder(qb.whereQueryBuilder.builder)
	qb.builder.SetWhereBuilder(qb.WhereQueryBuilder.builder)
	qb.builder.SetJoinBuilder(qb.joinQueryBuilder.builder)
	qb.builder.SetOrderByBuilder(qb.orderByQueryBuilder.builder)
	return qb.builder.Build()
}

func (qb *SelectBuilder) GetQuery() *structs.Query {
	return qb.builder.GetQuery()
}
