package db

import (
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type OrderByBaseBuilder struct {
	order *[]structs.Order
}

func NewOrderByBaseBuilder(order *[]structs.Order) *OrderByBaseBuilder {
	return &OrderByBaseBuilder{
		order: order,
	}
}

func (OrderByBaseBuilder) OrderBy(order *[]structs.Order) string {
	if len(*order) == 0 {
		return ""
	}

	var sb strings.Builder
	rawOrderQuerys := make([]string, 0, len(*order))
	orders := make([]string, 0, len(*order))
	for _, order := range *order {
		if order.Raw != "" {
			rawOrderQuerys = append(rawOrderQuerys, order.Raw)
			continue
		}
		if order.Column == "" {
			continue
		}
		desc := "DESC"
		if order.IsAsc {
			desc = "ASC"
		}
		orders = append(orders, order.Column+" "+desc)
	}
	if len(rawOrderQuerys) > 0 {
		sb.WriteString(" ORDER BY ")
		sb.WriteString(strings.Join(rawOrderQuerys, ", "))
	}
	if len(orders) > 0 {
		if len(rawOrderQuerys) > 0 {
			sb.WriteString(", ")
		} else {
			sb.WriteString(" ORDER BY ")
		}
		sb.WriteString(strings.Join(orders, ", "))
	}

	return sb.String()
}
