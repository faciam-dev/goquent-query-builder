package query

import (
	"strings"
	"sync"

	"github.com/faciam-dev/goquent-query-builder/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

/*
type QueryBuilder struct {
	whereBuilder   WhereBuilder[Builder]
	joinBuilder    JoinBuilder[Builder]
	orderByBuilder OrderByBuilder
}
*/

type SelectBuilder struct {
	dbBuilder     interfaces.QueryBuilderStrategy
	cache         cache.Cache
	query         *structs.Query
	selectQuery   *structs.SelectQuery
	selectValues  []interface{}
	groupByValues []interface{}
	*WhereBuilder[SelectBuilder]
	*JoinBuilder[SelectBuilder]
	*OrderByBuilder[SelectBuilder]
	BaseBuilder
}

func NewBuilder(dbBuilder interfaces.QueryBuilderStrategy, cache cache.Cache) *SelectBuilder {
	b := &SelectBuilder{
		dbBuilder: dbBuilder,
		cache:     cache,
		query: &structs.Query{
			Table:           structs.Table{},
			Columns:         &[]structs.Column{},
			ConditionGroups: []structs.WhereGroup{},
			Joins:           &structs.Joins{},
			Order:           &[]structs.Order{},
			Group:           &structs.GroupBy{},
			Limit:           structs.Limit{},
			Offset:          structs.Offset{},
			Lock:            &structs.Lock{},
		},
		selectQuery: &structs.SelectQuery{
			Table:   "",
			Columns: &[]structs.Column{},
			Limit:   structs.Limit{},
			Group:   &structs.GroupBy{},
			Offset:  structs.Offset{},
			Lock:    &structs.Lock{},
			Union:   &[]structs.Union{},
		},
		selectValues:  []interface{}{},
		groupByValues: []interface{}{},
		//joinBuilder:    NewJoinBuilder(dbBuilder, cache),
	}

	whereBuilder := NewWhereBuilder[SelectBuilder](dbBuilder, cache)
	whereBuilder.SetParent(b)
	b.WhereBuilder = whereBuilder

	joinBuilder := NewJoinBuilder[SelectBuilder](dbBuilder, cache)
	joinBuilder.SetParent(b)
	b.JoinBuilder = joinBuilder

	orderByBuilder := NewOrderByBuilder[SelectBuilder](dbBuilder, cache)
	orderByBuilder.SetParent(b)
	b.OrderByBuilder = orderByBuilder

	return b
}

var stringbufPool = sync.Pool{
	New: func() interface{} {
		return new(strings.Builder)
	},
}

/*
func (qb *SelectBuilder) SetBuilders(where *WhereBuilder[SelectBuilder], join *JoinBuilder[SelectBuilder], orderBy *OrderByBuilder[SelectBuilder]) {
	qb.WhereBuilder = where
	qb.JoinBuilder = join
	qb.OrderByBuilder = orderBy
}

func (b *SelectBuilder) SetWhereBuilder(whereBuilder *WhereBuilder[SelectBuilder]) {
	b.WhereBuilder = whereBuilder
}

func (b *SelectBuilder) SetJoinBuilder(joinBuilder *JoinBuilder[SelectBuilder]) {
	b.JoinBuilder = joinBuilder
}

func (b *SelectBuilder) SetOrderByBuilder(orderByBuilder *OrderByBuilder[SelectBuilder]) {
	b.OrderByBuilder = orderByBuilder
}
*/

func (b *SelectBuilder) Table(table string) *SelectBuilder {
	b.selectQuery.Table = table
	return b
}

func (b *SelectBuilder) Select(columns ...string) *SelectBuilder {
	for _, column := range columns {
		*b.selectQuery.Columns = append(*b.selectQuery.Columns, structs.Column{Name: column})
	}
	return b
}

func (b *SelectBuilder) SelectRaw(raw string, value ...interface{}) *SelectBuilder {
	b.selectValues = append(b.selectValues, value...)
	*b.selectQuery.Columns = append(*b.selectQuery.Columns, structs.Column{Raw: raw, Values: value})
	return b
}

// Count adds a COUNT aggregate function to the query.
func (b *SelectBuilder) Count(columns ...string) *SelectBuilder {
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

func (b *SelectBuilder) aggregate(column string, aggregateFunc string) *SelectBuilder {
	*b.selectQuery.Columns = append(*b.selectQuery.Columns, structs.Column{
		Name:     column,
		Function: aggregateFunc,
	})
	return b
}

// Max adds a MAX aggregate function to the query.
func (b *SelectBuilder) Max(column string) *SelectBuilder {
	return b.aggregate(column, "MAX")
}

// Min adds a MIN aggregate function to the query.
func (b *SelectBuilder) Min(column string) *SelectBuilder {
	return b.aggregate(column, "MIN")
}

// Sum adds a SUM aggregate function to the query.
func (b *SelectBuilder) Sum(column string) *SelectBuilder {
	return b.aggregate(column, "SUM")
}

// Avg adds an AVG aggregate function to the query.
func (b *SelectBuilder) Avg(column string) *SelectBuilder {
	return b.aggregate(column, "AVG")
}

func (b *SelectBuilder) Distinct(column ...string) *SelectBuilder {
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

func (b *SelectBuilder) Union(sb *SelectBuilder) *SelectBuilder {
	*b.selectQuery.Union = append(*b.selectQuery.Union, structs.Union{
		Query: sb.GetQuery(),
		IsAll: false,
	})

	return b
}

func (b *SelectBuilder) UnionAll(sb *SelectBuilder) *SelectBuilder {
	*b.selectQuery.Union = append(*b.selectQuery.Union, structs.Union{
		Query: sb.GetQuery(),
		IsAll: true,
	})

	return b
}

/*
// OrderBy adds an ORDER BY clause.
func (b *Builder) OrderBy(column string, ascDesc string) *Builder {
	b.OrderByBuilder.OrderBy(column, ascDesc)

	return b
}

// ReOrder removes all ORDER BY clauses.
func (b *Builder) ReOrder() *Builder {
	b.OrderByBuilder.ReOrder()
	return b
}

// OrderByRaw adds a raw ORDER BY clause.
func (b *Builder) OrderByRaw(raw string) *Builder {
	b.OrderByBuilder.OrderByRaw(raw)
	return b
}
*/

// GroupBy adds a GROUP BY clause.
func (b *SelectBuilder) GroupBy(columns ...string) *SelectBuilder {
	*b.selectQuery.Group = structs.GroupBy{
		Columns: columns,
		Having:  &[]structs.Having{},
	}
	return b
}

// Having adds a HAVING clause with an AND operator.
func (b *SelectBuilder) Having(column string, condition string, value interface{}) *SelectBuilder {
	*b.selectQuery.Group.Having = append(*b.selectQuery.Group.Having, structs.Having{
		Column:    column,
		Condition: condition,
		Value:     value,
		Operator:  consts.LogicalOperator_AND,
	})
	return b
}

// HavingRaw adds a raw HAVING clause with an AND operator.
func (b *SelectBuilder) HavingRaw(raw string) *SelectBuilder {
	*b.selectQuery.Group.Having = append(*b.selectQuery.Group.Having, structs.Having{
		Raw:      raw,
		Operator: consts.LogicalOperator_AND,
	})
	return b
}

// OrHaving adds a HAVING clause with an OR operator.
func (b *SelectBuilder) OrHaving(column string, condition string, value interface{}) *SelectBuilder {
	*b.selectQuery.Group.Having = append(*b.selectQuery.Group.Having, structs.Having{
		Column:    column,
		Condition: condition,
		Value:     value,
		Operator:  consts.LogicalOperator_OR,
	})
	return b
}

// OrHavingRaw adds a raw HAVING clause with an OR operator.
func (b *SelectBuilder) OrHavingRaw(raw string) *SelectBuilder {
	*b.selectQuery.Group.Having = append(*b.selectQuery.Group.Having, structs.Having{
		Raw:      raw,
		Operator: consts.LogicalOperator_OR,
	})
	return b
}

func (b *SelectBuilder) Limit(limit int64) *SelectBuilder {
	b.selectQuery.Limit.Limit = limit
	return b
}

func (b *SelectBuilder) Offset(offset int64) *SelectBuilder {
	b.selectQuery.Offset.Offset = offset
	return b
}

func (b *SelectBuilder) SharedLock() *SelectBuilder {
	b.selectQuery.Lock = &structs.Lock{
		LockType: consts.Lock_SHARE_MODE,
	}
	return b
}

func (b *SelectBuilder) LockForUpdate() *SelectBuilder {
	b.selectQuery.Lock = &structs.Lock{
		LockType: consts.Lock_FOR_UPDATE,
	}
	return b
}

// Build generates the SQL query string and parameter values based on the query builder's current state.
// It returns the generated query string and a slice of parameter values.
func (b *SelectBuilder) Build() (string, []interface{}, error) {
	// last query to be built and add to the union
	b.buildQuery()

	*b.selectQuery.Union = append(*b.selectQuery.Union, structs.Union{
		Query: b.query,
		IsAll: false,
	})

	sb := stringbufPool.Get().(*strings.Builder)
	sb.Reset()
	//sb.Grow(consts.StringBuffer_Short_Query_Grow)

	//query := ""
	values := make([]interface{}, 0)

	cacheKey := ""
	for i := range *b.selectQuery.Union {
		cacheKey += generateCacheKey(&(*b.selectQuery.Union)[i])
	}

	// grow the string builder based on the length of the cache key
	if len(cacheKey) < consts.StringBuffer_Short_Query_Grow {
		sb.Grow(consts.StringBuffer_Short_Query_Grow)
	} else if len(cacheKey) < consts.StringBuffer_Middle_Query_Grow {
		sb.Grow(consts.StringBuffer_Middle_Query_Grow)
	} else {
		sb.Grow(consts.StringBuffer_Long_Query_Grow)
	}

	for i := range *b.selectQuery.Union {
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
			_, v := b.dbBuilder.Build(sb, cacheKey, (*b.selectQuery.Union)[i].Query, i, b.selectQuery.Union)

			//sb.WriteString(q)
			values = append(values, v...)
		}
	}

	query := sb.String()
	b.cache.Set(cacheKey, query)

	// remove the last UNION
	*b.selectQuery.Union = (*b.selectQuery.Union)[:len(*b.selectQuery.Union)-1]

	stringbufPool.Put(sb)

	return query, values, nil
}

func (b *SelectBuilder) buildQuery() {
	// preprocess WHERE
	if len(*b.WhereBuilder.query.Conditions) > 0 {
		b.WhereBuilder.query.ConditionGroups = append(b.WhereBuilder.query.ConditionGroups,
			structs.WhereGroup{
				Conditions:   *b.WhereBuilder.query.Conditions,
				Operator:     consts.LogicalOperator_AND,
				IsDummyGroup: true,
			})
		b.WhereBuilder.query.Conditions = &[]structs.Where{}
	}

	// preprocess ORDER BY
	o := b.OrderByBuilder.Order

	b.query.Table = structs.Table{
		Name: b.selectQuery.Table,
	}
	b.query.Columns = b.selectQuery.Columns
	b.query.ConditionGroups = b.WhereBuilder.query.ConditionGroups
	b.query.Joins = b.JoinBuilder.Joins
	b.query.Order = o
	b.query.Group = b.selectQuery.Group
	b.query.Limit = b.selectQuery.Limit
	b.query.Offset = b.selectQuery.Offset
	b.query.Lock = b.selectQuery.Lock

}

func generateCacheKey(u *structs.Union) string {
	sb := stringbufPool.Get().(*strings.Builder)
	sb.Reset()
	sb.Grow(consts.StringBuffer_CacheKey_Grow)

	q := u.Query
	// Union key
	//sb.WriteString(fmt.Sprintf("%d", i))
	if u.IsAll {
		sb.WriteString("ALL")
	} else {
		sb.WriteString("UNION")
	}
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
		for j := range *q.Joins.Joins {
			sb.WriteString((*q.Joins.Joins)[j].Name)
			sb.WriteString(",")
			sb.WriteString((*q.Joins.Joins)[j].SearchColumn)
			sb.WriteString(",")
			sb.WriteString((*q.Joins.Joins)[j].SearchCondition)
			sb.WriteString(",")
			sb.WriteString((*q.Joins.Joins)[j].SearchTargetColumn)
			sb.WriteString(",")
		}
	}
	if q.Joins.JoinClauses != nil {
		for _, jc := range *q.Joins.JoinClauses {
			for _, on := range *jc.On {
				sb.WriteString(on.Column)
				sb.WriteString(",")
				sb.WriteString(on.Condition)
				sb.WriteString(",")
				sb.WriteString(on.Value.(string))
				sb.WriteString(",")
				if on.Operator == consts.LogicalOperator_OR {
					sb.WriteString("OR")
				} else {
					sb.WriteString("AND")
				}
				sb.WriteString(",")
			}
		}
	}
	sb.WriteString("|")

	// Condition key
	for c := range q.ConditionGroups {
		if (q.ConditionGroups)[c].Operator == consts.LogicalOperator_AND {
			sb.WriteString("AND")
		} else {
			sb.WriteString("OR")
		}
		sb.WriteString(",")
		for _, w := range (q.ConditionGroups)[c].Conditions {
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
	stringbufPool.Put(sb)

	return key
}

func (b *SelectBuilder) GetQuery() *structs.Query {
	b.buildQuery()
	return b.query
}

func (b *SelectBuilder) GetStrategy() interfaces.QueryBuilderStrategy {
	return b.dbBuilder
}

func (b *SelectBuilder) GetWhereBuilder() *WhereBuilder[SelectBuilder] {
	return b.WhereBuilder
}

func (b *SelectBuilder) GetJoinBuilder() *JoinBuilder[SelectBuilder] {
	return b.JoinBuilder
}

func (b *SelectBuilder) GetOrderByBuilder() *OrderByBuilder[SelectBuilder] {
	return b.OrderByBuilder
}

func (b *SelectBuilder) GetCache() cache.Cache {
	return b.cache
}
