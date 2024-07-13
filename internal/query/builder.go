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

func (b *Builder) Join(table string, my string, condition string, target string) *Builder {
	myTable := b.query.Table.Name
	// If a previous JOIN exists, retrieve the table name of that JOIN.
	if b.query.Joins != nil && len(*b.query.Joins) > 0 {
		myTable = (*b.query.Joins)[len(*b.query.Joins)-1].Name
	}
	*b.query.Joins = append(*b.query.Joins, structs.Join{
		Name: myTable,
		TargetNameMap: map[string]string{
			consts.Join_INNER: table,
		},
		SearchColumn:       my,
		SearchCondition:    condition,
		SearchTargetColumn: target,
	})
	return b
}

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
			Conditions: *b.query.Conditions,
			Operator:   consts.LogicalOperator_AND,
		})
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
