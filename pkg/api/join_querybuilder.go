package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type JoinQueryBuilder[T any, C any] struct {
	builder *query.JoinBuilder[C]
	parent  *T
}

func NewJoinQueryBuilder[T any, C any](strategy db.QueryBuilderStrategy, cache cache.Cache) *JoinQueryBuilder[T, C] {
	return &JoinQueryBuilder[T, C]{
		builder: query.NewJoinBuilder[C](strategy, cache),
	}
}

func (b *JoinQueryBuilder[T, C]) SetParent(parent *T) *T {
	b.parent = parent

	return b.parent
}

func (qb *JoinQueryBuilder[T, C]) Join(table, my, condition, target string) *T {
	qb.builder.Join(table, my, condition, target)
	return qb.parent
}

func (qb *JoinQueryBuilder[T, C]) LeftJoin(table, my, condition, target string) *T {
	qb.builder.LeftJoin(table, my, condition, target)
	return qb.parent
}

func (qb *JoinQueryBuilder[T, C]) RightJoin(table, my, condition, target string) *T {
	qb.builder.RightJoin(table, my, condition, target)
	return qb.parent
}

func (qb *JoinQueryBuilder[T, C]) CrossJoin(table string) *T {
	qb.builder.CrossJoin(table)
	return qb.parent
}

func (jb *JoinQueryBuilder[T, C]) JoinQuery(table string, fn func(b *query.JoinClauseBuilder) *query.JoinClauseBuilder) *T {
	jb.builder.JoinQuery(table, func(b *query.JoinClauseBuilder) *query.JoinClauseBuilder {
		return fn(b)
	})
	return jb.parent
}

func (jb *JoinQueryBuilder[T, C]) JoinSub(qb *SelectQueryBuilder, alias, my, condition, target string) *T {
	jb.builder.JoinSub(qb.builder, alias, my, condition, target)
	return jb.parent
}

func (jb *JoinQueryBuilder[T, C]) LeftJoinSub(qb *SelectQueryBuilder, alias, my, condition, target string) *T {
	jb.builder.LeftJoinSub(qb.builder, alias, my, condition, target)
	return jb.parent
}

func (jb *JoinQueryBuilder[T, C]) RightJoinSub(qb *SelectQueryBuilder, alias, my, condition, target string) *T {
	jb.builder.RightJoinSub(qb.builder, alias, my, condition, target)
	return jb.parent
}

func (jb *JoinQueryBuilder[T, C]) JoinLateral(qb *SelectQueryBuilder, alias string) *T {
	jb.builder.JoinLateral(qb.builder, alias)
	return jb.parent
}

func (jb *JoinQueryBuilder[T, C]) LeftJoinLateral(qb *SelectQueryBuilder, alias string) *T {
	jb.builder.LeftJoinLateral(qb.builder, alias)
	return jb.parent
}
