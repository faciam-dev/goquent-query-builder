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
	cache         *cache.AsyncQueryCache
	query         *structs.Query
	selectValues  []interface{}
	whereValues   []interface{}
	groupByValues []interface{}
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
		},
		selectValues:  []interface{}{},
		whereValues:   []interface{}{},
		groupByValues: []interface{}{},
	}
}

func NewBuilderWithQuery(dbBuilder db.QueryBuilderStrategy, query *structs.Query) *Builder {
	return &Builder{
		dbBuilder: dbBuilder,
		query:     query,
	}
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
	*b.query.Conditions = append(*b.query.Conditions, structs.Where{
		Column:    column,
		Condition: condition,
		Value:     value,
		Operator:  consts.LogicalOperator_AND,
	})
	b.whereValues = append(b.whereValues, value...)
	return b
}

func (b *Builder) OrWhere(column string, condition string, value ...interface{}) *Builder {
	*b.query.Conditions = append(*b.query.Conditions, structs.Where{
		Column:    column,
		Condition: condition,
		Value:     value,
		Operator:  consts.LogicalOperator_OR,
	})
	b.whereValues = append(b.whereValues, value...)
	return b
}

func (b *Builder) WhereQuery(column string, condition string, q *Builder) *Builder {
	return b.whereOrOrWhereQuery(column, condition, q, consts.LogicalOperator_AND)
}

func (b *Builder) OrWhereQuery(column string, condition string, q *Builder) *Builder {
	return b.whereOrOrWhereQuery(column, condition, q, consts.LogicalOperator_OR)
}

func (b *Builder) whereOrOrWhereQuery(column string, condition string, q *Builder, operator int) *Builder {
	*q.query.ConditionGroups = append(*q.query.ConditionGroups, structs.WhereGroup{
		Conditions:   *q.query.Conditions,
		IsDummyGroup: true,
	})

	sq := &structs.Query{
		ConditionGroups: q.query.ConditionGroups,
		Table:           q.query.Table,
		Columns:         q.query.Columns,
		Joins:           q.query.Joins,
		Order:           q.query.Order,
	}

	args := &structs.Where{
		Column:    column,
		Condition: condition,
		Query:     sq,
		Operator:  operator,
	}

	*b.query.Conditions = append(*b.query.Conditions, *args)

	return b
}

func (b *Builder) WhereGroup(fn func(b *Builder) *Builder) *Builder {
	if len(*b.query.Conditions) > 0 {
		*b.query.ConditionGroups = append(*b.query.ConditionGroups, structs.WhereGroup{
			Conditions:   *b.query.Conditions,
			Operator:     consts.LogicalOperator_AND,
			IsDummyGroup: true,
		})
		*b.query.Conditions = []structs.Where{}
	}

	cQ := fn(b)

	*b.query.ConditionGroups = append(*b.query.ConditionGroups, structs.WhereGroup{
		Conditions: *cQ.query.Conditions,
		Subgroups:  []structs.WhereGroup{},
		Operator:   consts.LogicalOperator_AND,
	})
	*cQ.query.Conditions = []structs.Where{}

	return b
}

func (b *Builder) OrWhereGroup(fn func(b *Builder) *Builder) *Builder {
	if len(*b.query.Conditions) > 0 {
		*b.query.ConditionGroups = append(*b.query.ConditionGroups, structs.WhereGroup{
			Conditions:   *b.query.Conditions,
			Operator:     consts.LogicalOperator_OR,
			IsDummyGroup: true,
		})
		*b.query.Conditions = []structs.Where{}
	}

	cQ := fn(b)

	*b.query.ConditionGroups = append(*b.query.ConditionGroups, structs.WhereGroup{
		Conditions: *cQ.query.Conditions,
		Subgroups:  []structs.WhereGroup{},
		Operator:   consts.LogicalOperator_OR,
	})
	*cQ.query.Conditions = []structs.Where{}

	return b
}

// Join adds a JOIN clause.
func (b *Builder) Join(table string, my string, condition string, target string) *Builder {
	return b.joinCommon(consts.Join_INNER, table, my, condition, target)
}

// LeftJoin adds a LEFT JOIN clause.
func (b *Builder) LeftJoin(table string, my string, condition string, target string) *Builder {
	return b.joinCommon(consts.Join_LEFT, table, my, condition, target)
}

// RightJoin adds a RIGHT JOIN clause.
func (b *Builder) RightJoin(table string, my string, condition string, target string) *Builder {
	return b.joinCommon(consts.Join_RIGHT, table, my, condition, target)
}

// joinCommon is a helper function for JOIN, LEFT JOIN, and RIGHT JOIN.
func (b *Builder) joinCommon(joinType string, table string, my string, condition string, target string) *Builder {
	myTable := b.query.Table.Name
	// If a previous JOIN exists, retrieve the table name of that JOIN.
	if b.query.Joins != nil && len(*b.query.Joins) > 0 {
		myTable = (*b.query.Joins)[len(*b.query.Joins)-1].Name
	}
	*b.query.Joins = append(*b.query.Joins, structs.Join{
		Name: myTable,
		TargetNameMap: map[string]string{
			joinType: table,
		},
		SearchColumn:       my,
		SearchCondition:    condition,
		SearchTargetColumn: target,
	})
	return b
}

// CrossJoin adds a CROSS JOIN clause.
func (b *Builder) CrossJoin(table string) *Builder {
	myTable := b.query.Table.Name
	// If a previous JOIN exists, retrieve the table name of that JOIN.
	if b.query.Joins != nil && len(*b.query.Joins) > 0 {
		myTable = (*b.query.Joins)[len(*b.query.Joins)-1].Name
	}
	*b.query.Joins = append(*b.query.Joins, structs.Join{
		Name: myTable,
		TargetNameMap: map[string]string{
			consts.Join_CROSS: table,
		},
	})
	return b
}

// OrderBy adds an ORDER BY clause.
func (b *Builder) OrderBy(column string, ascDesc string) *Builder {
	ascDesc = strings.ToUpper(ascDesc)

	if ascDesc == consts.Order_ASC {
		*b.query.Order = append(*b.query.Order, structs.Order{
			Column: column,
			IsAsc:  consts.Order_FLAG_ASC,
		})
	} else if ascDesc == consts.Order_DESC {
		*b.query.Order = append(*b.query.Order, structs.Order{
			Column: column,
			IsAsc:  consts.Order_FLAG_DESC,
		})
	}
	return b
}

// ReOrder removes all ORDER BY clauses.
func (b *Builder) ReOrder() *Builder {
	*b.query.Order = []structs.Order{}
	return b
}

// OrderByRaw adds a raw ORDER BY clause.
func (b *Builder) OrderByRaw(raw string) *Builder {
	*b.query.Order = append(*b.query.Order, structs.Order{
		Raw: raw,
	})
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
		values = append(values, b.whereValues...)
		values = append(values, b.groupByValues...)
		return cachedQuery, values
	}

	// preprocess WHERE
	if len(*b.query.Conditions) > 0 {
		*b.query.ConditionGroups = append(*b.query.ConditionGroups, structs.WhereGroup{
			Conditions:   *b.query.Conditions,
			Operator:     consts.LogicalOperator_AND,
			IsDummyGroup: true,
		})
		b.query.Conditions = &[]structs.Where{}
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
