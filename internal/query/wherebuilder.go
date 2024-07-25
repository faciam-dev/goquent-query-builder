package query

import (
	"log"
	"reflect"
	"time"

	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/sliceutils"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
)

type WhereBuilder[T any] struct {
	dbBuilder   db.QueryBuilderStrategy
	cache       cache.Cache
	query       *structs.Query
	whereValues []interface{}
	parent      *T
}

func NewWhereBuilder[T any](strategy db.QueryBuilderStrategy, cache cache.Cache) *WhereBuilder[T] {
	return &WhereBuilder[T]{
		dbBuilder: strategy,
		cache:     cache,
		query: &structs.Query{
			Conditions:      &[]structs.Where{},
			ConditionGroups: &[]structs.WhereGroup{},
		},
		whereValues: []interface{}{},
	}
}

func (b *WhereBuilder[T]) SetParent(parent *T) *T {
	b.parent = parent

	return b.parent
}

// Where adds a where clause with AND operator
func (b *WhereBuilder[T]) Where(column string, condition string, value ...interface{}) *T {
	*b.query.Conditions = append(*b.query.Conditions, structs.Where{
		Column:    column,
		Condition: condition,
		Value:     value,
		Operator:  consts.LogicalOperator_AND,
	})
	b.whereValues = append(b.whereValues, value...)
	return b.parent
}

// OrWhere adds a where clause with OR operator
func (b *WhereBuilder[T]) OrWhere(column string, condition string, value ...interface{}) *T {
	*b.query.Conditions = append(*b.query.Conditions, structs.Where{
		Column:    column,
		Condition: condition,
		Value:     value,
		Operator:  consts.LogicalOperator_OR,
	})
	b.whereValues = append(b.whereValues, value...)
	return b.parent
}

// WhereRaw adds a raw where clause with AND operator
func (b *WhereBuilder[T]) WhereRaw(column string, value ...interface{}) *T {
	*b.query.Conditions = append(*b.query.Conditions, structs.Where{
		Value:    value,
		Raw:      column,
		Operator: consts.LogicalOperator_AND,
	})
	b.whereValues = append(b.whereValues, value...)
	return b.parent
}

// OrWhereRaw adds a raw where clause with OR operator
func (b *WhereBuilder[T]) OrWhereRaw(column string, value ...interface{}) *T {
	*b.query.Conditions = append(*b.query.Conditions, structs.Where{
		Value:    value,
		Raw:      column,
		Operator: consts.LogicalOperator_OR,
	})
	b.whereValues = append(b.whereValues, value...)
	return b.parent
}

// WhereQuery adds a where clause with AND operator
func (b *WhereBuilder[T]) WhereQuery(column string, condition string, q *Builder) *T {
	return b.whereOrOrWhereQuery(column, condition, q, consts.LogicalOperator_AND)
}

// OrWhereQuery adds a where clause with OR operator
func (b *WhereBuilder[T]) OrWhereQuery(column string, condition string, q *Builder) *T {
	return b.whereOrOrWhereQuery(column, condition, q, consts.LogicalOperator_OR)
}

// whereOrOrWhereQuery adds a where clause with AND or OR operator
func (b *WhereBuilder[T]) whereOrOrWhereQuery(column string, condition string, q *Builder, operator int) *T {
	*q.WhereBuilder.query.ConditionGroups = append(*q.WhereBuilder.query.ConditionGroups, structs.WhereGroup{
		Conditions:   *q.WhereBuilder.query.Conditions,
		IsDummyGroup: true,
	})

	sq := &structs.Query{
		ConditionGroups: q.WhereBuilder.query.ConditionGroups,
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
	return b.parent
}

// WhereGroup adds a where group with AND operator
func (b *WhereBuilder[T]) WhereGroup(fn func(b *WhereBuilder[T]) *WhereBuilder[T]) *T {
	return b.addWhereGroup(fn, consts.LogicalOperator_AND, false)
}

// OrWhereGroup adds a where group with OR operator
func (b *WhereBuilder[T]) OrWhereGroup(fn func(b *WhereBuilder[T]) *WhereBuilder[T]) *T {
	return b.addWhereGroup(fn, consts.LogicalOperator_OR, false)
}

// WhereNot adds a not where group with AND operator
func (b *WhereBuilder[T]) WhereNot(fn func(b *WhereBuilder[T]) *WhereBuilder[T]) *T {
	return b.addWhereGroup(fn, consts.LogicalOperator_AND, true)
}

// OrWhereNot adds a not where group with OR operator
func (b *WhereBuilder[T]) OrWhereNot(fn func(b *WhereBuilder[T]) *WhereBuilder[T]) *T {
	return b.addWhereGroup(fn, consts.LogicalOperator_OR, true)
}

// addWhereGroup adds a where group with the specified operator
func (b *WhereBuilder[T]) addWhereGroup(fn func(b *WhereBuilder[T]) *WhereBuilder[T], operator int, isNot bool) *T {
	if len(*b.query.Conditions) > 0 {
		*b.query.ConditionGroups = append(*b.query.ConditionGroups, structs.WhereGroup{
			Conditions:   *b.query.Conditions,
			Operator:     operator,
			IsDummyGroup: true,
			IsNot:        false,
		})
		*b.query.Conditions = []structs.Where{}
	}

	cQ := fn(b)

	*b.query.ConditionGroups = append(*b.query.ConditionGroups, structs.WhereGroup{
		Conditions: *cQ.query.Conditions,
		Subgroups:  []structs.WhereGroup{},
		Operator:   operator,
		IsNot:      isNot,
	})
	*cQ.query.Conditions = []structs.Where{}

	return b.parent
}

// WhereAny adds where clauses with AND operator
func (b *WhereBuilder[T]) WhereAll(columns []string, condition string, value interface{}) *T {
	return b.addWhereConditions(columns, condition, value, consts.LogicalOperator_AND)
}

// OrWhereAny adds where clauses with OR operator
func (b *WhereBuilder[T]) WhereAny(columns []string, condition string, value interface{}) *T {
	return b.addWhereConditions(columns, condition, value, consts.LogicalOperator_OR)
}

func (b *WhereBuilder[T]) addWhereConditions(columns []string, condition string, value interface{}, operator int) *T {
	// already have conditions, add them to the query
	if len(*b.query.Conditions) > 0 {
		*b.query.ConditionGroups = append(*b.query.ConditionGroups, structs.WhereGroup{
			Conditions:   *b.query.Conditions,
			Operator:     operator,
			IsDummyGroup: true,
			IsNot:        false,
		})
		*b.query.Conditions = []structs.Where{}
	}

	conditions := []structs.Where{}
	for _, c := range columns {
		conditions = append(conditions, structs.Where{
			Column:    c,
			Condition: condition,
			Value:     []interface{}{value},
			Operator:  operator,
		})
		b.whereValues = append(b.whereValues, value)
	}

	*b.query.ConditionGroups = append(*b.query.ConditionGroups, structs.WhereGroup{
		Conditions: conditions,
		Subgroups:  []structs.WhereGroup{},
		Operator:   consts.LogicalOperator_AND,
	})

	return b.parent
}

// WhereIn adds a where in clause with AND operator
func (b *WhereBuilder[T]) WhereIn(column string, values interface{}) *T {
	return b.addWhereIn(column, consts.LogicalOperator_AND, consts.Condition_IN, values)
}

// WhereNotIn adds a not where in clause with AND operator
func (b *WhereBuilder[T]) WhereNotIn(column string, values interface{}) *T {
	return b.addWhereIn(column, consts.LogicalOperator_AND, consts.Condition_NOT_IN, values)
}

// OrWhereIn adds a where in clause with OR operator
func (b *WhereBuilder[T]) OrWhereIn(column string, values interface{}) *T {
	return b.addWhereIn(column, consts.LogicalOperator_OR, consts.Condition_IN, values)
}

// OrWhereNotIn adds a not where in clause with OR operator
func (b *WhereBuilder[T]) OrWhereNotIn(column string, values interface{}) *T {
	return b.addWhereIn(column, consts.LogicalOperator_OR, consts.Condition_NOT_IN, values)
}

// addWhereIn adds a where in clause with the specified operator
func (b *WhereBuilder[T]) addWhereIn(column string, operator int, condition string, values interface{}) *T {

	switch casted := values.(type) {
	case []interface{}:
		*b.query.Conditions = append(*b.query.Conditions, structs.Where{
			Value:     casted,
			Operator:  operator,
			Column:    column,
			Condition: condition,
		})
	case []bool, []int, []int32, []int64, []uint, []uint32, []uint64,
		[]float32, []float64, []string, []time.Time:
		nValues := sliceutils.ToInterfaceSlice(casted)
		*b.query.Conditions = append(*b.query.Conditions, structs.Where{
			Value:     nValues,
			Operator:  operator,
			Column:    column,
			Condition: condition,
		})

		b.whereValues = append(b.whereValues, nValues...)

	case *Builder:
		return b.addWhereInSubQuery(column, operator, condition, casted)
	default:
		log.Default().Printf("type: %T\n", reflect.TypeOf(values))
		log.Default().Println("values: ", values)
		//panic("Invalid type for values")
	}

	return b.parent
}

// WhereIn adds a where in clause with AND operator
func (b *WhereBuilder[T]) WhereInSubQuery(column string, q *Builder) *T {
	return b.addWhereInSubQuery(column, consts.LogicalOperator_AND, consts.Condition_IN, q)
}

// WhereNotIn adds a not where in clause with AND operator
func (b *WhereBuilder[T]) WhereNotInSubQuery(column string, q *Builder) *T {
	return b.addWhereInSubQuery(column, consts.LogicalOperator_AND, consts.Condition_NOT_IN, q)
}

// OrWhereIn adds a where in clause with OR operator
func (b *WhereBuilder[T]) OrWhereInSubQuery(column string, q *Builder) *T {
	return b.addWhereInSubQuery(column, consts.LogicalOperator_OR, consts.Condition_IN, q)
}

// OrWhereNotIn adds a not where in clause with OR operator
func (b *WhereBuilder[T]) OrWhereNotInSubQuery(column string, q *Builder) *T {
	return b.addWhereInSubQuery(column, consts.LogicalOperator_OR, consts.Condition_NOT_IN, q)
}

// addWhereIn adds a where in clause with the specified operator
func (b *WhereBuilder[T]) addWhereInSubQuery(column string, operator int, condition string, q *Builder) *T {
	*q.WhereBuilder.query.ConditionGroups = append(*q.WhereBuilder.query.ConditionGroups, structs.WhereGroup{
		Conditions:   *q.WhereBuilder.query.Conditions,
		IsDummyGroup: true,
	})

	sq := &structs.Query{
		ConditionGroups: q.WhereBuilder.query.ConditionGroups,
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
	return b.parent
}

func (b *WhereBuilder[T]) WhereNull(column string) *T {
	return b.addWhereNull(column, consts.LogicalOperator_AND, consts.Condition_IS_NULL)
}

func (b *WhereBuilder[T]) WhereNotNull(column string) *T {
	return b.addWhereNull(column, consts.LogicalOperator_AND, consts.Condition_IS_NOT_NULL)
}

func (b *WhereBuilder[T]) OrWhereNull(column string) *T {
	return b.addWhereNull(column, consts.LogicalOperator_OR, consts.Condition_IS_NULL)
}

func (b *WhereBuilder[T]) OrWhereNotNull(column string) *T {
	return b.addWhereNull(column, consts.LogicalOperator_OR, consts.Condition_IS_NOT_NULL)
}

func (b *WhereBuilder[T]) addWhereNull(column string, operator int, condition string) *T {
	*b.query.Conditions = append(*b.query.Conditions, structs.Where{
		Column:    column,
		Condition: condition,
		Operator:  operator,
		Value:     nil,
	})

	return b.parent
}

// WhereRawGroup adds a raw where group with AND operator
// BuildSq builds the query and returns the query string and values
func (b *WhereBuilder[T]) BuildSq(sq *structs.Query) (string, []interface{}) {
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

func (b *WhereBuilder[T]) GetQuery() *structs.Query {
	return b.query
}
