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
				return query.NewDeleteBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache(100)).
					Table("users").
					Delete()
			},
			"DELETE FROM users",
			[]interface{}{},
		},
		{
			"Delete_where",
			func() *query.DeleteBuilder {
				return query.NewDeleteBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache(100)).
					Table("users").
					Where("id", "=", 1).
					Delete()
			},
			"DELETE FROM users WHERE id = ?",
			[]interface{}{1},
		},
		{
			"Delete_where_not",
			func() *query.DeleteBuilder {
				return query.NewDeleteBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache(100)).
					Table("users").
					Where("id", "!=", 1).
					OrWhereNot(func(b *query.WhereBuilder) *query.WhereBuilder {
						return b.Where("age", ">", 18).Where("name", "=", "John")
					}).
					Delete()
			},
			"DELETE FROM users WHERE id != ? OR NOT (age > ? AND name = ?)",
			[]interface{}{1, 18, "John"},
		},
		{
			"Delete_where_all",
			func() *query.DeleteBuilder {
				return query.NewDeleteBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache(100)).
					Table("users").
					Where("id", ">", 10000).
					WhereAll([]string{"firstname", "lastname"}, "LIKE", "%test%").
					Delete()
			},
			"DELETE FROM users WHERE id > ? AND (firstname LIKE ? AND lastname LIKE ?)",
			[]interface{}{10000, "%test%", "%test%"},
		},
		{
			"Delete_JOINS",
			func() *query.DeleteBuilder {
				return query.NewDeleteBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache(100)).
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
					return query.NewDeleteBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache(100)).
						Table("users").
						Using(query.NewBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache(100)).
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
