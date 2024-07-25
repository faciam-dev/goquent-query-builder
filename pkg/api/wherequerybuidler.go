package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type WhereQueryBuilder[T any, C any] struct {
	builder *query.WhereBuilder[C]
	parent  *T
}

func NewWhereQueryBuilder[T any, C any](strategy db.QueryBuilderStrategy, cache cache.Cache) *WhereQueryBuilder[T, C] {
	return &WhereQueryBuilder[T, C]{
		builder: query.NewWhereBuilder[C](strategy, cache),
	}
}

func (b *WhereQueryBuilder[T, C]) SetParent(parent *T) *T {
	b.parent = parent

	return b.parent
}

// Where is a function that allows you to add a where condition
func (wb *WhereQueryBuilder[T, C]) Where(column string, condition string, value interface{}) *T {
	switch v := value.(type) {
	case SelectBuilder:
		wb.builder.WhereQuery(column, condition, v.builder)
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
	case SelectBuilder:
		wb.builder.OrWhereQuery(column, condition, v.builder)
	case []interface{}:
		wb.builder.OrWhere(column, condition, v...)
	default:
		wb.builder.OrWhere(column, condition, value)
	}

	return wb.parent
}

// WhereQuery is a function that allows you to add a where query condition
func (wb *WhereQueryBuilder[T, C]) WhereQuery(column string, condition string, q *SelectBuilder) *T {
	wb.builder.WhereQuery(column, condition, q.builder)
	return wb.parent
}

// OrWhereQuery is a function that allows you to add a or where query condition
func (wb *WhereQueryBuilder[T, C]) OrWhereQuery(column string, condition string, q *SelectBuilder) *T {
	wb.builder.OrWhereQuery(column, condition, q.builder)
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
func (wb *WhereQueryBuilder[T, C]) WhereGroup(fn func(wb *query.WhereBuilder[C]) *query.WhereBuilder[C]) *T {
	wb.builder.WhereGroup(func(b *query.WhereBuilder[C]) *query.WhereBuilder[C] {
		return fn(b)
	})
	return wb.parent
}

// OrWhereGroup is a function that allows you to group or where conditions
func (wb *WhereQueryBuilder[T, C]) OrWhereGroup(fn func(wb *query.WhereBuilder[C]) *query.WhereBuilder[C]) *T {
	wb.builder.OrWhereGroup(func(b *query.WhereBuilder[C]) *query.WhereBuilder[C] {
		return fn(b)
	})
	return wb.parent
}

// WhereNot is a function that allows you to add a where not condition
func (wb *WhereQueryBuilder[T, C]) WhereNot(fn func(wb *query.WhereBuilder[C]) *query.WhereBuilder[C]) *T {
	wb.builder.WhereNot(func(b *query.WhereBuilder[C]) *query.WhereBuilder[C] {
		return fn(b)
	})
	return wb.parent
}

// OrWhereNot is a function that allows you to add a or where not condition
func (wb *WhereQueryBuilder[T, C]) OrWhereNot(fn func(wb *query.WhereBuilder[C]) *query.WhereBuilder[C]) *T {
	wb.builder.OrWhereNot(func(b *query.WhereBuilder[C]) *query.WhereBuilder[C] {
		return fn(b)
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
	wb.builder.WhereIn(column, values)
	return wb.parent
}

// OrWhereAll is a function that allows you to add a or where all condition
func (wb *WhereQueryBuilder[T, C]) WhereNotIn(column string, values interface{}) *T {
	wb.builder.WhereNotIn(column, values)
	return wb.parent
}

// OrWhereIn is a function that allows you to add a or where in condition
func (wb *WhereQueryBuilder[T, C]) OrWhereIn(column string, values interface{}) *T {
	wb.builder.OrWhereIn(column, values)
	return wb.parent
}

// OrWhereNotIn is a function that allows you to add a or where not in condition
func (wb *WhereQueryBuilder[T, C]) OrWhereNotIn(column string, values interface{}) *T {
	wb.builder.OrWhereNotIn(column, values)
	return wb.parent
}

// WhereInSubQuery is a function that allows you to add a where in sub query condition
func (wb *WhereQueryBuilder[T, C]) WhereInSubQuery(column string, q *SelectBuilder) *T {
	wb.builder.WhereInSubQuery(column, q.builder)
	return wb.parent
}

// WhereNotInSubQuery is a function that allows you to add a where not in sub query condition
func (wb *WhereQueryBuilder[T, C]) WhereNotInSubQuery(column string, q *SelectBuilder) *T {
	wb.builder.WhereNotInSubQuery(column, q.builder)
	return wb.parent
}

// OrWhereInSubQuery is a function that allows you to add a or where in sub query condition
func (wb *WhereQueryBuilder[T, C]) OrWhereInSubQuery(column string, q *SelectBuilder) *T {
	wb.builder.OrWhereInSubQuery(column, q.builder)
	return wb.parent
}

// OrWhereNotInSubQuery is a function that allows you to add a or where not in sub query condition
func (wb *WhereQueryBuilder[T, C]) OrWhereNotInSubQuery(column string, q *SelectBuilder) *T {
	wb.builder.OrWhereNotInSubQuery(column, q.builder)
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) WhereNull(column string) *T {
	wb.builder.WhereNull(column)
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) WhereNotNull(column string) *T {
	wb.builder.WhereNotNull(column)
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) OrWhereNull(column string) *T {
	wb.builder.OrWhereNull(column)
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) OrWhereNotNull(column string) *T {
	wb.builder.OrWhereNotNull(column)
	return wb.parent
}
