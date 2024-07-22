package query

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
)

type WhereBuilder struct {
	dbBuilder   db.QueryBuilderStrategy
	cache       cache.Cache
	query       *structs.Query
	whereValues []interface{}
}

func NewWhereBuilder(strategy db.QueryBuilderStrategy, cache cache.Cache) *WhereBuilder {
	return &WhereBuilder{
		dbBuilder: strategy,
		cache:     cache,
		query: &structs.Query{
			Conditions:      &[]structs.Where{},
			ConditionGroups: &[]structs.WhereGroup{},
		},
	}
}

// Where adds a where clause with AND operator
func (b *WhereBuilder) Where(column string, condition string, value ...interface{}) *WhereBuilder {
	*b.query.Conditions = append(*b.query.Conditions, structs.Where{
		Column:    column,
		Condition: condition,
		Value:     value,
		Operator:  consts.LogicalOperator_AND,
	})
	b.whereValues = append(b.whereValues, value...)
	return b
}

// OrWhere adds a where clause with OR operator
func (b *WhereBuilder) OrWhere(column string, condition string, value ...interface{}) *WhereBuilder {
	*b.query.Conditions = append(*b.query.Conditions, structs.Where{
		Column:    column,
		Condition: condition,
		Value:     value,
		Operator:  consts.LogicalOperator_OR,
	})
	b.whereValues = append(b.whereValues, value...)
	return b
}

// WhereRaw adds a raw where clause with AND operator
func (b *WhereBuilder) WhereRaw(column string, value ...interface{}) *WhereBuilder {
	*b.query.Conditions = append(*b.query.Conditions, structs.Where{
		Value:    value,
		Raw:      column,
		Operator: consts.LogicalOperator_AND,
	})
	b.whereValues = append(b.whereValues, value...)
	return b
}

// OrWhereRaw adds a raw where clause with OR operator
func (b *WhereBuilder) OrWhereRaw(column string, value ...interface{}) *WhereBuilder {
	*b.query.Conditions = append(*b.query.Conditions, structs.Where{
		Value:    value,
		Raw:      column,
		Operator: consts.LogicalOperator_OR,
	})
	b.whereValues = append(b.whereValues, value...)
	return b
}

// WhereQuery adds a where clause with AND operator
func (b *WhereBuilder) WhereQuery(column string, condition string, q *Builder) *WhereBuilder {
	return b.whereOrOrWhereQuery(column, condition, q, consts.LogicalOperator_AND)
}

// OrWhereQuery adds a where clause with OR operator
func (b *WhereBuilder) OrWhereQuery(column string, condition string, q *Builder) *WhereBuilder {
	return b.whereOrOrWhereQuery(column, condition, q, consts.LogicalOperator_OR)
}

// whereOrOrWhereQuery adds a where clause with AND or OR operator
func (b *WhereBuilder) whereOrOrWhereQuery(column string, condition string, q *Builder, operator int) *WhereBuilder {
	*q.whereBuilder.query.ConditionGroups = append(*q.whereBuilder.query.ConditionGroups, structs.WhereGroup{
		Conditions:   *q.whereBuilder.query.Conditions,
		IsDummyGroup: true,
	})

	sq := &structs.Query{
		ConditionGroups: q.whereBuilder.query.ConditionGroups,
		Table:           structs.Table{Name: q.selectQuery.Table},
		Columns:         q.selectQuery.Columns,
		Joins:           q.joinBuilder.Joins,
		Order:           q.orderByBuilder.Order,
	}

	args := &structs.Where{
		Column:    column,
		Condition: condition,
		Query:     sq,
		Operator:  operator,
	}

	_, value := b.BuildSq(sq)

	*b.query.Conditions = append(*b.query.Conditions, *args)
	b.whereValues = append(b.whereValues, value...)
	return b
}

// WhereGroup adds a where group with AND operator
func (b *WhereBuilder) WhereGroup(fn func(b *WhereBuilder) *WhereBuilder) *WhereBuilder {
	if len(*b.query.Conditions) > 0 {
		*b.query.ConditionGroups = append(*b.query.ConditionGroups, structs.WhereGroup{
			Conditions:   *b.query.Conditions,
			Operator:     consts.LogicalOperator_AND,
			IsDummyGroup: true,
		})
		*b.query.Conditions = []structs.Where{}
	}

	cQ := fn(b)

	*b.query.ConditionGroups = append(*b.query.ConditionGroups, structs.WhereGroup{
		Conditions: *cQ.query.Conditions,
		Subgroups:  []structs.WhereGroup{},
		Operator:   consts.LogicalOperator_AND,
	})
	*cQ.query.Conditions = []structs.Where{}

	return b
}

// OrWhereGroup adds a where group with OR operator
func (b *WhereBuilder) OrWhereGroup(fn func(b *WhereBuilder) *WhereBuilder) *WhereBuilder {
	if len(*b.query.Conditions) > 0 {
		*b.query.ConditionGroups = append(*b.query.ConditionGroups, structs.WhereGroup{
			Conditions:   *b.query.Conditions,
			Operator:     consts.LogicalOperator_OR,
			IsDummyGroup: true,
		})
		*b.query.Conditions = []structs.Where{}
	}

	cQ := fn(b)

	*b.query.ConditionGroups = append(*b.query.ConditionGroups, structs.WhereGroup{
		Conditions: *cQ.query.Conditions,
		Subgroups:  []structs.WhereGroup{},
		Operator:   consts.LogicalOperator_OR,
	})
	*cQ.query.Conditions = []structs.Where{}

	return b
}

// BuildSq builds the query and returns the query string and values
func (b *WhereBuilder) BuildSq(sq *structs.Query) (string, []interface{}) {
	cacheKey := generateCacheKey(sq)

	if cachedQuery, found := b.cache.Get(cacheKey); found {
		values := []interface{}{}
		values = append(values, b.whereValues...)
		return cachedQuery, values
	}

	query, values := b.dbBuilder.Build(cacheKey, sq)

	b.cache.Set(cacheKey, query)

	return query, values
}

func (b *WhereBuilder) GetQuery() *structs.Query {
	return b.query
}
