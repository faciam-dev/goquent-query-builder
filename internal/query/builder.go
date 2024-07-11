package query

import (
	"fmt"
	"strings"
)

type Builder struct {
	table      string
	columns    []string
	conditions []string
	joins      []string
	orderBy    []string
	values     []interface{}
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Table(table string) *Builder {
	b.table = table
	return b
}

func (b *Builder) Select(columns ...string) *Builder {
	b.columns = append(b.columns, columns...)
	return b
}

func (b *Builder) Where(condition string, value interface{}) *Builder {
	b.conditions = append(b.conditions, condition)
	b.values = append(b.values, value)
	return b
}

func (b *Builder) Join(joinType, table, condition string) *Builder {
	b.joins = append(b.joins, fmt.Sprintf("%s JOIN %s ON %s", joinType, table, condition))
	return b
}

func (b *Builder) OrderBy(orderBy ...string) *Builder {
	b.orderBy = append(b.orderBy, orderBy...)
	return b
}

func (b *Builder) Build() (string, []interface{}) {
	query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(b.columns, ", "), b.table)
	if len(b.joins) > 0 {
		query += " " + strings.Join(b.joins, " ")
	}
	if len(b.conditions) > 0 {
		query += " WHERE " + strings.Join(b.conditions, " AND ")
	}
	if len(b.orderBy) > 0 {
		query += " ORDER BY " + strings.Join(b.orderBy, ", ")
	}
	return query, b.values
}
