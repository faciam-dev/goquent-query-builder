package db

import "github.com/faciam-dev/goquent-query-builder/internal/common/structs"

type QueryBuilderStrategy interface {
	Select(columns *[]structs.Column, joinedTablesForSelect *[]structs.Column) ([]string, []interface{})
	From(table string) string
	Where(WhereGroups *[]structs.WhereGroup) (string, []interface{})
	Join(tableName string, joins *structs.Joins) (*[]structs.Column, string, []interface{})
	OrderBy(order *[]structs.Order) string
	Build(q *structs.Query) (string, []interface{})

	Insert(q *structs.InsertQuery) (string, []interface{})
	InsertBatch(q *structs.InsertQuery) (string, []interface{})
	BuildInsert(q *structs.InsertQuery) (string, []interface{})

	BuildUpdate(q *structs.UpdateQuery) (string, []interface{})

	BuildDelete(q *structs.DeleteQuery) (string, []interface{})
}
