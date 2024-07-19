package query

import (
	"fmt"

	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
)

type Builder struct {
	dbBuilder      db.QueryBuilderStrategy
	cache          *cache.AsyncQueryCache
	query          *structs.Query
	selectValues   []interface{}
	groupByValues  []interface{}
	whereBuilder   *WhereBuilder
	joinBuilder    *JoinBuilder
	orderByBuilder *OrderByBuilder
}

func NewBuilder(dbBuilder db.QueryBuilderStrategy, cache *cache.AsyncQueryCache) *Builder {
	return &Builder{
		dbBuilder: dbBuilder,
		cache:     cache,
		query: &structs.Query{
			Columns:         &[]structs.Column{},
			Conditions:      &[]structs.Where{},
			ConditionGroups: &[]structs.WhereGroup{},
			Joins:           &[]structs.Join{},
			Order:           &[]structs.Order{},
			SubQuery:        &[]structs.Query{},
			Group: &structs.GroupBy{
				Columns: []string{},
				Having:  &[]structs.Having{},
			},
			Limit:  &structs.Limit{},
			Offset: &structs.Offset{},
			Lock:   &structs.Lock{},
		},
		selectValues:   []interface{}{},
		groupByValues:  []interface{}{},
		whereBuilder:   NewWhereBuilder(dbBuilder, cache),
		joinBuilder:    NewJoinBuilder(&[]structs.Join{}),
		orderByBuilder: NewOrderByBuilder(&[]structs.Order{}),
	}
}

func NewBuilderWithQuery(dbBuilder db.QueryBuilderStrategy, query *structs.Query) *Builder {
	return &Builder{
		dbBuilder: dbBuilder,
		query:     query,
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
	b.query.Table = structs.Table{
		Name: table,
	}
	return b
}

func (b *Builder) Select(columns ...string) *Builder {
	for _, column := range columns {
		*b.query.Columns = append(*b.query.Columns, structs.Column{Name: column})
		b.selectValues = append(b.selectValues, column)
	}
	return b
}

func (b *Builder) SelectRaw(raw string, value ...interface{}) *Builder {
	*b.query.Columns = append(*b.query.Columns, structs.Column{Raw: raw, Values: value})
	return b
}

func (b *Builder) Count(columns ...string) *Builder {
	if len(columns) == 0 {
		columns = append(columns, "*")
	}

	for _, column := range columns {
		*b.query.Columns = append(*b.query.Columns, structs.Column{
			Name: column,
			Raw:  fmt.Sprintf("COUNT(%s)", column),
		})
	}
	return b
}

func (b *Builder) aggregate(column string, aggregateFunc string) *Builder {
	*b.query.Columns = append(*b.query.Columns, structs.Column{
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
	*b.query.Group = structs.GroupBy{
		Columns: columns,
		Having:  &[]structs.Having{},
	}
	return b
}

// Having adds a HAVING clause with an AND operator.
func (b *Builder) Having(column string, condition string, value interface{}) *Builder {
	*b.query.Group.Having = append(*b.query.Group.Having, structs.Having{
		Column:    column,
		Condition: condition,
		Value:     value,
		Operator:  consts.LogicalOperator_AND,
	})
	return b
}

// HavingRaw adds a raw HAVING clause with an AND operator.
func (b *Builder) HavingRaw(raw string) *Builder {
	*b.query.Group.Having = append(*b.query.Group.Having, structs.Having{
		Raw:      raw,
		Operator: consts.LogicalOperator_AND,
	})
	return b
}

// OrHaving adds a HAVING clause with an OR operator.
func (b *Builder) OrHaving(column string, condition string, value interface{}) *Builder {
	*b.query.Group.Having = append(*b.query.Group.Having, structs.Having{
		Column:    column,
		Condition: condition,
		Value:     value,
		Operator:  consts.LogicalOperator_OR,
	})
	return b
}

// OrHavingRaw adds a raw HAVING clause with an OR operator.
func (b *Builder) OrHavingRaw(raw string) *Builder {
	*b.query.Group.Having = append(*b.query.Group.Having, structs.Having{
		Raw:      raw,
		Operator: consts.LogicalOperator_OR,
	})
	return b
}

func (b *Builder) Limit(limit int64) *Builder {
	b.query.Limit.Limit = limit
	return b
}

func (b *Builder) Offset(offset int64) *Builder {
	b.query.Offset.Offset = offset
	return b
}

func (b *Builder) SharedLock() *Builder {
	b.query.Lock = &structs.Lock{
		LockType: consts.Lock_SHARE_MODE,
	}
	return b
}

func (b *Builder) LockForUpdate() *Builder {
	b.query.Lock = &structs.Lock{
		LockType: consts.Lock_FOR_UPDATE,
	}
	return b
}

// Build generates the SQL query string and parameter values based on the query builder's current state.
// It returns the generated query string and a slice of parameter values.
func (b *Builder) Build() (string, []interface{}) {
	cacheKey := generateCacheKey(b.query)

	if cachedQuery, found := b.cache.Get(cacheKey); found {
		values := []interface{}{}
		values = append(values, b.selectValues...)
		values = append(values, b.whereBuilder.whereValues...)
		values = append(values, b.groupByValues...)
		return cachedQuery, values
	}

	// preprocess WHERE
	if len(*b.whereBuilder.query.Conditions) > 0 {
		*b.whereBuilder.query.ConditionGroups = append(*b.whereBuilder.query.ConditionGroups, structs.WhereGroup{
			Conditions:   *b.whereBuilder.query.Conditions,
			Operator:     consts.LogicalOperator_AND,
			IsDummyGroup: true,
		})
		b.whereBuilder.query.Conditions = &[]structs.Where{}
	}
	b.query.ConditionGroups = b.whereBuilder.query.ConditionGroups
	b.query.Conditions = b.whereBuilder.query.Conditions

	// preprocess JOIN
	if len(*b.joinBuilder.Joins) > 0 {
		*b.query.Joins = append(*b.query.Joins, *b.joinBuilder.Joins...)
	}

	// preprocess ORDER BY
	if len(*b.orderByBuilder.Order) > 0 {
		*b.query.Order = append(*b.query.Order, *b.orderByBuilder.Order...)
	}

	query, values := b.dbBuilder.Build(b.query)

	b.cache.Set(cacheKey, query)

	return query, values
}

func generateCacheKey(q *structs.Query) string {
	tableKey := q.Table.Name

	columnKey := ""
	for _, c := range *q.Columns {
		columnKey += c.Name + ","
	}

	orderKey := ""
	for _, o := range *q.Order {
		orderKey += o.Column + "," + o.Raw + "," + fmt.Sprint(o.IsAsc)
	}

	joinKey := ""
	for _, j := range *q.Joins {
		joinKey += j.Name + "," + j.SearchColumn + "," + j.SearchCondition + "," + j.SearchTargetColumn + ","
	}

	conditionKey := ""
	for _, c := range *q.ConditionGroups {
		conditionKey += fmt.Sprint(c.Operator) + ","
		for _, w := range c.Conditions {
			conditionKey += w.Column + "," + w.Condition + ","
			if w.Query != nil {
				conditionKey += generateCacheKey(w.Query) + ","
			}
		}
	}

	return fmt.Sprintf("%s|%s|%s|%s|%s",
		tableKey,
		columnKey,
		joinKey,
		conditionKey,
		orderKey,
	)
}

func (b *Builder) GetQuery() *structs.Query {
	b.query.ConditionGroups = b.whereBuilder.query.ConditionGroups
	b.query.Conditions = b.whereBuilder.query.Conditions

	return b.query
}
