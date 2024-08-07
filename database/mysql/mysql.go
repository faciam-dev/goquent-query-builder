package mysql

import (
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db/base"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

type MySQLQueryBuilder struct {
	base.BaseQueryBuilder
	base.DeleteBaseBuilder
	base.InsertBaseBuilder
	base.UpdateBaseBuilder

	selectBuilderStrategy  base.SelectBuilderStrategy
	FromBuilderStrategy    base.FromBuilderStrategy
	joinBuilderStrategy    base.JoinBuilderStrategy
	whereBuilderStrategy   base.WhereBuilderStrategy
	orderByBuilderStrategy base.OrderByBuilderStrategy
	groupByBuilderStrategy base.GroupByBuilderStrategy
	limitBuilderStrategy   base.LimitBuilderStrategy
	OffsetBuilderStrategy  base.OffsetBuilderStrategy
	insertBuilderStrategy  base.InsertBuilderStrategy
	updateBuilderStrategy  base.UpdateBuilderStrategy
	deleteBuilderStrategy  base.DeleteBuilderStrategy

	util interfaces.SQLUtils
}

func NewMySQLQueryBuilder() *MySQLQueryBuilder {
	queryBuilder := &MySQLQueryBuilder{}
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
	queryBuilder.selectBuilderStrategy = base.NewSelectBaseBuilder(u, &[]string{})
	queryBuilder.FromBuilderStrategy = base.NewFromBaseBuilder(u)
	queryBuilder.joinBuilderStrategy = base.NewJoinBaseBuilder(u, &structs.Joins{})
	queryBuilder.whereBuilderStrategy = NewWhereMySQLBuilder(u, []structs.WhereGroup{})
	queryBuilder.orderByBuilderStrategy = base.NewOrderByBaseBuilder(u, &[]structs.Order{})
	queryBuilder.groupByBuilderStrategy = base.NewGroupByBaseBuilder(u)
	queryBuilder.limitBuilderStrategy = base.NewLimitBaseBuilder()
	queryBuilder.OffsetBuilderStrategy = base.NewOffsetBaseBuilder()
	queryBuilder.insertBuilderStrategy = base.NewInsertBaseBuilder(u, &structs.InsertQuery{})
	queryBuilder.updateBuilderStrategy = base.NewUpdateBaseBuilder(u, &structs.UpdateQuery{})
	queryBuilder.deleteBuilderStrategy = base.NewDeleteBaseBuilder(u, &structs.DeleteQuery{})
	return queryBuilder
}

// Build builds the query.
func (m MySQLQueryBuilder) Build(sb *strings.Builder, cacheKey string, q *structs.Query, number int, unions *[]structs.Union) (string, []interface{}) {
	// SELECT
	sb.WriteString("SELECT ")
	colValues := m.selectBuilderStrategy.Select(sb, q.Columns, q.Table.Name, q.Joins)

	sb.WriteString(" ")
	m.FromBuilderStrategy.From(sb, q.Table.Name)
	values := colValues

	// JOIN
	if q.Joins.JoinClauses != nil && (len(*q.Joins.JoinClauses) > 0 || len(*q.Joins.LateralJoins) > 0 || len(*q.Joins.Joins) > 0) {
		joinValues := m.joinBuilderStrategy.Join(sb, q.Joins)
		values = append(values, joinValues...)
	}

	// WHERE
	if len(q.ConditionGroups) > 0 {
		whereValues := m.whereBuilderStrategy.Where(sb, q.ConditionGroups)
		values = append(values, whereValues...)
	}

	// GROUP BY / HAVING
	if q.Group != nil && len(q.Group.Columns) > 0 {
		groupByValues := m.groupByBuilderStrategy.GroupBy(sb, q.Group)
		values = append(values, groupByValues...)
	}

	// ORDER BY
	if len(*q.Order) > 0 {
		m.orderByBuilderStrategy.OrderBy(sb, q.Order)
	}

	// LIMIT
	if q.Limit != nil && q.Limit.Limit > 0 {
		m.limitBuilderStrategy.Limit(sb, q.Limit)
	}

	// OFFSET
	if q.Offset != nil && q.Offset.Offset > 0 {
		m.OffsetBuilderStrategy.Offset(sb, q.Offset)
	}

	// LOCK
	if q.Lock != nil && q.Lock.LockType != "" {
		m.Lock(sb, q.Lock)
	}

	// UNION
	m.Union(sb, unions, number)

	//query := sb.String()
	//sb.Reset()

	query := ""

	return query, values
}

func (m MySQLQueryBuilder) Where(sb *strings.Builder, c []structs.WhereGroup) []interface{} {
	return m.whereBuilderStrategy.Where(sb, c)
}
