package base

import (
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type UnionBaseBuilder struct {
}

func NewUnionBaseBuilder() *UnionBaseBuilder {
	return &UnionBaseBuilder{}
}

func (ub *UnionBaseBuilder) Union(sb *strings.Builder, unions *[]structs.Union, number int) {
	if unions == nil {
		return
	}

	ub.buildUnionStatement(sb, unions, number)
}

func (ub *UnionBaseBuilder) buildUnionStatement(sb *strings.Builder, unions *[]structs.Union, number int) {
	if (*unions)[number].Query != nil {
		if len(*unions) > number+1 {
			if (*unions)[number].IsAll {
				sb.WriteString(" UNION ALL ")
			} else {
				sb.WriteString(" UNION ")
			}
		}
	}
}
