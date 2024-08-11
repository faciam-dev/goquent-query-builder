package base

import (
	"strconv"

	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type OffsetBaseBuilder struct {
}

func NewOffsetBaseBuilder() *OffsetBaseBuilder {
	return &OffsetBaseBuilder{}
}

func (OffsetBaseBuilder) Offset(sb *[]byte, offset structs.Offset) {
	if offset.Offset == 0 {
		return
	}

	*sb = append(*sb, " OFFSET "...)
	*sb = strconv.AppendInt(*sb, offset.Offset, 10)
}
