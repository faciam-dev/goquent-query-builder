package db

import "strings"

type MySQLQueryBuilder struct{}

func (MySQLQueryBuilder) Select(columns ...string) string {
	return "SELECT " + strings.Join(columns, ", ")
}

func (MySQLQueryBuilder) From(table string) string {
	return "FROM " + table
}

func (MySQLQueryBuilder) Where(condition string) string {
	return "WHERE " + condition
}

func (MySQLQueryBuilder) Join(joinType, table, condition string) string {
	return joinType + " JOIN " + table + " ON " + condition
}

func (MySQLQueryBuilder) OrderBy(orderBy ...string) string {
	return "ORDER BY " + strings.Join(orderBy, ", ")
}
