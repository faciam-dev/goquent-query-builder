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

func (wb *WhereQueryBuilder[T, C]) WhereQuery(column string, condition string, q *SelectBuilder) *T {
	wb.builder.WhereQuery(column, condition, q.builder)
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) OrWhereQuery(column string, condition string, q *SelectBuilder) *T {
	wb.builder.OrWhereQuery(column, condition, q.builder)
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) WhereRaw(raw string, value interface{}) *T {
	wb.builder.WhereRaw(raw, value)
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) OrWhereRaw(raw string, value interface{}) *T {
	wb.builder.OrWhereRaw(raw, value)
	return wb.parent
}

// WhereGroup
func (wb *WhereQueryBuilder[T, C]) WhereGroup(fn func(wb *query.WhereBuilder[C]) *query.WhereBuilder[C]) *T {
	wb.builder.WhereGroup(func(b *query.WhereBuilder[C]) *query.WhereBuilder[C] {
		return fn(b)
	})
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) OrWhereGroup(fn func(wb *query.WhereBuilder[C]) *query.WhereBuilder[C]) *T {
	wb.builder.OrWhereGroup(func(b *query.WhereBuilder[C]) *query.WhereBuilder[C] {
		return fn(b)
	})
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) WhereNot(fn func(wb *query.WhereBuilder[C]) *query.WhereBuilder[C]) *T {
	wb.builder.WhereNot(func(b *query.WhereBuilder[C]) *query.WhereBuilder[C] {
		return fn(b)
	})
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) OrWhereNot(fn func(wb *query.WhereBuilder[C]) *query.WhereBuilder[C]) *T {
	wb.builder.OrWhereNot(func(b *query.WhereBuilder[C]) *query.WhereBuilder[C] {
		return fn(b)
	})
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) WhereAny(columns []string, condition string, value interface{}) *T {
	wb.builder.WhereAny(columns, condition, value)
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) WhereAll(columns []string, condition string, value interface{}) *T {
	wb.builder.WhereAll(columns, condition, value)
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) WhereIn(column string, values interface{}) *T {
	wb.builder.WhereIn(column, values)
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) WhereNotIn(column string, values interface{}) *T {
	wb.builder.WhereNotIn(column, values)
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) OrWhereIn(column string, values interface{}) *T {
	wb.builder.OrWhereIn(column, values)
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) OrWhereNotIn(column string, values interface{}) *T {
	wb.builder.OrWhereNotIn(column, values)
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) WhereInSubQuery(column string, q *SelectBuilder) *T {
	wb.builder.WhereInSubQuery(column, q.builder)
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) WhereNotInSubQuery(column string, q *SelectBuilder) *T {
	wb.builder.WhereNotInSubQuery(column, q.builder)
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) OrWhereInSubQuery(column string, q *SelectBuilder) *T {
	wb.builder.OrWhereInSubQuery(column, q.builder)
	return wb.parent
}

func (wb *WhereQueryBuilder[T, C]) OrWhereNotInSubQuery(column string, q *SelectBuilder) *T {
	wb.builder.OrWhereNotInSubQuery(column, q.builder)
	return wb.parent
}
