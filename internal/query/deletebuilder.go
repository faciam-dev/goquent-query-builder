package query

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
)

type DeleteBuilder struct {
	dbBuilder db.QueryBuilderStrategy
	cache     cache.Cache
	query     *structs.DeleteQuery
	WhereBuilder[DeleteBuilder]
	joinBuilder    *JoinBuilder
	orderByBuilder *OrderByBuilder
}

func NewDeleteBuilder(strategy db.QueryBuilderStrategy, cache cache.Cache) *DeleteBuilder {
	db := &DeleteBuilder{
		dbBuilder: strategy,
		cache:     cache,
		query: &structs.DeleteQuery{
			Query: &structs.Query{},
		},
		WhereBuilder:   *NewWhereBuilder[DeleteBuilder](strategy, cache),
		joinBuilder:    NewJoinBuilder(strategy, cache),
		orderByBuilder: NewOrderByBuilder(&[]structs.Order{}),
	}

	whereBuilder := NewWhereBuilder[DeleteBuilder](strategy, cache)
	whereBuilder.SetParent(db)
	db.WhereBuilder = *whereBuilder
	return db
}

func (b *DeleteBuilder) SetWhereBuilder(whereBuilder *WhereBuilder[DeleteBuilder]) {
	b.WhereBuilder = *whereBuilder
}

func (b *DeleteBuilder) SetJoinBuilder(joinBuilder *JoinBuilder) {
	b.joinBuilder = joinBuilder
}

func (b *DeleteBuilder) SetOrderByBuilder(orderByBuilder *OrderByBuilder) {
	b.orderByBuilder = orderByBuilder
}

func (b *DeleteBuilder) Table(table string) *DeleteBuilder {
	b.query.Table = table
	b.joinBuilder.Table.Name = table
	return b
}

// Delete
func (b *DeleteBuilder) Delete() *DeleteBuilder {
	return b
}

func (u *DeleteBuilder) Build() (string, []interface{}) {
	// If there are conditions, add them to the query
	if len(*u.WhereBuilder.query.Conditions) > 0 {
		*u.WhereBuilder.query.ConditionGroups = append(*u.WhereBuilder.query.ConditionGroups, structs.WhereGroup{
			Conditions:   *u.WhereBuilder.query.Conditions,
			Operator:     consts.LogicalOperator_AND,
			IsDummyGroup: true,
		})
		u.WhereBuilder.query.Conditions = &[]structs.Where{}
	}

	u.query.Query.Conditions = u.WhereBuilder.query.Conditions
	u.query.Query.ConditionGroups = u.WhereBuilder.query.ConditionGroups
	u.query.Query.Joins = u.joinBuilder.Joins
	u.query.Query.Order = u.orderByBuilder.Order

	query, values := u.dbBuilder.BuildDelete(u.query)
	return query, values
}

func (b *DeleteBuilder) Join(table, my, condition, target string) *DeleteBuilder {
	b.joinBuilder.Join(table, my, condition, target)
	return b
}

func (b *DeleteBuilder) LeftJoin(table, my, condition, target string) *DeleteBuilder {
	b.joinBuilder.LeftJoin(table, my, condition, target)
	return b
}

func (b *DeleteBuilder) RightJoin(table, my, condition, target string) *DeleteBuilder {
	b.joinBuilder.RightJoin(table, my, condition, target)
	return b
}

func (b *DeleteBuilder) CrossJoin(table string) *DeleteBuilder {
	b.joinBuilder.CrossJoin(table)
	return b
}

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

func (b *DeleteBuilder) GetQuery() *structs.DeleteQuery {
	return b.query
}
