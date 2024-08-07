package api

import (
	"github.com/faciam-dev/goquent-query-builder/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type DeleteQueryBuilder struct {
	WhereQueryBuilder[*DeleteQueryBuilder, query.DeleteBuilder]
	JoinQueryBuilder[*DeleteQueryBuilder, query.DeleteBuilder]
	OrderByQueryBuilder[*DeleteQueryBuilder, query.DeleteBuilder]
	builder *query.DeleteBuilder
	QueryBuilderStrategy[DeleteQueryBuilder, query.DeleteBuilder]
}

func NewDeleteQueryBuilder(strategy interfaces.QueryBuilderStrategy, cache cache.Cache) *DeleteQueryBuilder {
	db := &DeleteQueryBuilder{
		builder: query.NewDeleteBuilder(strategy, cache),
	}

	whereBuilder := NewWhereQueryBuilder[*DeleteQueryBuilder, query.DeleteBuilder](strategy, cache)
	whereBuilder.SetParent(&db)
	db.WhereQueryBuilder = *whereBuilder

	joinBuilder := NewJoinQueryBuilder[*DeleteQueryBuilder, query.DeleteBuilder](strategy, cache)
	joinBuilder.SetParent(&db)
	db.JoinQueryBuilder = *joinBuilder

	orderByBuilder := NewOrderByQueryBuilder[*DeleteQueryBuilder, query.DeleteBuilder](strategy, cache)
	orderByBuilder.SetParent(&db)
	db.OrderByQueryBuilder = *orderByBuilder

	return db
}

func (qb *DeleteQueryBuilder) Delete() *DeleteQueryBuilder {
	qb.builder.Delete()

	return qb
}

// Using
/*
func (qb *UpdateQueryBuilder) Using(qb *QueryBuilder) *DeleteQueryBuilder {
	qb.builder.Using(qb)

	return ub
}
*/

func (qb *DeleteQueryBuilder) Table(table string) *DeleteQueryBuilder {
	qb.builder.Table(table)
	return qb
}

func (ub *DeleteQueryBuilder) Dump() (string, []interface{}, error) {
	b := query.NewDebugBuilder[*query.DeleteBuilder, DeleteQueryBuilder](ub.builder)

	return b.Dump()
}

func (ub *DeleteQueryBuilder) RawSql() (string, error) {
	b := query.NewDebugBuilder[*query.DeleteBuilder, DeleteQueryBuilder](ub.builder)

	return b.RawSql()
}

func (qb *DeleteQueryBuilder) Build() (string, []interface{}, error) {
	return qb.builder.Build()
}

func (qb *DeleteQueryBuilder) GetQueryBuilder() *DeleteQueryBuilder {
	return qb
}

func (qb *DeleteQueryBuilder) GetWhereBuilder() *query.WhereBuilder[query.DeleteBuilder] {
	return &qb.builder.WhereBuilder
}

func (qb *DeleteQueryBuilder) GetJoinBuilder() *query.JoinBuilder[query.DeleteBuilder] {
	return &qb.builder.JoinBuilder
}

func (qb *DeleteQueryBuilder) GetOrderByBuilder() *query.OrderByBuilder[query.DeleteBuilder] {
	return &qb.builder.OrderByBuilder
}
