package query_test

import (
	"testing"

	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

func TestInsertBuilder(t *testing.T) {
	tests := []struct {
		name           string
		setup          func() *query.InsertBuilder
		expectedQuery  string
		expectedValues []interface{}
	}{
		{
			"Insert",
			func() *query.InsertBuilder {
				return query.NewInsertBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).
					Table("users").
					Insert(map[string]interface{}{
						"name": "John Doe",
						"age":  30,
					})
			},
			"INSERT INTO users (age, name) VALUES (?, ?)",
			[]interface{}{30, "John Doe"},
		},
		{
			"InsertBatch",
			func() *query.InsertBuilder {
				return query.NewInsertBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).
					Table("users").
					InsertBatch([]map[string]interface{}{
						{
							"name": "John Doe",
							"age":  30,
						},
						{
							"name": "Jane Doe",
							"age":  25,
						},
					})
			},
			"INSERT INTO users (age, name) VALUES (?, ?), (?, ?)",
			[]interface{}{30, "John Doe", 25, "Jane Doe"},
		},
		{
			"InsertUsing",
			func() *query.InsertBuilder {
				return query.NewInsertBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).
					Table("users").
					InsertUsing([]string{"name", "age"}, query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).
						Table("profiles").
						Select("name", "age").
						Where("age", ">", 18))
			},
			"INSERT INTO users (name, age) SELECT name, age FROM profiles WHERE age > ?",
			[]interface{}{18},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			builder := tt.setup()
			query, values, _ := builder.Build()

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
