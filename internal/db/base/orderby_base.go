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

	for i, order := range *order {
		if i > 0 {
			sb.WriteString(", ")
		}
		if order.Raw != "" {
			sb.WriteString(order.Raw)
			continue
		}
		if order.Column == "" {
			continue
		}

		desc := "DESC"
		if order.IsAsc {
			desc = "ASC"
		}
		sb.WriteString(o.u.EscapeIdentifier(order.Column))
		sb.WriteString(" ")
		sb.WriteString(desc)
	}
}
