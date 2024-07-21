package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type InsertQueryBuilder struct {
	builder *query.InsertBuilder
}

func NewInsertQueryBuilder(strategy db.QueryBuilderStrategy, cache cache.Cache) *InsertQueryBuilder {
	return &InsertQueryBuilder{
		builder: query.NewInsertBuilder(strategy, cache),
	}
}

func (ib *InsertQueryBuilder) Table(table string) *InsertQueryBuilder {
	ib.builder.Table(table)
	return ib
}

func (ib *InsertQueryBuilder) Insert(data map[string]interface{}) *InsertQueryBuilder {
	ib.builder.Insert(data)
	return ib
}

func (ib *InsertQueryBuilder) InsertBatch(data []map[string]interface{}) *InsertQueryBuilder {
	ib.builder.InsertBatch(data)
	return ib
}

func (ib *InsertQueryBuilder) InsertUsing(columns []string, qb *QueryBuilder) *InsertQueryBuilder {
	ib.builder.InsertUsing(columns, qb.builder.GetQuery())
	return ib
}

func (ib *InsertQueryBuilder) Build() (string, []interface{}) {
	return ib.builder.Build()
}
