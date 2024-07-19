package query

import (
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/consts"
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type OrderByBuilder struct {
	Order *[]structs.Order
}

func NewOrderByBuilder(o *[]structs.Order) *OrderByBuilder {
	return &OrderByBuilder{
		Order: o,
	}
}

// OrderBy adds an ORDER BY clause.
func (b *OrderByBuilder) OrderBy(column string, ascDesc string) *OrderByBuilder {
	ascDesc = strings.ToUpper(ascDesc)

	if ascDesc == consts.Order_ASC {
		*b.Order = append(*b.Order, structs.Order{
			Column: column,
			IsAsc:  consts.Order_FLAG_ASC,
		})
	} else if ascDesc == consts.Order_DESC {
		*b.Order = append(*b.Order, structs.Order{
			Column: column,
			IsAsc:  consts.Order_FLAG_DESC,
		})
	}
	return b
}

// ReOrder removes all ORDER BY clauses.
func (b *OrderByBuilder) ReOrder() *OrderByBuilder {
	*b.Order = []structs.Order{}
	return b
}

// OrderByRaw adds a raw ORDER BY clause.
func (b *OrderByBuilder) OrderByRaw(raw string) *OrderByBuilder {
	*b.Order = append(*b.Order, structs.Order{
		Raw: raw,
	})
	return b
}
