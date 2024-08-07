package query

import (
	"github.com/faciam-dev/goquent-query-builder/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

type DeleteBuilder struct {
	dbBuilder interfaces.QueryBuilderStrategy
	cache     cache.Cache
	query     *structs.DeleteQuery
	WhereBuilder[DeleteBuilder]
	JoinBuilder[DeleteBuilder]
	OrderByBuilder[DeleteBuilder]
}

func NewDeleteBuilder(strategy interfaces.QueryBuilderStrategy, cache cache.Cache) *DeleteBuilder {
	db := &DeleteBuilder{
		dbBuilder: strategy,
		cache:     cache,
		query: &structs.DeleteQuery{
			Query: &structs.Query{},
		},
	}

	whereBuilder := NewWhereBuilder[DeleteBuilder](strategy, cache)
	whereBuilder.SetParent(db)
	db.WhereBuilder = *whereBuilder

	joinBuilder := NewJoinBuilder[DeleteBuilder](strategy, cache)
	joinBuilder.SetParent(db)
	db.JoinBuilder = *joinBuilder

	orderByBuilder := NewOrderByBuilder[DeleteBuilder](strategy, cache)
	orderByBuilder.SetParent(db)
	db.OrderByBuilder = *orderByBuilder

	return db
}

/*
	func (b *DeleteBuilder) SetWhereBuilder(whereBuilder *WhereBuilder[DeleteBuilder]) {
		b.WhereBuilder = *whereBuilder
	}

	func (b *DeleteBuilder) SetJoinBuilder(joinBuilder *JoinBuilder[DeleteBuilder]) {
		b.JoinBuilder = *joinBuilder
	}

	func (b *DeleteBuilder) SetOrderByBuilder(orderByBuilder *OrderByBuilder[DeleteBuilder]) {
		b.OrderByBuilder = *orderByBuilder
	}
*/
func (b *DeleteBuilder) Table(table string) *DeleteBuilder {
	b.query.Table = table
	b.JoinBuilder.Table.Name = table
	return b
}

// Delete
func (b *DeleteBuilder) Delete() *DeleteBuilder {
	return b
}

func (d *DeleteBuilder) Build() (string, []interface{}, error) {
	// If there are conditions, add them to the query
	if len(*d.WhereBuilder.query.Conditions) > 0 {
		d.WhereBuilder.query.ConditionGroups = append(d.WhereBuilder.query.ConditionGroups, structs.WhereGroup{
			Conditions:   *d.WhereBuilder.query.Conditions,
			Operator:     consts.LogicalOperator_AND,
			IsDummyGroup: true,
		})
		d.WhereBuilder.query.Conditions = &[]structs.Where{}
	}

	d.query.Query.Conditions = d.WhereBuilder.query.Conditions
	d.query.Query.ConditionGroups = d.WhereBuilder.query.ConditionGroups
	d.query.Query.Joins = d.JoinBuilder.Joins
	d.query.Query.Order = d.OrderByBuilder.Order

	query, values, err := d.dbBuilder.BuildDelete(d.query)
	return query, values, err
}

/*
func (b *DeleteBuilder) OrderBy(column string, direction string) *DeleteBuilder {
	b.orderByBuilder.OrderBy(column, direction)
	return b
}

func (b *DeleteBuilder) OrderByRaw(raw string) *DeleteBuilder {
	b.orderByBuilder.OrderByRaw(raw)
	return b
}

func (b *DeleteBuilder) ReOrder() *DeleteBuilder {
	b.orderByBuilder.ReOrder()
	return b
}
*/

func (b *DeleteBuilder) GetQuery() *structs.DeleteQuery {
	return b.query
}
