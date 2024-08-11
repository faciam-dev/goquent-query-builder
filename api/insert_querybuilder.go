package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type InsertQueryBuilder struct {
	builder *query.InsertBuilder
}

func NewInsertQueryBuilder(strategy interfaces.QueryBuilderStrategy) *InsertQueryBuilder {
	return &InsertQueryBuilder{
		builder: query.NewInsertBuilder(strategy),
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

func (ib *InsertQueryBuilder) InsertUsing(columns []string, qb *SelectQueryBuilder) *InsertQueryBuilder {
	ib.builder.InsertUsing(columns, qb.builder)
	return ib
}

func (ib *InsertQueryBuilder) Dump() (string, []interface{}, error) {
	b := query.NewDebugBuilder[*query.InsertBuilder, InsertQueryBuilder](ib.builder)

	return b.Dump()
}

func (ib *InsertQueryBuilder) RawSql() (string, error) {
	b := query.NewDebugBuilder[*query.InsertBuilder, InsertQueryBuilder](ib.builder)

	return b.RawSql()
}

func (ib *InsertQueryBuilder) Build() (string, []interface{}, error) {
	return ib.builder.Build()
}
