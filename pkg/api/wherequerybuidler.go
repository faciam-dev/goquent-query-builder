package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type WhereQueryBuilder struct {
	builder *query.WhereBuilder
}

func NewWhereQueryBuilder(strategy db.QueryBuilderStrategy, cache cache.Cache) *WhereQueryBuilder {
	return &WhereQueryBuilder{
		builder: query.NewWhereBuilder(strategy, cache),
	}
}

func (wb *WhereQueryBuilder) Where(column string, condition string, value interface{}) *WhereQueryBuilder {
	switch v := value.(type) {
	case SelectBuilder:
		wb.builder.WhereQuery(column, condition, v.builder)
	case []interface{}:
		wb.builder.Where(column, condition, v...)
	default:
		wb.builder.Where(column, condition, value)
	}
	return wb
}

func (wb *WhereQueryBuilder) OrWhere(column string, condition string, value interface{}) *WhereQueryBuilder {
	switch v := value.(type) {
	case SelectBuilder:
		wb.builder.OrWhereQuery(column, condition, v.builder)
	case []interface{}:
		wb.builder.OrWhere(column, condition, v...)
	default:
		wb.builder.OrWhere(column, condition, value)
	}

	return wb
}

func (wb *WhereQueryBuilder) WhereQuery(column string, condition string, q *SelectBuilder) *WhereQueryBuilder {
	wb.builder.WhereQuery(column, condition, q.builder)
	return wb
}

func (wb *WhereQueryBuilder) OrWhereQuery(column string, condition string, q *SelectBuilder) *WhereQueryBuilder {
	wb.builder.OrWhereQuery(column, condition, q.builder)
	return wb
}

func (wb *WhereQueryBuilder) WhereRaw(raw string, value interface{}) *WhereQueryBuilder {
	wb.builder.WhereRaw(raw, value)
	return wb
}

func (wb *WhereQueryBuilder) OrWhereRaw(raw string, value interface{}) *WhereQueryBuilder {
	wb.builder.OrWhereRaw(raw, value)
	return wb
}

// WhereGroup
func (wb *WhereQueryBuilder) WhereGroup(fn func(wb *query.WhereBuilder) *query.WhereBuilder) *WhereQueryBuilder {
	wb.builder.WhereGroup(func(b *query.WhereBuilder) *query.WhereBuilder {
		return fn(b)
	})
	return wb
}

func (wb *WhereQueryBuilder) OrWhereGroup(fn func(wb *query.WhereBuilder) *query.WhereBuilder) *WhereQueryBuilder {
	wb.builder.OrWhereGroup(func(b *query.WhereBuilder) *query.WhereBuilder {
		return fn(b)
	})
	return wb
}

func (wb *WhereQueryBuilder) WhereNot(fn func(wb *query.WhereBuilder) *query.WhereBuilder) *WhereQueryBuilder {
	wb.builder.WhereNot(func(b *query.WhereBuilder) *query.WhereBuilder {
		return fn(b)
	})
	return wb
}

func (wb *WhereQueryBuilder) OrWhereNot(fn func(wb *query.WhereBuilder) *query.WhereBuilder) *WhereQueryBuilder {
	wb.builder.OrWhereNot(func(b *query.WhereBuilder) *query.WhereBuilder {
		return fn(b)
	})
	return wb
}

func (wb *WhereQueryBuilder) WhereAny(columns []string, condition string, value interface{}) *WhereQueryBuilder {
	wb.builder.WhereAny(columns, condition, value)
	return wb
}

func (wb *WhereQueryBuilder) WhereAll(columns []string, condition string, value interface{}) *WhereQueryBuilder {
	wb.builder.WhereAll(columns, condition, value)
	return wb
}
