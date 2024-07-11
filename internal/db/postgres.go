package db

import "strings"

type PostgreSQLQueryBuilder struct{}

func (PostgreSQLQueryBuilder) Select(columns ...string) string {
	return "SELECT " + strings.Join(columns, ", ")
}

func (PostgreSQLQueryBuilder) From(table string) string {
	return "FROM " + table
}

func (PostgreSQLQueryBuilder) Where(condition string) string {
	return "WHERE " + condition
}

func (PostgreSQLQueryBuilder) Join(joinType, table, condition string) string {
	return joinType + " JOIN " + table + " ON " + condition
}

func (PostgreSQLQueryBuilder) OrderBy(orderBy ...string) string {
	return "ORDER BY " + strings.Join(orderBy, ", ")
}
