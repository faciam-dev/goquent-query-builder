package base

import (
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
	"github.com/faciam-dev/goquent-query-builder/internal/db/interfaces"
)

type OrderByBaseBuilder struct {
	u     interfaces.SQLUtils
	order *[]structs.Order
}

func NewOrderByBaseBuilder(util interfaces.SQLUtils, order *[]structs.Order) *OrderByBaseBuilder {
	return &OrderByBaseBuilder{
		u:     util,
		order: order,
	}
}

func (o OrderByBaseBuilder) OrderBy(sb *strings.Builder, order *[]structs.Order) {
	if order == nil || len(*order) == 0 {
		return
	}

	sb.WriteString(" ORDER BY ")

	for i := range *order {
		if i > 0 {
			sb.WriteString(", ")
		}
		if (*order)[i].Raw != "" {
			sb.WriteString((*order)[i].Raw)
			continue
		}
		if (*order)[i].Column == "" {
			continue
		}

		desc := "DESC"
		if (*order)[i].IsAsc {
			desc = "ASC"
		}
		sb.WriteString(o.u.EscapeIdentifier(sb, (*order)[i].Column))
		sb.WriteString(" ")
		sb.WriteString(desc)
	}
}
