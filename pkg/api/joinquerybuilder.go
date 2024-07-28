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

func (jb *JoinQueryBuilder) JoinQuery(table string, fn func(b *query.JoinClauseBuilder) *query.JoinClauseBuilder) *JoinQueryBuilder {
	jb.builder.JoinQuery(table, func(b *query.JoinClauseBuilder) *query.JoinClauseBuilder {
		return fn(b)
	})
	return jb
}

func (jb *JoinQueryBuilder) JoinSub(qb *SelectBuilder, alias, my, condition, target string) *JoinQueryBuilder {
	jb.builder.JoinSub(qb.builder, alias, my, condition, target)
	return jb
}

func (jb *JoinQueryBuilder) LeftJoinSub(qb *SelectBuilder, alias, my, condition, target string) *JoinQueryBuilder {
	jb.builder.LeftJoinSub(qb.builder, alias, my, condition, target)
	return jb
}

func (jb *JoinQueryBuilder) RightJoinSub(qb *SelectBuilder, alias, my, condition, target string) *JoinQueryBuilder {
	jb.builder.RightJoinSub(qb.builder, alias, my, condition, target)
	return jb
}

func (jb *JoinQueryBuilder) JoinLateral(qb *SelectBuilder, alias string) *JoinQueryBuilder {
	jb.builder.JoinLateral(qb.builder, alias)
	return jb
}

func (jb *JoinQueryBuilder) LeftJoinLateral(qb *SelectBuilder, alias string) *JoinQueryBuilder {
	jb.builder.LeftJoinLateral(qb.builder, alias)
	return jb
}
