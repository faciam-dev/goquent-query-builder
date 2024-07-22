package db

type MySQLQueryBuilder struct {
	BaseQueryBuilder
}

func NewMySQLQueryBuilder() *MySQLQueryBuilder {
	queryBuilder := &MySQLQueryBuilder{}
	queryBuilder.columnNames = &[]string{}
	return queryBuilder
}
