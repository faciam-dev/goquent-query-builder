package mysql

import (
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db/base"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

type MySQLQueryBuilder struct {
	base.BaseQueryBuilder
	base.DeleteBaseBuilder
	base.InsertBaseBuilder
	base.UpdateBaseBuilder

	WhereMySQLBuilder

	util interfaces.SQLUtils
}

func NewMySQLQueryBuilder() *MySQLQueryBuilder {
	queryBuilder := &MySQLQueryBuilder{}
	u := NewSQLUtils()
	queryBuilder.util = u
	queryBuilder.SelectBaseBuilder = *base.NewSelectBaseBuilder(u, &[]string{})
	queryBuilder.JoinBaseBuilder = *base.NewJoinBaseBuilder(u, &structs.Joins{})
	queryBuilder.FromBaseBuilder = *base.NewFromBaseBuilder(u)
	queryBuilder.GroupByBaseBuilder = *base.NewGroupByBaseBuilder(u)
	queryBuilder.OrderByBaseBuilder = *base.NewOrderByBaseBuilder(u, &[]structs.Order{})
	queryBuilder.WhereMySQLBuilder = *NewWhereMySQLBuilder(u, []structs.WhereGroup{})
	queryBuilder.UpdateBaseBuilder = *base.NewUpdateBaseBuilder(u, &structs.UpdateQuery{})
	queryBuilder.InsertBaseBuilder = *base.NewInsertBaseBuilder(u, &structs.InsertQuery{})
	queryBuilder.DeleteBaseBuilder = *base.NewDeleteBaseBuilder(u, &structs.DeleteQuery{})
	return queryBuilder
}

// Build builds the query.
func (m MySQLQueryBuilder) Build(sb *[]byte, q *structs.Query, number int, unions *[]structs.Union) ([]interface{}, error) {
	// SELECT
	*sb = append(*sb, "SELECT "...)
	colValues := m.Select(sb, q.Columns, q.Table.Name, q.Joins)

	*sb = append(*sb, " "...)
	m.From(sb, q.Table.Name)
	values := colValues

	// JOIN
	if q.Joins.JoinClauses != nil && (len(*q.Joins.JoinClauses) > 0 || len(*q.Joins.LateralJoins) > 0 || len(*q.Joins.Joins) > 0) {
		joinValues := m.Join(sb, q.Joins)
		values = append(values, joinValues...)
	}

	// WHERE
	if len(q.ConditionGroups) > 0 {
		whereValues, err := m.Where(sb, q.ConditionGroups)
		if err != nil {
			return nil, err
		}
		values = append(values, whereValues...)
	}

	// GROUP BY / HAVING
	if q.Group != nil && len(q.Group.Columns) > 0 {
		groupByValues := m.GroupBy(sb, q.Group)
		values = append(values, groupByValues...)
	}

	// ORDER BY
	if len(*q.Order) > 0 {
		m.OrderBy(sb, q.Order)
	}

	// LIMIT
	if q.Limit.Limit > 0 {
		m.Limit(sb, q.Limit)
	}

	// OFFSET
	if q.Offset.Offset > 0 {
		m.Offset(sb, q.Offset)
	}

	// LOCK
	if q.Lock != nil && q.Lock.LockType != "" {
		m.Lock(sb, q.Lock)
	}

	// UNION
	if unions != nil && len(*unions) > 0 {
		m.Union(sb, unions, number)
	}

	return values, nil
}

func (m MySQLQueryBuilder) Where(sb *[]byte, c []structs.WhereGroup) ([]interface{}, error) {
	return m.WhereMySQLBuilder.Where(sb, c)
}
