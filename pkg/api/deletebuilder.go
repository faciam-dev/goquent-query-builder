package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type DeleteBuilder struct {
	WhereQueryBuilder[DeleteBuilder, query.DeleteBuilder]
	JoinQueryBuilder[DeleteBuilder, query.DeleteBuilder]
	builder             *query.DeleteBuilder
	orderByQueryBuilder *OrderByQueryBuilder
}

func NewDeleteBuilder(strategy db.QueryBuilderStrategy, cache cache.Cache) *DeleteBuilder {
	db := &DeleteBuilder{
		builder: query.NewDeleteBuilder(strategy, cache),
		orderByQueryBuilder: &OrderByQueryBuilder{
			builder: query.NewOrderByBuilder(&[]structs.Order{}),
		},
	}

	whereBuilder := NewWhereQueryBuilder[DeleteBuilder, query.DeleteBuilder](strategy, cache)
	whereBuilder.SetParent(db)
	db.WhereQueryBuilder = *whereBuilder

	joinBuilder := NewJoinQueryBuilder[DeleteBuilder, query.DeleteBuilder](strategy, cache)
	joinBuilder.SetParent(db)
	db.JoinQueryBuilder = *joinBuilder

	return db
}

func (qb *DeleteBuilder) Delete() *DeleteBuilder {
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

func (qb *DeleteBuilder) Table(table string) *DeleteBuilder {
	qb.builder.Table(table)
	return qb
}

// OrderBy
func (qb *DeleteBuilder) OrderBy(column, ascDesc string) *DeleteBuilder {
	qb.orderByQueryBuilder.OrderBy(column, ascDesc)
	return qb
}

func (qb *DeleteBuilder) OrderByRaw(raw string) *DeleteBuilder {
	qb.orderByQueryBuilder.OrderByRaw(raw)
	return qb
}

func (qb *DeleteBuilder) ReOrder() *DeleteBuilder {
	qb.orderByQueryBuilder.ReOrder()
	return qb
}

func (qb *DeleteBuilder) Build() (string, []interface{}) {
	qb.builder.SetWhereBuilder(qb.WhereQueryBuilder.builder)
	qb.builder.SetJoinBuilder(qb.JoinQueryBuilder.builder)
	qb.builder.SetOrderByBuilder(qb.orderByQueryBuilder.builder)

	return qb.builder.Build()
}
