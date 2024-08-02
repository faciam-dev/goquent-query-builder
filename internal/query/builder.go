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
	dbBuilder     db.QueryBuilderStrategy
	cache         cache.Cache
	query         *structs.Query
	selectQuery   *structs.SelectQuery
	selectValues  []interface{}
	groupByValues []interface{}
	WhereBuilder[Builder]
	JoinBuilder[Builder]
	orderByBuilder *OrderByBuilder
	BaseBuilder
}

func NewBuilder(dbBuilder db.QueryBuilderStrategy, cache cache.Cache) *Builder {
	b := &Builder{
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
			Union:    &[]structs.Union{},
		},
		selectValues:  []interface{}{},
		groupByValues: []interface{}{},
		//joinBuilder:    NewJoinBuilder(dbBuilder, cache),
		orderByBuilder: NewOrderByBuilder(&[]structs.Order{}),
	}

	whereBuilder := NewWhereBuilder[Builder](dbBuilder, cache)
	whereBuilder.SetParent(b)
	b.WhereBuilder = *whereBuilder

	joinBuilder := NewJoinBuilder[Builder](dbBuilder, cache)
	joinBuilder.SetParent(b)
	b.JoinBuilder = *joinBuilder

	return b
}

func (b *Builder) SetWhereBuilder(whereBuilder *WhereBuilder[Builder]) {
	b.WhereBuilder = *whereBuilder
}

func (b *Builder) SetJoinBuilder(joinBuilder *JoinBuilder[Builder]) {
	b.JoinBuilder = *joinBuilder
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

// Count adds a COUNT aggregate function to the query.
func (b *Builder) Count(columns ...string) *Builder {
	if len(columns) == 0 {
		columns = append(columns, "*")
	}

	for i, c := range *b.selectQuery.Columns {
		for _, col := range columns {
			if c.Name == col {
				(*b.selectQuery.Columns)[i].Count = true
			}
		}
	}

out:
	for _, column := range columns {
		for _, c := range *b.selectQuery.Columns {
			if c.Count {
				continue out
			}
		}

		*b.selectQuery.Columns = append(*b.selectQuery.Columns, structs.Column{
			Name:  column,
			Count: true,
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

// Max adds a MAX aggregate function to the query.
func (b *Builder) Max(column string) *Builder {
	return b.aggregate(column, "MAX")
}

// Min adds a MIN aggregate function to the query.
func (b *Builder) Min(column string) *Builder {
	return b.aggregate(column, "MIN")
}

// Sum adds a SUM aggregate function to the query.
func (b *Builder) Sum(column string) *Builder {
	return b.aggregate(column, "SUM")
}

// Avg adds an AVG aggregate function to the query.
func (b *Builder) Avg(column string) *Builder {
	return b.aggregate(column, "AVG")
}

func (b *Builder) Distinct(column ...string) *Builder {
	for i, c := range *b.selectQuery.Columns {
		for _, col := range column {
			if c.Name == col {
				(*b.selectQuery.Columns)[i].Distinct = true
			}
		}
	}

out:
	for _, c := range column {
		for _, c := range *b.selectQuery.Columns {
			if c.Count {
				continue out
			}
		}
		*b.selectQuery.Columns = append(*b.selectQuery.Columns, structs.Column{
			Name:     c,
			Distinct: true,
		})
	}

	return b
}

func (b *Builder) Union(sb *Builder) *Builder {
	*b.selectQuery.Union = append(*b.selectQuery.Union, structs.Union{
		Query: sb.GetQuery(),
		IsAll: false,
	})

	return b
}

func (b *Builder) UnionAll(sb *Builder) *Builder {
	*b.selectQuery.Union = append(*b.selectQuery.Union, structs.Union{
		Query: sb.GetQuery(),
		IsAll: true,
	})

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
func (b *Builder) Build() (string, []interface{}, error) {
	// last query to be built and add to the union
	b.buildQuery()

	*b.selectQuery.Union = append(*b.selectQuery.Union, structs.Union{
		Query: b.query,
		IsAll: false,
	})

	sb := &strings.Builder{}
	sb.Grow(consts.StringBuffer_Middle_Query_Grow) // todo check if this is correct

	query := ""
	values := make([]interface{}, 0)
	for i, u := range *b.selectQuery.Union {
		cacheKey := generateCacheKey(&u)
		if cachedQuery, found := b.cache.Get(cacheKey); found {
			vals := make([]interface{}, 0, len(b.selectValues)+len(b.JoinBuilder.joinValues)+len(b.WhereBuilder.whereValues)+len(b.groupByValues))
			vals = append(vals, b.selectValues...)
			vals = append(vals, b.JoinBuilder.joinValues...)
			vals = append(vals, b.WhereBuilder.whereValues...)
			vals = append(vals, b.groupByValues...)
			sb.WriteString(cachedQuery)
			values = append(values, vals...)
			continue
		} else {
			q, v := b.dbBuilder.Build(cacheKey, u.Query, i, b.selectQuery.Union)

			b.cache.Set(cacheKey, q)

			sb.WriteString(q)
			values = append(values, v...)
		}
	}

	// remove the last UNION
	*b.selectQuery.Union = (*b.selectQuery.Union)[:len(*b.selectQuery.Union)-1]

	query = sb.String()

	return query, values, nil
}

func (b *Builder) buildQuery() {
	// preprocess WHERE
	wg := b.WhereBuilder.query.ConditionGroups
	if len(*b.WhereBuilder.query.Conditions) > 0 {
		*wg = append(*wg, structs.WhereGroup{
			Conditions:   *b.WhereBuilder.query.Conditions,
			Operator:     consts.LogicalOperator_AND,
			IsDummyGroup: true,
		})
		b.WhereBuilder.query.Conditions = &[]structs.Where{}
	}

	// preprocess ORDER BY
	o := b.orderByBuilder.Order

	b.query.Table = structs.Table{
		Name: b.selectQuery.Table,
	}
	b.query.Columns = b.selectQuery.Columns
	b.query.ConditionGroups = wg
	b.query.Joins = b.JoinBuilder.Joins
	b.query.Order = o
	b.query.Group = b.selectQuery.Group
	b.query.Limit = b.selectQuery.Limit
	b.query.Offset = b.selectQuery.Offset
	b.query.Lock = b.selectQuery.Lock
	b.query.SubQuery = b.selectQuery.SubQuery

}

func generateCacheKey(u *structs.Union) string {
	var sb strings.Builder
	sb.Grow(consts.StringBuffer_CacheKey_Grow)

	q := u.Query
	// Union key
	//sb.WriteString(fmt.Sprintf("%d", i))
	sb.WriteString(fmt.Sprintf("%t", u.IsAll))
	sb.WriteString("|")

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
				sb.WriteString(generateCacheKey(&structs.Union{Query: w.Query}))
				sb.WriteString(",")
			}
		}
	}

	key := sb.String()

	return key
}

func (b *Builder) GetQuery() *structs.Query {
	b.buildQuery()
	return b.query
}

func (b *Builder) GetStrategy() db.QueryBuilderStrategy {
	return b.dbBuilder
}

func (b *Builder) GetWhereBuilder() *WhereBuilder[Builder] {
	return &b.WhereBuilder
}

func (b *Builder) GetJoinBuilder() *JoinBuilder[Builder] {
	return &b.JoinBuilder
}

func (b *Builder) GetOrderByBuilder() *OrderByBuilder {
	return b.orderByBuilder
}

func (b *Builder) GetCache() cache.Cache {
	return b.cache
}
