package interfaces

import (
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type QueryBuilderStrategy interface {
	/*
		Select(sb *[]byte, columns *[]structs.Column, tableName string, joins *structs.Joins) []interface{}
		From(sb *[]byte, table string)
		Where(sb *[]byte, WhereGroups []structs.WhereGroup) []interface{}
		Join(sb *[]byte, joins *structs.Joins) []interface{}
		Union(sb *[]byte, unions *[]structs.Union, number int)
		OrderBy(sb *[]byte, order *[]structs.Order)
	*/
	Build(sb *[]byte, q *structs.Query, number int, unions *[]structs.Union) []interface{}

	Insert(q *structs.InsertQuery) (string, []interface{}, error)
	InsertBatch(q *structs.InsertQuery) (string, []interface{}, error)
	BuildInsert(q *structs.InsertQuery) (string, []interface{}, error)

	BuildUpdate(q *structs.UpdateQuery) (string, []interface{}, error)

	BuildDelete(q *structs.DeleteQuery) (string, []interface{}, error)
}
