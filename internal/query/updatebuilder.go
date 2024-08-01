package query

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
)

type UpdateBuilder struct {
	dbBuilder      db.QueryBuilderStrategy
	cache          cache.Cache
	query          *structs.UpdateQuery
	orderByBuilder *OrderByBuilder
	JoinBuilder[UpdateBuilder]
	WhereBuilder[UpdateBuilder]
}

func NewUpdateBuilder(strategy db.QueryBuilderStrategy, cache cache.Cache) *UpdateBuilder {
	ub := &UpdateBuilder{
		dbBuilder: strategy,
		cache:     cache,
		query: &structs.UpdateQuery{
			Query: &structs.Query{},
		},
		orderByBuilder: NewOrderByBuilder(&[]structs.Order{}),
	}

	whereBuilder := NewWhereBuilder[UpdateBuilder](strategy, cache)
	whereBuilder.SetParent(ub)
	ub.WhereBuilder = *whereBuilder

	joinBuilder := NewJoinBuilder[UpdateBuilder](strategy, cache)
	joinBuilder.SetParent(ub)
	ub.JoinBuilder = *joinBuilder

	return ub
}

func (b *UpdateBuilder) SetWhereBuilder(whereBuilder WhereBuilder[UpdateBuilder]) {
	b.WhereBuilder = whereBuilder
}

func (b *UpdateBuilder) SetJoinBuilder(joinBuilder JoinBuilder[UpdateBuilder]) {
	b.JoinBuilder = joinBuilder
}

func (b *UpdateBuilder) SetOrderByBuilder(orderByBuilder *OrderByBuilder) {
	b.orderByBuilder = orderByBuilder
}

func (b *UpdateBuilder) Table(table string) *UpdateBuilder {
	b.query.Table = table
	b.JoinBuilder.Table.Name = table
	return b
}

func (b *UpdateBuilder) Update(data map[string]interface{}) *UpdateBuilder {
	b.query.Values = data

	return b
}

func (u *UpdateBuilder) Build() (string, []interface{}, error) {
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
	u.query.Query.Joins = u.JoinBuilder.Joins
	u.query.Query.Order = u.orderByBuilder.Order

	query, values, err := u.dbBuilder.BuildUpdate(u.query)
	return query, values, err
}

func (b *UpdateBuilder) OrderBy(column string, direction string) *UpdateBuilder {
	b.orderByBuilder.OrderBy(column, direction)
	return b
}

func (b *UpdateBuilder) OrderByRaw(raw string) *UpdateBuilder {
	b.orderByBuilder.OrderByRaw(raw)
	return b
}

func (b *UpdateBuilder) ReOrder() *UpdateBuilder {
	b.orderByBuilder.ReOrder()
	return b
}

func (b *UpdateBuilder) GetQuery() *structs.UpdateQuery {
	return b.query
}
