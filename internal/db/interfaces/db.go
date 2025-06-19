package interfaces

import (
	"github.com/faciam-dev/goquent-query-builder/internal/common/structs"
)

type QueryBuilderStrategy interface {
	Build(sb *[]byte, q *structs.Query, number int, unions *[]structs.Union) ([]interface{}, error)

	Insert(q *structs.InsertQuery) (string, []interface{}, error)
	InsertBatch(q *structs.InsertQuery) (string, []interface{}, error)
	BuildInsert(q *structs.InsertQuery) (string, []interface{}, error)
	InsertIgnore(q *structs.InsertQuery) (string, []interface{}, error)
	Upsert(q *structs.InsertQuery) (string, []interface{}, error)

	BuildUpdate(q *structs.UpdateQuery) (string, []interface{}, error)

	BuildDelete(q *structs.DeleteQuery) (string, []interface{}, error)
}
