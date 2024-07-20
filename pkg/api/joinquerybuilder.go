package api

import "github.com/faciam-dev/goquent-query-builder/internal/query"

type JoinQueryBuilder struct {
	builder          *query.JoinBuilder
	joinQueryBuilder *JoinQueryBuilder
}

func (qb *JoinQueryBuilder) Join(table, my, condition, target string) *JoinQueryBuilder {
	qb.builder.Join(table, my, condition, target)
	return qb
}

func (qb *JoinQueryBuilder) LeftJoin(table, my, condition, target string) *JoinQueryBuilder {
	qb.builder.LeftJoin(table, my, condition, target)
	return qb
}

func (qb *JoinQueryBuilder) RightJoin(table, my, condition, target string) *JoinQueryBuilder {
	qb.builder.RightJoin(table, my, condition, target)
	return qb
}

func (qb *JoinQueryBuilder) CrossJoin(table string) *JoinQueryBuilder {
	qb.builder.CrossJoin(table)
	return qb
}

func (jb *JoinQueryBuilder) JoinQuery(table string, fn func(b *query.JoinClauseBuilder) *query.JoinClauseBuilder) *JoinQueryBuilder {
	jb.builder.JoinQuery(table, func(b *query.JoinClauseBuilder) *query.JoinClauseBuilder {
		return fn(b)
	})
	return jb
}

func (jb *JoinQueryBuilder) JoinSub(table string, qb *QueryBuilder) *JoinQueryBuilder {
	jb.builder.JoinSub(table, qb.builder)
	return jb
}

func (jb *JoinQueryBuilder) LeftJoinSub(table string, qb *QueryBuilder) *JoinQueryBuilder {
	jb.builder.LeftJoinSub(table, qb.builder)
	return jb
}

func (jb *JoinQueryBuilder) RightJoinSub(table string, qb *QueryBuilder) *JoinQueryBuilder {
	jb.builder.RightJoinSub(table, qb.builder)
	return jb
}
