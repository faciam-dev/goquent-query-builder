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
	whereBuilder   *WhereBuilder
	joinBuilder    *JoinBuilder
	orderByBuilder *OrderByBuilder
}

func NewUpdateBuilder(strategy db.QueryBuilderStrategy, cache cache.Cache) *UpdateBuilder {
	return &UpdateBuilder{
		dbBuilder: strategy,
		cache:     cache,
		query: &structs.UpdateQuery{
			Query: &structs.Query{},
		},
		whereBuilder:   NewWhereBuilder(strategy, cache),
		joinBuilder:    NewJoinBuilder(strategy, cache),
		orderByBuilder: NewOrderByBuilder(&[]structs.Order{}),
	}
}

func (b *UpdateBuilder) SetWhereBuilder(whereBuilder *WhereBuilder) {
	b.whereBuilder = whereBuilder
}

func (b *UpdateBuilder) SetJoinBuilder(joinBuilder *JoinBuilder) {
	b.joinBuilder = joinBuilder
}

func (b *UpdateBuilder) SetOrderByBuilder(orderByBuilder *OrderByBuilder) {
	b.orderByBuilder = orderByBuilder
}

func (b *UpdateBuilder) Table(table string) *UpdateBuilder {
	b.query.Table = table
	b.joinBuilder.Table.Name = table
	return b
}

func (b *UpdateBuilder) Where(column string, condition string, value ...interface{}) *UpdateBuilder {
	b.whereBuilder.Where(column, condition, value...)
	return b
}

func (b *UpdateBuilder) OrWhere(column string, condition string, value ...interface{}) *UpdateBuilder {
	b.whereBuilder.OrWhere(column, condition, value...)
	return b
}

func (b *UpdateBuilder) WhereQuery(column string, condition string, q *Builder) *UpdateBuilder {
	b.whereBuilder.WhereQuery(column, condition, q)

	return b
}

func (b *UpdateBuilder) OrWhereQuery(column string, condition string, q *Builder) *UpdateBuilder {
	b.whereBuilder.OrWhereQuery(column, condition, q)

	return b
}

func (b *UpdateBuilder) WhereGroup(fn func(b *WhereBuilder) *WhereBuilder) *UpdateBuilder {
	b.whereBuilder.WhereGroup(fn)

	return b
}

func (b *UpdateBuilder) OrWhereGroup(fn func(b *WhereBuilder) *WhereBuilder) *UpdateBuilder {
	b.whereBuilder.OrWhereGroup(fn)

	return b
}

func (b *UpdateBuilder) WhereNot(fn func(b *WhereBuilder) *WhereBuilder) *UpdateBuilder {
	b.whereBuilder.WhereNot(fn)

	return b
}

func (b *UpdateBuilder) OrWhereNot(fn func(b *WhereBuilder) *WhereBuilder) *UpdateBuilder {
	b.whereBuilder.OrWhereNot(fn)

	return b
}

func (b *UpdateBuilder) WhereAny(columns []string, condition string, value interface{}) *UpdateBuilder {
	b.whereBuilder.WhereAny(columns, condition, value)
	return b
}

func (b *UpdateBuilder) WhereAll(columns []string, condition string, value interface{}) *UpdateBuilder {
	b.whereBuilder.WhereAll(columns, condition, value)
	return b
}

func (b *UpdateBuilder) Update(data map[string]interface{}) *UpdateBuilder {
	b.query.Values = data

	return b
}

func (u *UpdateBuilder) Build() (string, []interface{}) {
	// If there are conditions, add them to the query
	if len(*u.whereBuilder.query.Conditions) > 0 {
		*u.whereBuilder.query.ConditionGroups = append(*u.whereBuilder.query.ConditionGroups, structs.WhereGroup{
			Conditions:   *u.whereBuilder.query.Conditions,
			Operator:     consts.LogicalOperator_AND,
			IsDummyGroup: true,
		})
		u.whereBuilder.query.Conditions = &[]structs.Where{}
	}

	u.query.Query.Conditions = u.whereBuilder.query.Conditions
	u.query.Query.ConditionGroups = u.whereBuilder.query.ConditionGroups
	u.query.Query.Joins = u.joinBuilder.Joins
	u.query.Query.Order = u.orderByBuilder.Order

	query, values := u.dbBuilder.BuildUpdate(u.query)
	return query, values
}

func (b *UpdateBuilder) Join(table, my, condition, target string) *UpdateBuilder {
	b.joinBuilder.Join(table, my, condition, target)
	return b
}

func (b *UpdateBuilder) LeftJoin(table, my, condition, target string) *UpdateBuilder {
	b.joinBuilder.LeftJoin(table, my, condition, target)
	return b
}

func (b *UpdateBuilder) RightJoin(table, my, condition, target string) *UpdateBuilder {
	b.joinBuilder.RightJoin(table, my, condition, target)
	return b
}

func (b *UpdateBuilder) CrossJoin(table string) *UpdateBuilder {
	b.joinBuilder.CrossJoin(table)
	return b
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
