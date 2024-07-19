package query_test

import (
	"testing"

	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

func TestDeleteBuilder(t *testing.T) {
	tests := []struct {
		name           string
		setup          func() *query.DeleteBuilder
		expectedQuery  string
		expectedValues []interface{}
	}{
		{
			"Delete_all",
			func() *query.DeleteBuilder {
				return query.NewDeleteBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache()).
					Table("users").
					Delete()
			},
			"DELETE FROM users",
			[]interface{}{},
		},
		{
			"Delete_where",
			func() *query.DeleteBuilder {
				return query.NewDeleteBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache()).
					Table("users").
					Where("id", "=", 1).
					Delete()
			},
			"DELETE FROM users WHERE id = ?",
			[]interface{}{1},
		},
		{
			"Delete_JOINS",
			func() *query.DeleteBuilder {
				return query.NewDeleteBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache()).
					Table("users").
					Join("profiles", "users.id", "=", "profiles.user_id").
					Where("age", ">", 18).
					Delete()
			},
			"DELETE users FROM users INNER JOIN profiles ON users.id = profiles.user_id WHERE age > ?",
			[]interface{}{18},
		},
		/*
			{
				"Delete_using",
				func() *query.DeleteBuilder {
					return query.NewDeleteBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache()).
						Table("users").
						Using(query.NewBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache()).
							Table("profiles").
							Select("name", "age").
							Where("age", ">", 18).GetQuery()).
						Delete()
				},
				"DELETE users FROM users USING (SELECT name, age FROM profiles WHERE age > ?)",
				[]interface{}{18},
			},
		*/
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
