package query

import (
	"fmt"
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
)

type Builder struct {
	dbBuilder      db.QueryBuilderStrategy
	cache          cache.Cache
	query          *structs.Query
	selectQuery    *structs.SelectQuery
	selectValues   []interface{}
	groupByValues  []interface{}
	whereBuilder   *WhereBuilder
	joinBuilder    *JoinBuilder
	orderByBuilder *OrderByBuilder
}

func NewBuilder(dbBuilder db.QueryBuilderStrategy, cache cache.Cache) *Builder {
	return &Builder{
		dbBuilder: dbBuilder,
		cache:     cache,
		query: &structs.Query{
			Table:           structs.Table{},
			Columns:         &[]structs.Column{},
			ConditionGroups: &[]structs.WhereGroup{},
			Joins:           &structs.Joins{},
			Order:           &[]structs.Order{},
			Group:           &structs.GroupBy{},
			Limit:           &structs.Limit{},
			Offset:          &structs.Offset{},
			Lock:            &structs.Lock{},
			SubQuery:        &[]structs.Query{},
		},
		selectQuery: &structs.SelectQuery{
			Table:    "",
			Columns:  &[]structs.Column{},
			Limit:    &structs.Limit{},
			SubQuery: &[]structs.Query{},
			Group:    &structs.GroupBy{},
			Offset:   &structs.Offset{},
			Lock:     &structs.Lock{},
		},
		selectValues:   []interface{}{},
		groupByValues:  []interface{}{},
		whereBuilder:   NewWhereBuilder(dbBuilder, cache),
		joinBuilder:    NewJoinBuilder(dbBuilder, cache),
		orderByBuilder: NewOrderByBuilder(&[]structs.Order{}),
	}
}

func (b *Builder) SetWhereBuilder(whereBuilder *WhereBuilder) {
	b.whereBuilder = whereBuilder
}

func (b *Builder) SetJoinBuilder(joinBuilder *JoinBuilder) {
	b.joinBuilder = joinBuilder
}

func (b *Builder) SetOrderByBuilder(orderByBuilder *OrderByBuilder) {
	b.orderByBuilder = orderByBuilder
}

func (b *Builder) Table(table string) *Builder {
	b.selectQuery.Table = table
	return b
}

func (b *Builder) Select(columns ...string) *Builder {
	for _, column := range columns {
		*b.selectQuery.Columns = append(*b.selectQuery.Columns, structs.Column{Name: column})
	}
	return b
}

func (b *Builder) SelectRaw(raw string, value ...interface{}) *Builder {
	b.selectValues = append(b.selectValues, value...)
	*b.selectQuery.Columns = append(*b.selectQuery.Columns, structs.Column{Raw: raw, Values: value})
	return b
}

func (b *Builder) Count(columns ...string) *Builder {
	if len(columns) == 0 {
		columns = append(columns, "*")
	}

	for _, column := range columns {
		*b.selectQuery.Columns = append(*b.selectQuery.Columns, structs.Column{
			Name: column,
			Raw:  fmt.Sprintf("COUNT(%s)", column),
		})
	}
	return b
}

func (b *Builder) aggregate(column string, aggregateFunc string) *Builder {
	*b.selectQuery.Columns = append(*b.selectQuery.Columns, structs.Column{
		Name: column,
		Raw:  fmt.Sprintf("%s(%s)", aggregateFunc, column),
	})
	return b
}

func (b *Builder) Max(column string) *Builder {
	return b.aggregate(column, "MAX")
}

func (b *Builder) Min(column string) *Builder {
	return b.aggregate(column, "MIN")
}

func (b *Builder) Sum(column string) *Builder {
	return b.aggregate(column, "SUM")
}

func (b *Builder) Avg(column string) *Builder {
	return b.aggregate(column, "AVG")
}

func (b *Builder) Where(column string, condition string, value ...interface{}) *Builder {
	b.whereBuilder.Where(column, condition, value...)
	return b
}

func (b *Builder) OrWhere(column string, condition string, value ...interface{}) *Builder {
	b.whereBuilder.OrWhere(column, condition, value...)
	return b
}

func (b *Builder) WhereRaw(raw string, value ...interface{}) *Builder {
	b.whereBuilder.WhereRaw(raw, value...)
	return b
}

func (b *Builder) OrWhereRaw(raw string, value ...interface{}) *Builder {
	b.whereBuilder.OrWhereRaw(raw, value...)
	return b
}

func (b *Builder) WhereQuery(column string, condition string, q *Builder) *Builder {
	b.whereBuilder.WhereQuery(column, condition, q)

	return b
}

func (b *Builder) OrWhereQuery(column string, condition string, q *Builder) *Builder {
	b.whereBuilder.OrWhereQuery(column, condition, q)

	return b
}

func (b *Builder) WhereGroup(fn func(b *WhereBuilder) *WhereBuilder) *Builder {
	b.whereBuilder.WhereGroup(fn)

	return b
}

func (b *Builder) OrWhereGroup(fn func(b *WhereBuilder) *WhereBuilder) *Builder {
	b.whereBuilder.OrWhereGroup(fn)

	return b
}

// Join adds a JOIN clause.
func (b *Builder) Join(table string, my string, condition string, target string) *Builder {
	b.joinBuilder.Join(table, my, condition, target)

	return b
}

// LeftJoin adds a LEFT JOIN clause.
func (b *Builder) LeftJoin(table string, my string, condition string, target string) *Builder {
	b.joinBuilder.LeftJoin(table, my, condition, target)

	return b
}

// RightJoin adds a RIGHT JOIN clause.
func (b *Builder) RightJoin(table string, my string, condition string, target string) *Builder {
	b.joinBuilder.RightJoin(table, my, condition, target)

	return b
}

// CrossJoin adds a CROSS JOIN clause.
func (b *Builder) CrossJoin(table string) *Builder {
	b.joinBuilder.CrossJoin(table)

	return b
}

func (b *Builder) JoinQuery(table string, fn func(b *JoinClauseBuilder) *JoinClauseBuilder) *Builder {
	b.joinBuilder.JoinQuery(table, fn)

	return b
}

func (b *Builder) LeftJoinQuery(table string, fn func(b *JoinClauseBuilder) *JoinClauseBuilder) *Builder {
	b.joinBuilder.LeftJoinQuery(table, fn)

	return b
}

func (b *Builder) RightJoinQuery(table string, fn func(b *JoinClauseBuilder) *JoinClauseBuilder) *Builder {
	b.joinBuilder.RightJoinQuery(table, fn)

	return b
}

func (b *Builder) JoinSub(q *Builder, alias, my, condition, target string) *Builder {
	b.joinBuilder.JoinSub(q, alias, my, condition, target)

	return b
}

func (b *Builder) LeftJoinSub(q *Builder, alias, my, condition, target string) *Builder {
	b.joinBuilder.LeftJoinSub(q, alias, my, condition, target)

	return b
}

func (b *Builder) RightJoinSub(q *Builder, alias, my, condition, target string) *Builder {
	b.joinBuilder.RightJoinSub(q, alias, my, condition, target)

	return b
}

// OrderBy adds an ORDER BY clause.
func (b *Builder) OrderBy(column string, ascDesc string) *Builder {
	b.orderByBuilder.OrderBy(column, ascDesc)

	return b
}

// ReOrder removes all ORDER BY clauses.
func (b *Builder) ReOrder() *Builder {
	b.orderByBuilder.ReOrder()
	return b
}

// OrderByRaw adds a raw ORDER BY clause.
func (b *Builder) OrderByRaw(raw string) *Builder {
	b.orderByBuilder.OrderByRaw(raw)
	return b
}

// GroupBy adds a GROUP BY clause.
func (b *Builder) GroupBy(columns ...string) *Builder {
	*b.selectQuery.Group = structs.GroupBy{
		Columns: columns,
		Having:  &[]structs.Having{},
	}
	return b
}

// Having adds a HAVING clause with an AND operator.
func (b *Builder) Having(column string, condition string, value interface{}) *Builder {
	*b.selectQuery.Group.Having = append(*b.selectQuery.Group.Having, structs.Having{
		Column:    column,
		Condition: condition,
		Value:     value,
		Operator:  consts.LogicalOperator_AND,
	})
	return b
}

// HavingRaw adds a raw HAVING clause with an AND operator.
func (b *Builder) HavingRaw(raw string) *Builder {
	*b.selectQuery.Group.Having = append(*b.selectQuery.Group.Having, structs.Having{
		Raw:      raw,
		Operator: consts.LogicalOperator_AND,
	})
	return b
}

// OrHaving adds a HAVING clause with an OR operator.
func (b *Builder) OrHaving(column string, condition string, value interface{}) *Builder {
	*b.selectQuery.Group.Having = append(*b.selectQuery.Group.Having, structs.Having{
		Column:    column,
		Condition: condition,
		Value:     value,
		Operator:  consts.LogicalOperator_OR,
	})
	return b
}

// OrHavingRaw adds a raw HAVING clause with an OR operator.
func (b *Builder) OrHavingRaw(raw string) *Builder {
	*b.selectQuery.Group.Having = append(*b.selectQuery.Group.Having, structs.Having{
		Raw:      raw,
		Operator: consts.LogicalOperator_OR,
	})
	return b
}

func (b *Builder) Limit(limit int64) *Builder {
	b.selectQuery.Limit.Limit = limit
	return b
}

func (b *Builder) Offset(offset int64) *Builder {
	b.selectQuery.Offset.Offset = offset
	return b
}

func (b *Builder) SharedLock() *Builder {
	b.selectQuery.Lock = &structs.Lock{
		LockType: consts.Lock_SHARE_MODE,
	}
	return b
}

func (b *Builder) LockForUpdate() *Builder {
	b.selectQuery.Lock = &structs.Lock{
		LockType: consts.Lock_FOR_UPDATE,
	}
	return b
}

// Build generates the SQL query string and parameter values based on the query builder's current state.
// It returns the generated query string and a slice of parameter values.
func (b *Builder) Build() (string, []interface{}) {
	b.buildQuery()

	cacheKey := generateCacheKey(b.query)

	if cachedQuery, found := b.cache.Get(cacheKey); found {
		values := make([]interface{}, 0, len(b.selectValues)+len(b.joinBuilder.joinValues)+len(b.whereBuilder.whereValues)+len(b.groupByValues))
		values = append(values, b.selectValues...)
		values = append(values, b.joinBuilder.joinValues...)
		values = append(values, b.whereBuilder.whereValues...)
		values = append(values, b.groupByValues...)
		return cachedQuery, values
	}

	query, values := b.dbBuilder.Build(cacheKey, b.query)

	b.cache.Set(cacheKey, query)

	return query, values
}

func (b *Builder) buildQuery() {
	// preprocess WHERE
	wg := b.whereBuilder.query.ConditionGroups
	if len(*b.whereBuilder.query.Conditions) > 0 {
		*wg = append(*wg, structs.WhereGroup{
			Conditions:   *b.whereBuilder.query.Conditions,
			Operator:     consts.LogicalOperator_AND,
			IsDummyGroup: true,
		})
		b.whereBuilder.query.Conditions = &[]structs.Where{}
	}

	// preprocess ORDER BY
	o := b.orderByBuilder.Order

	b.query.Table = structs.Table{
		Name: b.selectQuery.Table,
	}
	b.query.Columns = b.selectQuery.Columns
	b.query.ConditionGroups = wg
	b.query.Joins = b.joinBuilder.Joins
	b.query.Order = o
	b.query.Group = b.selectQuery.Group
	b.query.Limit = b.selectQuery.Limit
	b.query.Offset = b.selectQuery.Offset
	b.query.Lock = b.selectQuery.Lock
	b.query.SubQuery = b.selectQuery.SubQuery

}

func generateCacheKey(q *structs.Query) string {
	var sb strings.Builder
	sb.Grow(consts.StringBuffer_CacheKey_Grow)

	// Table key
	sb.WriteString(q.Table.Name)

	// Column key
	for _, c := range *q.Columns {
		sb.WriteString(c.Name)
		sb.WriteString(",")
	}
	sb.WriteString("|")

	// Order key
	for _, o := range *q.Order {
		sb.WriteString(o.Column)
		sb.WriteString(",")
		sb.WriteString(o.Raw)
		sb.WriteString(",")
		if o.IsAsc {
			sb.WriteString("ASC")
		} else {
			sb.WriteString("DESC")
		}
	}
	sb.WriteString("|")

	// Join key
	if q.Joins.Joins != nil {
		for _, j := range *q.Joins.Joins {
			sb.WriteString(j.Name)
			sb.WriteString(",")
			sb.WriteString(j.SearchColumn)
			sb.WriteString(",")
			sb.WriteString(j.SearchCondition)
			sb.WriteString(",")
			sb.WriteString(j.SearchTargetColumn)
			sb.WriteString(",")
		}
	}
	if q.Joins.JoinClause != nil {
		for _, o := range *q.Joins.JoinClause.On {
			sb.WriteString(o.Column)
			sb.WriteString(",")
			sb.WriteString(o.Condition)
			sb.WriteString(",")
			if o.Operator == consts.LogicalOperator_OR {
				sb.WriteString("OR")
			} else {
				sb.WriteString("AND")
			}
			sb.WriteString(",")
		}
	}
	sb.WriteString("|")

	// Condition key
	for _, c := range *q.ConditionGroups {
		if c.Operator == consts.LogicalOperator_AND {
			sb.WriteString("AND")
		} else {
			sb.WriteString("OR")
		}
		sb.WriteString(",")
		for _, w := range c.Conditions {
			if w.Operator == consts.LogicalOperator_AND {
				sb.WriteString("AND")
			} else {
				sb.WriteString("OR")
			}
			sb.WriteString(w.Column)
			sb.WriteString(",")
			sb.WriteString(w.Condition)
			sb.WriteString(",")
			if w.Query != nil {
				sb.WriteString(generateCacheKey(w.Query))
				sb.WriteString(",")
			}
		}
	}

	key := sb.String()

	return key
}

func (b *Builder) GetQuery() *structs.Query {
	return b.query
}
