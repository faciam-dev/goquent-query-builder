package base

import (
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

type SelectBuilderStrategy interface {
	Select(sb *[]byte, columns *[]structs.Column, tableName string, joins *structs.Joins) []interface{}
}

type FromBuilderStrategy interface {
	From(sb *[]byte, tableName string)
}

type WhereBuilderStrategy interface {
	Where(sb *[]byte, wg []structs.WhereGroup) []interface{}
	ProcessFullText(sb *[]byte, c structs.Where) []interface{}
}

type OrderByBuilderStrategy interface {
	OrderBy(sb *[]byte, order *[]structs.Order)
}

type GroupByBuilderStrategy interface {
	GroupBy(sb *[]byte, groupBy *structs.GroupBy) []interface{}
}

type JoinBuilderStrategy interface {
	Join(sb *[]byte, joins *structs.Joins) []interface{}
}

type UnionBuilderStrategy interface {
	Union(sb *[]byte, union *structs.Union, number int) []interface{}
}

type LimitBuilderStrategy interface {
	Limit(sb *[]byte, limit structs.Limit)
}

type OffsetBuilderStrategy interface {
	Offset(sb *[]byte, offset structs.Offset)
}

type InsertBuilderStrategy interface {
	Insert(q *structs.InsertQuery) (string, []interface{}, error)
	BuildInsert(q *structs.InsertQuery) (string, []interface{}, error)
}

type UpdateBuilderStrategy interface {
	Update(q *structs.UpdateQuery) *UpdateBaseBuilder
	BuildUpdate(q *structs.UpdateQuery) (string, []interface{}, error)
}

type DeleteBuilderStrategy interface {
	Delete(q *structs.DeleteQuery) *DeleteBaseBuilder
	BuildDelete(q *structs.DeleteQuery) (string, []interface{}, error)
}

type BaseQueryBuilder struct {
	UnionBaseBuilder
	SelectBaseBuilder
	FromBaseBuilder
	WhereBaseBuilder
	JoinBaseBuilder
	OrderByBaseBuilder
	GroupByBaseBuilder
	LimitBaseBuilder
	OffsetBaseBuilder
	InsertBaseBuilder
	UpdateBaseBuilder
	DeleteBaseBuilder

	util interfaces.SQLUtils
}

func NewBaseQueryBuilder() *BaseQueryBuilder {
	queryBuilder := &BaseQueryBuilder{}
	u := NewSQLUtils()
	queryBuilder.util = u
	queryBuilder.SelectBaseBuilder = *NewSelectBaseBuilder(u, &[]string{})
	queryBuilder.FromBaseBuilder = *NewFromBaseBuilder(u)
	queryBuilder.JoinBaseBuilder = *NewJoinBaseBuilder(u, &structs.Joins{})
	queryBuilder.WhereBaseBuilder = *NewWhereBaseBuilder(u, []structs.WhereGroup{})
	queryBuilder.OrderByBaseBuilder = *NewOrderByBaseBuilder(u, &[]structs.Order{})
	queryBuilder.GroupByBaseBuilder = *NewGroupByBaseBuilder(u)
	queryBuilder.LimitBaseBuilder = *NewLimitBaseBuilder()
	queryBuilder.OffsetBaseBuilder = *NewOffsetBaseBuilder()
	queryBuilder.InsertBaseBuilder = *NewInsertBaseBuilder(u, &structs.InsertQuery{})
	queryBuilder.UpdateBaseBuilder = *NewUpdateBaseBuilder(u, &structs.UpdateQuery{})
	queryBuilder.DeleteBaseBuilder = *NewDeleteBaseBuilder(u, &structs.DeleteQuery{})
	return queryBuilder
}

// Lock returns the lock statement.
func (BaseQueryBuilder) Lock(sb *[]byte, lock *structs.Lock) {
	if lock == nil || lock.LockType == "" {
		return
	}

	*sb = append(*sb, " "...)
	*sb = append(*sb, lock.LockType...)
}

// Build builds the query.
func (m BaseQueryBuilder) Build(sb *[]byte, q *structs.Query, number int, unions *[]structs.Union) []interface{} {
	values := make([]interface{}, 0)

	// SELECT
	*sb = append(*sb, "SELECT "...)
	colValues := m.Select(sb, q.Columns, q.Table.Name, q.Joins)

	// FROM
	*sb = append(*sb, " "...)
	m.From(sb, q.Table.Name)
	values = append(values, colValues...)

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
