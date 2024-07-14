package query_test

import (
	"testing"

	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

func TestBuilder(t *testing.T) {
	tests := []struct {
		name           string
		setup          func() *query.Builder
		expectedQuery  string
		expectedValues []interface{}
	}{
		{
			"Select",
			func() *query.Builder {
				return query.NewBuilder(db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache()).Select("id", "name")
			},
			"SELECT id, name FROM ",
			nil,
		},
		{
			"From",
			func() *query.Builder {
				return query.NewBuilder(db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache()).Table("users")
			},
			"SELECT  FROM users",
			nil,
		},
		{
			"Where",
			func() *query.Builder {
				return query.NewBuilder(db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache()).Where("age", ">", 18)
			},
			"SELECT  FROM  WHERE age > ?",
			[]interface{}{18},
		},
		{
			"WhereQuery",
			func() *query.Builder {
				sq := query.NewBuilder(db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache()).Select("id").Table("users").Where("name", "=", "John")
				return query.NewBuilder(db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache()).WhereQuery("user_id", "IN", sq).Where("city", "=", "New York")
			},
			"SELECT  FROM  WHERE user_id IN (SELECT id FROM users WHERE name = ?) AND city = ?",
			[]interface{}{"John", "New York"},
		},
		{
			"WhereGroup",
			func() *query.Builder {
				return query.NewBuilder(db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache()).
					WhereGroup(func(b *query.Builder) *query.Builder {
						return b.Where("age", ">", 18).Where("name", "=", "John")
					})
			},
			"SELECT  FROM  WHERE (age > ? AND name = ?)",
			[]interface{}{18, "John"},
		},
		{
			"WhereGroup_And",
			func() *query.Builder {
				return query.NewBuilder(db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache()).
					WhereGroup(func(b *query.Builder) *query.Builder {
						return b.Where("age", ">", 18).Where("name", "=", "John")
					}).Where("city", "=", "New York")
			},
			"SELECT  FROM  WHERE (age > ? AND name = ?) AND city = ?",
			[]interface{}{18, "John", "New York"},
		},
		{
			"Join",
			func() *query.Builder {
				return query.NewBuilder(db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache()).Join("orders", "users.id", "=", "orders.user_id")
			},
			"SELECT orders.*, .* FROM  INNER JOIN orders ON users.id = orders.user_id",
			nil,
		},
		{
			"OrderBy",
			func() *query.Builder {
				return query.NewBuilder(db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache()).OrderBy("name", "asc")
			},
			"SELECT  FROM  ORDER BY name ASC",
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			builder := tt.setup()
			query, values := builder.Build()

			if query != tt.expectedQuery {
				t.Errorf("expected '%s' but got '%s'", tt.expectedQuery, query)
			}

			if len(values) != len(tt.expectedValues) {
				t.Errorf("expected values %v but got %v", tt.expectedValues, values)
			}

			for i := range values {
				if values[i] != tt.expectedValues[i] {
					t.Errorf("expected value %v at index %d but got %v", tt.expectedValues[i], i, values[i])
				}
			}
		})
	}
}
