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

func (OrderByBaseBuilder) OrderBy(sb *strings.Builder, order *[]structs.Order) {
	if order == nil || len(*order) == 0 {
		return
	}

	//	sb.Grow(consts.StringBuffer_OrderBy_Grow)
	//rawOrderQuerys := make([]string, 0, len(*order))
	//orders := make([]string, 0, len(*order))
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
		sb.WriteString(order.Column)
		sb.WriteString(" ")
		sb.WriteString(desc)
	}

	/*

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
			for i, rawOrderQuery := range rawOrderQuerys {
				if i > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(rawOrderQuery)
			}
		}

		if len(orders) > 0 {
			if len(rawOrderQuerys) == 0 {
				sb.WriteString(" ORDER BY ")
			}
			for i, order := range orders {
				if i > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(order)
			}
		}
	*/
}
