package db

type QueryBuilderStrategy interface {
	Select(columns ...string) string
	From(table string) string
	Where(condition string) string
	Join(joinType, table, condition string) string
	OrderBy(orderBy ...string) string
}

type QueryBuilder struct {
	strategy QueryBuilderStrategy
}

func NewQueryBuilder(strategy QueryBuilderStrategy) *QueryBuilder {
	return &QueryBuilder{strategy: strategy}
}

func (qb *QueryBuilder) Select(columns ...string) string {
	return qb.strategy.Select(columns...)
}

func (qb *QueryBuilder) From(table string) string {
	return qb.strategy.From(table)
}

func (qb *QueryBuilder) Where(condition string) string {
	return qb.strategy.Where(condition)
}

func (qb *QueryBuilder) Join(joinType, table, condition string) string {
	return qb.strategy.Join(joinType, table, condition)
}

func (qb *QueryBuilder) OrderBy(orderBy ...string) string {
	return qb.strategy.OrderBy(orderBy...)
}
