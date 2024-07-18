package api

import "github.com/faciam-dev/goquent-query-builder/internal/query"

type JoinQueryBuilder struct {
	builder *query.JoinBuilder
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
