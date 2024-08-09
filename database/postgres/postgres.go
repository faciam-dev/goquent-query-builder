package postgres

import (
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db/base"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

type PostgreSQLQueryBuilder struct {
	base.BaseQueryBuilder
	base.DeleteBaseBuilder
	base.InsertBaseBuilder
	base.UpdateBaseBuilder

	/*
		selectBuilderStrategy base.SelectBuilderStrategy
		FromBuilderStrategy   base.FromBuilderStrategy
		joinBuilderStrategy   base.JoinBuilderStrategy
	*/
	//whereBuilderStrategy   base.WhereBuilderStrategy
	/*
		orderByBuilderStrategy base.OrderByBuilderStrategy
		groupByBuilderStrategy base.GroupByBuilderStrategy
		limitBuilderStrategy   base.LimitBuilderStrategy
		OffsetBuilderStrategy  base.OffsetBuilderStrategy
		insertBuilderStrategy  base.InsertBuilderStrategy
		updateBuilderStrategy  base.UpdateBuilderStrategy
		deleteBuilderStrategy  base.DeleteBuilderStrategy
	*/

	WherePostgreSQLBuilder

	util interfaces.SQLUtils
}

func NewPostgreSQLQueryBuilder() *PostgreSQLQueryBuilder {
	queryBuilder := &PostgreSQLQueryBuilder{}
	u := NewSQLUtils()
	queryBuilder.util = u
	//queryBuilder.BaseQueryBuilder = *base.NewBaseQueryBuilder()
	queryBuilder.SelectBaseBuilder = *base.NewSelectBaseBuilder(u, &[]string{})
	queryBuilder.JoinBaseBuilder = *base.NewJoinBaseBuilder(u, &structs.Joins{})
	queryBuilder.FromBaseBuilder = *base.NewFromBaseBuilder(u)
	queryBuilder.GroupByBaseBuilder = *base.NewGroupByBaseBuilder(u)
	queryBuilder.OrderByBaseBuilder = *base.NewOrderByBaseBuilder(u, &[]structs.Order{})
	queryBuilder.DeleteBaseBuilder = *base.NewDeleteBaseBuilder(u, &structs.DeleteQuery{})
	queryBuilder.InsertBaseBuilder = *base.NewInsertBaseBuilder(u, &structs.InsertQuery{})
	queryBuilder.UpdateBaseBuilder = *base.NewUpdateBaseBuilder(u, &structs.UpdateQuery{})
	queryBuilder.WherePostgreSQLBuilder = *NewWherePostgreSQLBuilder(u, []structs.WhereGroup{})
	/*
		queryBuilder.selectBuilderStrategy = base.NewSelectBaseBuilder(u, &[]string{})
		queryBuilder.FromBuilderStrategy = base.NewFromBaseBuilder(u)
		queryBuilder.joinBuilderStrategy = base.NewJoinBaseBuilder(u, &structs.Joins{})
	*/
	//queryBuilder.whereBuilderStrategy = NewWherePostgreSQLBuilder(u, []structs.WhereGroup{})
	/*
		queryBuilder.orderByBuilderStrategy = base.NewOrderByBaseBuilder(u, &[]structs.Order{})
		queryBuilder.groupByBuilderStrategy = base.NewGroupByBaseBuilder(u)
		queryBuilder.limitBuilderStrategy = base.NewLimitBaseBuilder()
		queryBuilder.OffsetBuilderStrategy = base.NewOffsetBaseBuilder()
		queryBuilder.insertBuilderStrategy = base.NewInsertBaseBuilder(u, &structs.InsertQuery{})
		queryBuilder.updateBuilderStrategy = base.NewUpdateBaseBuilder(u, &structs.UpdateQuery{})
		queryBuilder.deleteBuilderStrategy = base.NewDeleteBaseBuilder(u, &structs.DeleteQuery{})
	*/
	return queryBuilder
}

// Build builds the query.
func (m PostgreSQLQueryBuilder) Build(sb *[]byte, q *structs.Query, number int, unions *[]structs.Union) []interface{} {
	// SELECT
	*sb = append(*sb, "SELECT "...)
	colValues := m.Select(sb, q.Columns, q.Table.Name, q.Joins)

	*sb = append(*sb, " "...)
	m.From(sb, q.Table.Name)
	values := colValues

	// JOIN
	joinValues := m.Join(sb, q.Joins)
	values = append(values, joinValues...)

	// WHERE
	whereValues := m.Where(sb, q.ConditionGroups)
	values = append(values, whereValues...)

	// GROUP BY / HAVING
	groupByValues := m.GroupBy(sb, q.Group)
	values = append(values, groupByValues...)

	// ORDER BY
	m.OrderBy(sb, q.Order)

	// LIMIT
	m.Limit(sb, q.Limit)

	// OFFSET
	m.Offset(sb, q.Offset)

	// LOCK
	m.Lock(sb, q.Lock)

	// UNION
	m.Union(sb, unions, number)

	//query := sb.String()
	//sb.Reset()

	return values
}

func (m PostgreSQLQueryBuilder) Where(sb *[]byte, conditionGroups []structs.WhereGroup) []interface{} {
	return m.WherePostgreSQLBuilder.Where(sb, conditionGroups)
}
