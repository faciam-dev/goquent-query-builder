package base

import (
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

type FromBaseBuilder struct {
	u interfaces.SQLUtils
}

func NewFromBaseBuilder(util interfaces.SQLUtils) *FromBaseBuilder {
	return &FromBaseBuilder{
		u: util,
	}
}

func (f FromBaseBuilder) From(sb *strings.Builder, table string) {
	sb.WriteString("FROM ")
	f.u.EscapeIdentifier(sb, table)
}
