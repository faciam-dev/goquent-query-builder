package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type WhereQueryBuilder[T any, C any] struct {
	builder *query.WhereBuilder[C]
	parent  *T
}

func NewWhereQueryBuilder[T any, C any](strategy interfaces.QueryBuilderStrategy, cache cache.Cache) *WhereQueryBuilder[T, C] {
	return &WhereQueryBuilder[T, C]{
		builder: query.NewWhereBuilder[C](strategy, cache),
	}
}

// WhereSelectQueryBuilder is a type that represents a where select builder
type WhereSelectQueryBuilder = WhereQueryBuilder[SelectQueryBuilder, query.Builder]

// WhereInsertBuilder is a type that represents a where insert builder
type WhereUpdateQueryBuilder = WhereQueryBuilder[UpdateQueryBuilder, query.UpdateBuilder]

// WhereDeleteQueryBuilder is a type that represents a where delete builder
type WhereDeleteQueryBuilder = WhereQueryBuilder[DeleteQueryBuilder, query.DeleteBuilder]

func (b *WhereQueryBuilder[T, C]) SetParent(parent *T) *T {
	b.parent = parent

	return b.parent
}

// Where is a function that allows you to add a where condition
func (wb *WhereQueryBuilder[T, C]) Where(column string, condition string, value interface{}) *T {
	switch v := value.(type) {
	case *SelectQueryBuilder:
		wb.builder.WhereSubQuery(column, condition, v.builder)
	case []interface{}:
		wb.builder.Where(column, condition, v...)
	default:
		wb.builder.Where(column, condition, value)
	}
	return wb.parent
}

// OrWhere is a function that allows you to add a or where condition
func (wb *WhereQueryBuilder[T, C]) OrWhere(column string, condition string, value interface{}) *T {
	switch v := value.(type) {
	case *SelectQueryBuilder:
		wb.builder.OrWhereSubQuery(column, condition, v.builder)
	case []interface{}:
		wb.builder.OrWhere(column, condition, v...)
	default:
		wb.builder.OrWhere(column, condition, value)
	}

	return wb.parent
}

// WhereSubQuery is a function that allows you to add a where query condition
func (wb *WhereQueryBuilder[T, C]) WhereSubQuery(column string, condition string, qb *SelectQueryBuilder) *T {
	qb.builder.SetWhereBuilder(qb.WhereQueryBuilder.builder)
	qb.builder.SetJoinBuilder(qb.JoinQueryBuilder.builder)
	qb.builder.SetOrderByBuilder(qb.orderByQueryBuilder.builder)
	wb.builder.WhereSubQuery(column, condition, qb.builder)
	return wb.parent
}

// OrWhereSubQuery is a function that allows you to add a or where query condition
func (wb *WhereQueryBuilder[T, C]) OrWhereSubQuery(column string, condition string, qb *SelectQueryBuilder) *T {
	qb.builder.SetWhereBuilder(qb.WhereQueryBuilder.builder)
	qb.builder.SetJoinBuilder(qb.JoinQueryBuilder.builder)
	qb.builder.SetOrderByBuilder(qb.orderByQueryBuilder.builder)
	wb.builder.OrWhereSubQuery(column, condition, qb.builder)
	return wb.parent
}

// WhereRaw is a function that allows you to add a where raw condition
func (wb *WhereQueryBuilder[T, C]) WhereRaw(raw string, value interface{}) *T {
	wb.builder.WhereRaw(raw, value)
	return wb.parent
}

// OrWhereRaw is a function that allows you to add a or where raw condition
func (wb *WhereQueryBuilder[T, C]) OrWhereRaw(raw string, value interface{}) *T {
	wb.builder.OrWhereRaw(raw, value)
	return wb.parent
}

// WhereGroup is a function that allows you to group where conditions
func (wb *WhereQueryBuilder[T, C]) WhereGroup(fn func(wqb *WhereQueryBuilder[T, C])) *T {
	wb.builder.WhereGroup(func(b *query.WhereBuilder[C]) {
		fn(&WhereQueryBuilder[T, C]{builder: b, parent: wb.parent})
	})
	return wb.parent
}

// OrWhereGroup is a function that allows you to group or where conditions
func (wb *WhereQueryBuilder[T, C]) OrWhereGroup(fn func(wb *WhereQueryBuilder[T, C])) *T {
	wb.builder.OrWhereGroup(func(b *query.WhereBuilder[C]) {
		fn(&WhereQueryBuilder[T, C]{builder: b, parent: wb.parent})
	})
	return wb.parent
}

// WhereNot is a function that allows you to add a where not condition
func (wb *WhereQueryBuilder[T, C]) WhereNot(fn func(wb *WhereQueryBuilder[T, C])) *T {
	wb.builder.WhereNot(func(b *query.WhereBuilder[C]) {
		fn(&WhereQueryBuilder[T, C]{builder: b, parent: wb.parent})
	})
	return wb.parent
}

// OrWhereNot is a function that allows you to add a or where not condition
func (wb *WhereQueryBuilder[T, C]) OrWhereNot(fn func(wb *WhereQueryBuilder[T, C])) *T {
	wb.builder.OrWhereNot(func(b *query.WhereBuilder[C]) {
		fn(&WhereQueryBuilder[T, C]{builder: b, parent: wb.parent})
	})
	return wb.parent
}

// WhereAny is a function that allows you to add a where any condition
func (wb *WhereQueryBuilder[T, C]) WhereAny(columns []string, condition string, value interface{}) *T {
	wb.builder.WhereAny(columns, condition, value)
	return wb.parent
}

// WhereAll is a function that allows you to add a where all condition
func (wb *WhereQueryBuilder[T, C]) WhereAll(columns []string, condition string, value interface{}) *T {
	wb.builder.WhereAll(columns, condition, value)
	return wb.parent
}

// OrWhereAny is a function that allows you to add a or where any condition
func (wb *WhereQueryBuilder[T, C]) WhereIn(column string, values interface{}) *T {
	switch casted := values.(type) {
	case *SelectQueryBuilder:
		casted.builder.SetWhereBuilder(casted.WhereQueryBuilder.builder)
		casted.builder.SetJoinBuilder(casted.JoinQueryBuilder.builder)
		casted.builder.SetOrderByBuilder(casted.orderByQueryBuilder.builder)
		wb.builder.WhereInSubQuery(column, casted.builder)
	default:
		wb.builder.WhereIn(column, values)
	}
	return wb.parent
}

// OrWhereAll is a function that allows you to add a or where all condition
func (wb *WhereQueryBuilder[T, C]) WhereNotIn(column string, values interface{}) *T {
	switch casted := values.(type) {
	case *SelectQueryBuilder:
		casted.builder.SetWhereBuilder(casted.WhereQueryBuilder.builder)
		casted.builder.SetJoinBuilder(casted.JoinQueryBuilder.builder)
		casted.builder.SetOrderByBuilder(casted.orderByQueryBuilder.builder)
		wb.builder.WhereNotInSubQuery(column, casted.builder)
	default:
		wb.builder.WhereNotIn(column, values)
	}
	return wb.parent
}

// OrWhereIn is a function that allows you to add a or where in condition
func (wb *WhereQueryBuilder[T, C]) OrWhereIn(column string, values interface{}) *T {
	switch casted := values.(type) {
	case *SelectQueryBuilder:
		casted.builder.SetWhereBuilder(casted.WhereQueryBuilder.builder)
		casted.builder.SetJoinBuilder(casted.JoinQueryBuilder.builder)
		casted.builder.SetOrderByBuilder(casted.orderByQueryBuilder.builder)
		wb.builder.OrWhereInSubQuery(column, casted.builder)
	default:
		wb.builder.OrWhereIn(column, values)
	}
	return wb.parent
}

// OrWhereNotIn is a function that allows you to add a or where not in condition
func (wb *WhereQueryBuilder[T, C]) OrWhereNotIn(column string, values interface{}) *T {
	switch casted := values.(type) {
	case *SelectQueryBuilder:
		casted.builder.SetWhereBuilder(casted.WhereQueryBuilder.builder)
		casted.builder.SetJoinBuilder(casted.JoinQueryBuilder.builder)
		casted.builder.SetOrderByBuilder(casted.orderByQueryBuilder.builder)
		wb.builder.OrWhereNotInSubQuery(column, casted.builder)
	default:
		wb.builder.OrWhereNotIn(column, values)
	}
	return wb.parent
}

// WhereInSubQuery is a function that allows you to add a where in sub query condition
func (wb *WhereQueryBuilder[T, C]) WhereInSubQuery(column string, qb *SelectQueryBuilder) *T {
	qb.builder.SetWhereBuilder(qb.WhereQueryBuilder.builder)
	qb.builder.SetJoinBuilder(qb.JoinQueryBuilder.builder)
	qb.builder.SetOrderByBuilder(qb.orderByQueryBuilder.builder)
	wb.builder.WhereInSubQuery(column, qb.builder)
	return wb.parent
}

// WhereNotInSubQuery is a function that allows you to add a where not in sub query condition
func (wb *WhereQueryBuilder[T, C]) WhereNotInSubQuery(column string, qb *SelectQueryBuilder) *T {
	qb.builder.SetWhereBuilder(qb.WhereQueryBuilder.builder)
	qb.builder.SetJoinBuilder(qb.JoinQueryBuilder.builder)
	qb.builder.SetOrderByBuilder(qb.orderByQueryBuilder.builder)
	wb.builder.WhereNotInSubQuery(column, qb.builder)
	return wb.parent
}

// OrWhereInSubQuery is a function that allows you to add a or where in sub query condition
func (wb *WhereQueryBuilder[T, C]) OrWhereInSubQuery(column string, qb *SelectQueryBuilder) *T {
	qb.builder.SetWhereBuilder(qb.WhereQueryBuilder.builder)
	qb.builder.SetJoinBuilder(qb.JoinQueryBuilder.builder)
	qb.builder.SetOrderByBuilder(qb.orderByQueryBuilder.builder)
	wb.builder.OrWhereInSubQuery(column, qb.builder)
	return wb.parent
}

// OrWhereNotInSubQuery is a function that allows you to add a or where not in sub query condition
func (wb *WhereQueryBuilder[T, C]) OrWhereNotInSubQuery(column string, qb *SelectQueryBuilder) *T {
	qb.builder.SetWhereBuilder(qb.WhereQueryBuilder.builder)
	qb.builder.SetJoinBuilder(qb.JoinQueryBuilder.builder)
	qb.builder.SetOrderByBuilder(qb.orderByQueryBuilder.builder)
	wb.builder.OrWhereNotInSubQuery(column, qb.builder)
	return wb.parent
}

// WhereNull is a function that allows you to add a where null condition
func (wb *WhereQueryBuilder[T, C]) WhereNull(column string) *T {
	wb.builder.WhereNull(column)
	return wb.parent
}

// WhereNotNull is a function that allows you to add a where not null condition
func (wb *WhereQueryBuilder[T, C]) WhereNotNull(column string) *T {
	wb.builder.WhereNotNull(column)
	return wb.parent
}

// OrWhereNull is a function that allows you to add a or where null condition
func (wb *WhereQueryBuilder[T, C]) OrWhereNull(column string) *T {
	wb.builder.OrWhereNull(column)
	return wb.parent
}

// OrWhereNotNull is a function that allows you to add a or where not null condition
func (wb *WhereQueryBuilder[T, C]) OrWhereNotNull(column string) *T {
	wb.builder.OrWhereNotNull(column)
	return wb.parent
}

// WhereColumn is a function that allows you to add a where column condition
func (wb *WhereQueryBuilder[T, C]) WhereColumn(allColumns []string, column string, cond ...string) *T {
	operator := consts.Condition_EQUAL
	valueColumn := column
	if len(cond) > 0 {
		valueColumn = cond[0]
	}
	if len(cond) > 1 {
		operator = cond[0]
		valueColumn = cond[1]
	}

	wb.builder.WhereColumn(allColumns, column, operator, valueColumn)
	return wb.parent
}

// OrWhereColumn is a function that allows you to add a or where column condition
func (wb *WhereQueryBuilder[T, C]) OrWhereColumn(allColumns []string, column string, cond ...string) *T {
	operator := consts.Condition_EQUAL
	valueColumn := column
	if len(cond) > 0 {
		valueColumn = cond[0]
	}
	if len(cond) > 1 {
		operator = cond[0]
		valueColumn = cond[1]
	}

	wb.builder.OrWhereColumn(allColumns, column, operator, valueColumn)
	return wb.parent
}

// WhereColumns is a function that allows you to add a where columns condition
func (wb *WhereQueryBuilder[T, C]) WhereColumns(allColumns []string, columns [][]string) *T {
	wb.builder.WhereColumns(allColumns, columns)
	return wb.parent
}

// OrWhereColumns is a function that allows you to add a or where columns condition
func (wb *WhereQueryBuilder[T, C]) OrWhereColumns(allColumns []string, columns [][]string) *T {
	wb.builder.OrWhereColumns(allColumns, columns)
	return wb.parent
}

// WhereBetween is a function that allows you to add a where between condition
func (wb *WhereQueryBuilder[T, C]) WhereBetween(column string, min interface{}, max interface{}) *T {
	wb.builder.WhereBetween(column, min, max)
	return wb.parent
}

// OrWhereBetween is a function that allows you to add a or where between condition
func (wb *WhereQueryBuilder[T, C]) OrWhereBetween(column string, min interface{}, max interface{}) *T {
	wb.builder.OrWhereBetween(column, min, max)
	return wb.parent
}

// WhereNotBetween is a function that allows you to add a where not between condition
func (wb *WhereQueryBuilder[T, C]) WhereNotBetween(column string, min interface{}, max interface{}) *T {
	wb.builder.WhereNotBetween(column, min, max)
	return wb.parent
}

// OrWhereNotBetween is a function that allows you to add a or where not between condition
func (wb *WhereQueryBuilder[T, C]) OrWhereNotBetween(column string, min interface{}, max interface{}) *T {
	wb.builder.OrWhereNotBetween(column, min, max)
	return wb.parent
}

// WhereBetweenColumns is a function that allows you to add a where between columns condition
func (wb *WhereQueryBuilder[T, C]) WhereBetweenColumns(allColumns []string, column string, min string, max string) *T {
	wb.builder.WhereBetweenColumns(allColumns, column, min, max)
	return wb.parent
}

// OrWhereBetweenColumns is a function that allows you to add a or where between columns condition
func (wb *WhereQueryBuilder[T, C]) OrWhereBetweenColumns(allColumns []string, column string, min string, max string) *T {
	wb.builder.OrWhereBetweenColumns(allColumns, column, min, max)
	return wb.parent
}

// WhereNotBetweenColumns is a function that allows you to add a where not between columns condition
func (wb *WhereQueryBuilder[T, C]) WhereNotBetweenColumns(allColumns []string, column string, min string, max string) *T {
	wb.builder.WhereNotBetweenColumns(allColumns, column, min, max)
	return wb.parent
}

// OrWhereNotBetweenColumns is a function that allows you to add a or where not between columns condition
func (wb *WhereQueryBuilder[T, C]) OrWhereNotBetweenColumns(allColumns []string, column string, min string, max string) *T {
	wb.builder.OrWhereNotBetweenColumns(allColumns, column, min, max)
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) WhereExists(fn func(q *SelectQueryBuilder)) *T {
	wb.builder.WhereExists(func(b *query.Builder) {
		sqb := NewSelectQueryBuilder(b.GetStrategy(), b.GetCache())
		sqb.builder = b
		sqb.builder.SetJoinBuilder(sqb.JoinQueryBuilder.builder)
		sqb.builder.SetWhereBuilder(sqb.WhereQueryBuilder.builder)
		sqb.builder.SetOrderByBuilder(sqb.orderByQueryBuilder.builder)
		fn(sqb)
	})
	return wb.parent
}

// WhereDateQuery is a function that allows you to add a where date condition
func (wb *WhereQueryBuilder[T, C]) WhereExistsSubQuery(qb *SelectQueryBuilder) *T {
	qb.builder.SetWhereBuilder(qb.WhereQueryBuilder.builder)
	qb.builder.SetJoinBuilder(qb.JoinQueryBuilder.builder)
	qb.builder.SetOrderByBuilder(qb.orderByQueryBuilder.builder)
	wb.builder.WhereExistsQuery(qb.builder)
	return wb.parent
}

// OrWhereExists is a function that allows you to add a or where exists condition
func (wb *WhereQueryBuilder[T, C]) OrWhereExists(fn func(q *SelectQueryBuilder)) *T {
	wb.builder.OrWhereExists(func(b *query.Builder) {
		sqb := NewSelectQueryBuilder(b.GetStrategy(), b.GetCache())
		sqb.builder = b
		sqb.builder.SetJoinBuilder(sqb.JoinQueryBuilder.builder)
		sqb.builder.SetWhereBuilder(sqb.WhereQueryBuilder.builder)
		sqb.builder.SetOrderByBuilder(sqb.orderByQueryBuilder.builder)
		fn(sqb)
	})
	return wb.parent
}

// OrWhereExistsSubQuery is a function that allows you to add a or where exists condition
func (wb *WhereQueryBuilder[T, C]) OrWhereExistsSubQuery(qb *SelectQueryBuilder) *T {
	qb.builder.SetWhereBuilder(qb.WhereQueryBuilder.builder)
	qb.builder.SetJoinBuilder(qb.JoinQueryBuilder.builder)
	qb.builder.SetOrderByBuilder(qb.orderByQueryBuilder.builder)
	wb.builder.OrWhereExistsQuery(qb.builder)
	return wb.parent
}

// WhereNotExists is a function that allows you to add a where not exists condition
func (wb *WhereQueryBuilder[T, C]) WhereNotExists(fn func(q *SelectQueryBuilder)) *T {
	wb.builder.WhereNotExists(func(b *query.Builder) {
		sqb := NewSelectQueryBuilder(b.GetStrategy(), b.GetCache())
		sqb.builder = b
		sqb.builder.SetJoinBuilder(sqb.JoinQueryBuilder.builder)
		sqb.builder.SetWhereBuilder(sqb.WhereQueryBuilder.builder)
		sqb.builder.SetOrderByBuilder(sqb.orderByQueryBuilder.builder)
		fn(sqb)
	})
	return wb.parent
}

// WhereNotExistsQuery is a function that allows you to add a where not exists condition
func (wb *WhereQueryBuilder[T, C]) WhereNotExistsQuery(qb *SelectQueryBuilder) *T {
	qb.builder.SetWhereBuilder(qb.WhereQueryBuilder.builder)
	qb.builder.SetJoinBuilder(qb.JoinQueryBuilder.builder)
	qb.builder.SetOrderByBuilder(qb.orderByQueryBuilder.builder)
	wb.builder.WhereNotExistsQuery(qb.builder)
	return wb.parent
}

// OrWhereNotExists is a function that allows you to add a or where not exists condition
func (wb *WhereQueryBuilder[T, C]) OrWhereNotExists(fn func(q *SelectQueryBuilder)) *T {
	wb.builder.OrWhereNotExists(func(b *query.Builder) {
		sqb := NewSelectQueryBuilder(b.GetStrategy(), b.GetCache())
		sqb.builder = b
		sqb.builder.SetJoinBuilder(sqb.JoinQueryBuilder.builder)
		sqb.builder.SetWhereBuilder(sqb.WhereQueryBuilder.builder)
		sqb.builder.SetOrderByBuilder(sqb.orderByQueryBuilder.builder)
		fn(sqb)
	})
	return wb.parent
}

// OrWhereNotExistsQuery is a function that allows you to add a or where not exists condition
func (wb *WhereQueryBuilder[T, C]) OrWhereNotExistsQuery(qb *SelectQueryBuilder) *T {
	qb.builder.SetWhereBuilder(qb.WhereQueryBuilder.builder)
	qb.builder.SetJoinBuilder(qb.JoinQueryBuilder.builder)
	qb.builder.SetOrderByBuilder(qb.orderByQueryBuilder.builder)
	wb.builder.OrWhereNotExistsQuery(qb.builder)
	return wb.parent
}

// WhereFullText is a function that allows you to add a where full text condition
func (wb *WhereQueryBuilder[T, C]) WhereFullText(columns []string, value string, options map[string]interface{}) *T {
	wb.builder.WhereFullText(columns, value, options)
	return wb.parent
}

// OrWhereFullText is a function that allows you to add a or where full text condition
func (wb *WhereQueryBuilder[T, C]) OrWhereFullText(columns []string, value string, options map[string]interface{}) *T {
	wb.builder.OrWhereFullText(columns, value, options)
	return wb.parent
}

// WhereDate is a function that allows you to add a where date condition
func (wb *WhereQueryBuilder[T, C]) WhereDate(column string, cond string, date string) *T {
	wb.builder.WhereDate(column, cond, date)
	return wb.parent
}

// OrWhereDate is a function that allows you to add a or where date condition
func (wb *WhereQueryBuilder[T, C]) OrWhereDate(column string, cond string, date string) *T {
	wb.builder.OrWhereDate(column, cond, date)
	return wb.parent
}

// WhereTime is a function that allows you to add a where time condition
func (wb *WhereQueryBuilder[T, C]) WhereTime(column string, cond string, time string) *T {
	wb.builder.WhereTime(column, cond, time)
	return wb.parent
}

// OrWhereTime is a function that allows you to add a or where time condition
func (wb *WhereQueryBuilder[T, C]) OrWhereTime(column string, cond string, time string) *T {
	wb.builder.OrWhereTime(column, cond, time)
	return wb.parent
}

// WhereDay is a function that allows you to add a where day condition
func (wb *WhereQueryBuilder[T, C]) WhereDay(column string, cond string, day string) *T {
	wb.builder.WhereDay(column, cond, day)
	return wb.parent
}

// OrWhereDay is a function that allows you to add a or where day condition
func (wb *WhereQueryBuilder[T, C]) OrWhereDay(column string, cond string, day string) *T {
	wb.builder.OrWhereDay(column, cond, day)
	return wb.parent
}

// WhereMonth is a function that allows you to add a where month condition
func (wb *WhereQueryBuilder[T, C]) WhereMonth(column string, cond string, month string) *T {
	wb.builder.WhereMonth(column, cond, month)
	return wb.parent
}

// OrWhereMonth is a function that allows you to add a or where month condition
func (wb *WhereQueryBuilder[T, C]) OrWhereMonth(column string, cond string, month string) *T {
	wb.builder.OrWhereMonth(column, cond, month)
	return wb.parent
}

// WhereYear is a function that allows you to add a where year condition
func (wb *WhereQueryBuilder[T, C]) WhereYear(column string, cond string, year string) *T {
	wb.builder.WhereYear(column, cond, year)
	return wb.parent
}

// OrWhereYear is a function that allows you to add a or where year condition
func (wb *WhereQueryBuilder[T, C]) OrWhereYear(column string, cond string, year string) *T {
	wb.builder.OrWhereYear(column, cond, year)
	return wb.parent
}

// GetBuilder is a function that allows you to get the where builder
func (wb *WhereQueryBuilder[T, C]) GetBuilder() *query.WhereBuilder[C] {
	return wb.builder
}
