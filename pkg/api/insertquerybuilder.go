package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type InsertBuilder struct {
	builder *query.InsertBuilder
}

func NewInsertBuilder(strategy db.QueryBuilderStrategy, cache cache.Cache) *InsertBuilder {
	return &InsertBuilder{
		builder: query.NewInsertBuilder(strategy, cache),
	}
}

func (ib *InsertBuilder) Table(table string) *InsertBuilder {
	ib.builder.Table(table)
	return ib
}

func (ib *InsertBuilder) Insert(data map[string]interface{}) *InsertBuilder {
	ib.builder.Insert(data)
	return ib
}

func (ib *InsertBuilder) InsertBatch(data []map[string]interface{}) *InsertBuilder {
	ib.builder.InsertBatch(data)
	return ib
}

func (ib *InsertBuilder) InsertUsing(columns []string, qb *QueryBuilder) *InsertBuilder {
	ib.builder.InsertUsing(columns, qb.builder)
	return ib
}

func (ib *InsertBuilder) Build() (string, []interface{}) {
	return ib.builder.Build()
}
