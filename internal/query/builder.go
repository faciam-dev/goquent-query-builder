package query

import (
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
)

type Builder struct {
	dbBuilder db.QueryBuilderStrategy
	query     *structs.Query
}

func NewBuilder(dbBuilder db.QueryBuilderStrategy) *Builder {
	return &Builder{
		dbBuilder: dbBuilder,
		query: &structs.Query{
			Columns:         &[]structs.Column{},
			Conditions:      &[]structs.Where{},
			ConditionGroups: &[]structs.WhereGroup{},
			Joins:           &[]structs.Join{},
			Order:           &[]structs.Order{},
			SubQuery:        &[]structs.Query{},
		},
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
	// preprocess WHERE
	if len(*b.query.Conditions) > 0 {
		*b.query.ConditionGroups = append(*b.query.ConditionGroups, structs.WhereGroup{
			Conditions: *b.query.Conditions,
			Operator:   consts.LogicalOperator_AND,
		})
	}

	return b.dbBuilder.Build(b.query)
}
