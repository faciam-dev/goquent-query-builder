package api

import (
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

type JoinClauseBuilder struct {
	builder           *query.JoinClauseBuilder
	joinClauseBuilder *JoinClauseBuilder
}

func NewJoinClauseBuilder() *JoinClauseBuilder {
	return &JoinClauseBuilder{
		builder: query.NewJoinClauseBuilder(),
	}
}

func (qb *JoinClauseBuilder) On(my, condition, target string) *JoinClauseBuilder {
	qb.builder.On(my, condition, target)
	return qb
}

func (qb *JoinClauseBuilder) OrOn(my, condition, target string) *JoinClauseBuilder {
	qb.builder.OrOn(my, condition, target)
	return qb
}

func (qb *JoinClauseBuilder) Where(column, condition string, value interface{}) *JoinClauseBuilder {
	qb.builder.Where(column, condition, value)
	return qb
}

func (qb *JoinClauseBuilder) OrWhere(column, condition string, value interface{}) *JoinClauseBuilder {
	qb.builder.OrWhere(column, condition, value)
	return qb
}
