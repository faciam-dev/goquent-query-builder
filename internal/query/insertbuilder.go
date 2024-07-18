package query

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
)

type InsertBuilder struct {
	dbBuilder db.QueryBuilderStrategy
	cache     *cache.AsyncQueryCache
	query     *structs.InsertQuery
}

func NewInsertBuilder(dbBuilder db.QueryBuilderStrategy, cache *cache.AsyncQueryCache) *InsertBuilder {
	return &InsertBuilder{
		dbBuilder: dbBuilder,
		cache:     cache,
		query:     &structs.InsertQuery{},
	}
}

func (ib *InsertBuilder) Table(table string) *InsertBuilder {
	ib.query.Table = table
	return ib
}

func (ib *InsertBuilder) Insert(data map[string]interface{}) *InsertBuilder {
	ib.query.Values = data
	return ib
}

func (ib *InsertBuilder) InsertBatch(data []map[string]interface{}) *InsertBuilder {
	ib.query.ValuesBatch = data
	return ib
}

func (ib *InsertBuilder) InsertUsing(columns []string, q *structs.Query) *InsertBuilder {
	ib.query.Columns = columns

	// If there are conditions, add them to the query
	if len(*q.Conditions) > 0 {
		*q.ConditionGroups = append(*q.ConditionGroups, structs.WhereGroup{
			Conditions:   *q.Conditions,
			Operator:     consts.LogicalOperator_AND,
			IsDummyGroup: true,
		})
		q.Conditions = &[]structs.Where{}
	}

	ib.query.SelectQuery = q

	return ib
}

func (ib *InsertBuilder) Build() (string, []interface{}) {
	query, values := ib.dbBuilder.BuildInsert(ib.query)
	return query, values
}
