package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type DeleteQueryBuilder struct {
	WhereQueryBuilder[DeleteQueryBuilder, query.DeleteBuilder]
	JoinQueryBuilder[DeleteQueryBuilder, query.DeleteBuilder]
	builder             *query.DeleteBuilder
	orderByQueryBuilder *OrderByQueryBuilder
}

func NewDeleteBuilder(strategy db.QueryBuilderStrategy, cache cache.Cache) *DeleteQueryBuilder {
	db := &DeleteQueryBuilder{
		builder: query.NewDeleteBuilder(strategy, cache),
		orderByQueryBuilder: &OrderByQueryBuilder{
			builder: query.NewOrderByBuilder(&[]structs.Order{}),
		},
	}

	whereBuilder := NewWhereQueryBuilder[DeleteQueryBuilder, query.DeleteBuilder](strategy, cache)
	whereBuilder.SetParent(db)
	db.WhereQueryBuilder = *whereBuilder

	joinBuilder := NewJoinQueryBuilder[DeleteQueryBuilder, query.DeleteBuilder](strategy, cache)
	joinBuilder.SetParent(db)
	db.JoinQueryBuilder = *joinBuilder

	return db
}

func (qb *DeleteQueryBuilder) Delete() *DeleteQueryBuilder {
	qb.builder.Delete()

	return qb
}

// Using
/*
func (ub *UpdateQueryBuilder) Using(qb *QueryBuilder) *UpdateQueryBuilder {
	ub.builder.Using(qb)

	return ub
}
*/

func (qb *DeleteQueryBuilder) Table(table string) *DeleteQueryBuilder {
	qb.builder.Table(table)
	return qb
}

// OrderBy
func (qb *DeleteQueryBuilder) OrderBy(column, ascDesc string) *DeleteQueryBuilder {
	qb.orderByQueryBuilder.OrderBy(column, ascDesc)
	return qb
}

func (qb *DeleteQueryBuilder) OrderByRaw(raw string) *DeleteQueryBuilder {
	qb.orderByQueryBuilder.OrderByRaw(raw)
	return qb
}

func (qb *DeleteQueryBuilder) ReOrder() *DeleteQueryBuilder {
	qb.orderByQueryBuilder.ReOrder()
	return qb
}

func (ub *DeleteQueryBuilder) Dump() (string, []interface{}, error) {
	ub.builder.SetOrderByBuilder(ub.orderByQueryBuilder.builder)
	ub.builder.SetWhereBuilder(ub.WhereQueryBuilder.builder)
	ub.builder.SetJoinBuilder(ub.JoinQueryBuilder.builder)

	b := query.NewDebugBuilder[*query.DeleteBuilder, DeleteQueryBuilder](ub.builder)

	return b.Dump()
}

func (ub *DeleteQueryBuilder) RawSql() (string, error) {
	ub.builder.SetOrderByBuilder(ub.orderByQueryBuilder.builder)
	ub.builder.SetWhereBuilder(ub.WhereQueryBuilder.builder)
	ub.builder.SetJoinBuilder(ub.JoinQueryBuilder.builder)

	b := query.NewDebugBuilder[*query.DeleteBuilder, DeleteQueryBuilder](ub.builder)

	return b.RawSql()
}

func (qb *DeleteQueryBuilder) Build() (string, []interface{}, error) {
	qb.builder.SetWhereBuilder(qb.WhereQueryBuilder.builder)
	qb.builder.SetJoinBuilder(qb.JoinQueryBuilder.builder)
	qb.builder.SetOrderByBuilder(qb.orderByQueryBuilder.builder)

	return qb.builder.Build()
}
