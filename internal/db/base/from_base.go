package base

import "strings"

type FromBaseBuilder struct {
}

func NewFromBaseBuilder() *FromBaseBuilder {
	return &FromBaseBuilder{}
}

func (FromBaseBuilder) From(sb *strings.Builder, table string) {
	sb.WriteString("FROM ")
	sb.WriteString(table)
}
