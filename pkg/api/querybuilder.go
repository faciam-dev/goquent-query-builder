package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type QueryBuilder struct {
	builder   *query.Builder
	dbBuilder db.QueryBuilder
}

func NewQueryBuilder(strategy db.QueryBuilderStrategy) *QueryBuilder {
	return &QueryBuilder{
		builder:   query.NewBuilder(),
		dbBuilder: db.NewQueryBuilder(strategy),
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

func (qb *QueryBuilder) Where(condition string, value interface{}) *QueryBuilder {
	qb.builder.Where(condition, value)
	return qb
}

func (qb *QueryBuilder) Join(joinType, table, condition string) *QueryBuilder {
	qb.builder.Join(joinType, table, condition)
	return qb
}

func (qb *QueryBuilder) OrderBy(orderBy ...string) *QueryBuilder {
	qb.builder.OrderBy(orderBy...)
	return qb
}

func (qb *QueryBuilder) Build() (string, []interface{}) {
	query := qb.dbBuilder.Select(qb.builder.Columns()...)
	query += " " + qb.dbBuilder.From(qb.builder.Table())
	if len(qb.builder.Joins()) > 0 {
		for _, join := range qb.builder.Joins() {
			query += " " + qb.dbBuilder.Join(join.Type, join.Table, join.Condition)
		}
	}
	if len(qb.builder.Conditions()) > 0 {
		query += " " + qb.dbBuilder.Where(qb.builder.Conditions()...)
	}
	if len(qb.builder.OrderBy()) > 0 {
		query += " " + qb.dbBuilder.OrderBy(qb.builder.OrderBy()...)
	}
	return query, qb.builder.Values()
}
