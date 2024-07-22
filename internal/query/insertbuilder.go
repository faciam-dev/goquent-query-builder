package query

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
)

type InsertBuilder struct {
	dbBuilder db.QueryBuilderStrategy
	cache     cache.Cache
	query     *structs.InsertQuery
}

func NewInsertBuilder(dbBuilder db.QueryBuilderStrategy, cache cache.Cache) *InsertBuilder {
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

func (ib *InsertBuilder) InsertUsing(columns []string, b *Builder) *InsertBuilder {
	ib.query.Columns = columns

	// If there are conditions, add them to the query
	if b.whereBuilder.query.Conditions != nil && len(*b.whereBuilder.query.Conditions) > 0 {
		*b.whereBuilder.query.ConditionGroups = append(*b.whereBuilder.query.ConditionGroups, structs.WhereGroup{
			Conditions:   *b.whereBuilder.query.Conditions,
			Operator:     consts.LogicalOperator_AND,
			IsDummyGroup: true,
		})
		b.whereBuilder.query.Conditions = &[]structs.Where{}
	}

	b.buildQuery()
	ib.query.Query = b.GetQuery()

	return ib
}

func (ib *InsertBuilder) Build() (string, []interface{}) {
	query, values := ib.dbBuilder.BuildInsert(ib.query)
	return query, values
}
