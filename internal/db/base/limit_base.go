package base

import (
	"fmt"
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type LimitBaseBuilder struct {
}

func NewLimitBaseBuilder() *LimitBaseBuilder {
	return &LimitBaseBuilder{}
}

func (LimitBaseBuilder) Limit(sb *strings.Builder, limit *structs.Limit) {
	if limit == nil || limit.Limit == 0 {
		return
	}

	sb.WriteString(" LIMIT ")
	sb.WriteString(fmt.Sprint(limit.Limit))
}
