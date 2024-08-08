package interfaces

import (
	"strings"

	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type QueryBuilderStrategy interface {
	Select(sb *strings.Builder, columns *[]structs.Column, tableName string, joins *structs.Joins) []interface{}
	From(sb *strings.Builder, table string)
	Where(sb *strings.Builder, WhereGroups []structs.WhereGroup) []interface{}
	Join(sb *strings.Builder, joins *structs.Joins) []interface{}
	Union(sb *strings.Builder, unions *[]structs.Union, number int)
	OrderBy(sb *strings.Builder, order *[]structs.Order)
	Build(sb *strings.Builder, q *structs.Query, number int, unions *[]structs.Union) []interface{}

	Insert(q *structs.InsertQuery) (string, []interface{}, error)
	InsertBatch(q *structs.InsertQuery) (string, []interface{}, error)
	BuildInsert(q *structs.InsertQuery) (string, []interface{}, error)

	BuildUpdate(q *structs.UpdateQuery) (string, []interface{}, error)

	BuildDelete(q *structs.DeleteQuery) (string, []interface{}, error)
}
