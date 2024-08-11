package query

import (
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

type InsertBuilder struct {
	BaseBuilder
	dbBuilder interfaces.QueryBuilderStrategy
	query     *structs.InsertQuery
}

func NewInsertBuilder(dbBuilder interfaces.QueryBuilderStrategy) *InsertBuilder {
	return &InsertBuilder{
		dbBuilder: dbBuilder,
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

func (ib *InsertBuilder) InsertUsing(columns []string, b *SelectBuilder) *InsertBuilder {
	ib.query.Columns = columns

	// If there are conditions, add them to the query
	if b.WhereBuilder.query.Conditions != nil && len(*b.WhereBuilder.query.Conditions) > 0 {
		b.WhereBuilder.query.ConditionGroups = append(b.WhereBuilder.query.ConditionGroups, structs.WhereGroup{
			Conditions:   *b.WhereBuilder.query.Conditions,
			Operator:     consts.LogicalOperator_AND,
			IsDummyGroup: true,
		})
		b.WhereBuilder.query.Conditions = &[]structs.Where{}
	}

	b.buildQuery()
	ib.query.Query = b.GetQuery()

	return ib
}

func (ib *InsertBuilder) Build() (string, []interface{}, error) {
	query, values, err := ib.dbBuilder.BuildInsert(ib.query)
	return query, values, err
}
