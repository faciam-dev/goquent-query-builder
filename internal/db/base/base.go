package base

import (
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type SelectBuilderStrategy interface {
	Select(sb *strings.Builder, columns *[]structs.Column, tableName string, joins *structs.Joins) []interface{}
}

type FromBuilderStrategy interface {
	From(sb *strings.Builder, tableName string)
}

type WhereBuilderStrategy interface {
	Where(sb *strings.Builder, wg *[]structs.WhereGroup) []interface{}
	ProcessFullText(sb *strings.Builder, c structs.Where) []interface{}
}

type OrderByBuilderStrategy interface {
	OrderBy(sb *strings.Builder, order *[]structs.Order)
}

type GroupByBuilderStrategy interface {
	GroupBy(sb *strings.Builder, groupBy *structs.GroupBy) []interface{}
}

type JoinBuilderStrategy interface {
	Join(sb *strings.Builder, joins *structs.Joins) []interface{}
}

type LimitBuilderStrategy interface {
	Limit(sb *strings.Builder, limit *structs.Limit)
}

type OffsetBuilderStrategy interface {
	Offset(sb *strings.Builder, offset *structs.Offset)
}

type InsertBuilderStrategy interface {
	Insert(q *structs.InsertQuery) (string, []interface{})
	BuildInsert(q *structs.InsertQuery) (string, []interface{})
}

type UpdateBuilderStrategy interface {
	Update(q *structs.UpdateQuery) *UpdateBaseBuilder
	BuildUpdate(q *structs.UpdateQuery) (string, []interface{})
}

type DeleteBuilderStrategy interface {
	Delete(q *structs.DeleteQuery) *DeleteBaseBuilder
	BuildDelete(q *structs.DeleteQuery) (string, []interface{})
}

type BaseQueryBuilder struct {
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
}

// Lock returns the lock statement.
func (BaseQueryBuilder) Lock(lock *structs.Lock) string {
	if lock == nil || lock.LockType == "" {
		return ""
	}

	return " " + lock.LockType
}

// Build builds the query.
func (m BaseQueryBuilder) Build(cacheKey string, q *structs.Query) (string, []interface{}) {
	sb := &strings.Builder{}

	// grow the string builder based on the length of the cache key
	if len(cacheKey) < consts.StringBuffer_Short_Query_Grow {
		sb.Grow(consts.StringBuffer_Short_Query_Grow)
	} else if len(cacheKey) < consts.StringBuffer_Middle_Query_Grow {
		sb.Grow(consts.StringBuffer_Middle_Query_Grow)
	} else {
		sb.Grow(consts.StringBuffer_Long_Query_Grow)
	}

	// SELECT
	sb.WriteString("SELECT ")
	colValues := m.Select(sb, q.Columns, q.Table.Name, q.Joins)

	// FROM
	sb.WriteString(" ")
	m.From(sb, q.Table.Name)
	values := colValues

	// JOIN
	joinValues := m.Join(sb, q.Joins)
	values = append(values, joinValues...)

	// WHERE
	whereValues := m.Where(sb, q.ConditionGroups)
	values = append(values, whereValues...)

	// ORDER BY
	m.OrderBy(sb, q.Order)

	// GROUP BY / HAVING
	groupByValues := m.GroupBy(sb, q.Group)
	values = append(values, groupByValues...)

	// LIMIT
	m.Limit(sb, q.Limit)

	// OFFSET
	m.Offset(sb, q.Offset)

	// LOCK
	lock := m.Lock(q.Lock)
	sb.WriteString(lock)

	query := sb.String()
	sb.Reset()

	return query, values
}
