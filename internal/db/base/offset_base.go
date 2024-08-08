package base

import (
	"strconv"
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type OffsetBaseBuilder struct {
}

func NewOffsetBaseBuilder() *OffsetBaseBuilder {
	return &OffsetBaseBuilder{}
}

func (OffsetBaseBuilder) Offset(sb *strings.Builder, offset structs.Offset) {
	if offset.Offset == 0 {
		return
	}

	sb.WriteString(" OFFSET ")
	sb.WriteString(strconv.FormatInt(offset.Offset, 10))
}
