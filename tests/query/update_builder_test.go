package query_test

import (
	"testing"

	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

func TestUpdateBuilder(t *testing.T) {
	tests := []struct {
		name           string
		setup          func() *query.UpdateBuilder
		expectedQuery  string
		expectedValues []interface{}
	}{
		{
			"Update_all",
			func() *query.UpdateBuilder {
				return query.NewUpdateBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache(100)).
					Table("users").
					Update(map[string]interface{}{
						"name": "Joe",
						"age":  31,
					})
			},
			"UPDATE users SET age = ?, name = ?",
			[]interface{}{31, "Joe"},
		},
		{
			"Update_where",
			func() *query.UpdateBuilder {
				return query.NewUpdateBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache(100)).
					Table("users").
					Where("id", "=", 1).
					Update(map[string]interface{}{
						"name": "Joe",
						"age":  31,
					})
			},
			"UPDATE users SET age = ?, name = ? WHERE id = ?",
			[]interface{}{31, "Joe", 1},
		},
		{
			"Update_where_not",
			func() *query.UpdateBuilder {

				return query.NewUpdateBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache(100)).
					//					SetParent(query.NewUpdateBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache(100))).
					Table("users").
					Where("id", "!=", 1).
					OrWhereNot(func(b *query.WhereBuilder[query.UpdateBuilder]) *query.WhereBuilder[query.UpdateBuilder] {
						return &b.Where("age", ">", 18).Where("name", "=", "John").WhereBuilder
					}).
					Update(map[string]interface{}{
						"name": "Joe",
						"age":  31,
					})
			},
			"UPDATE users SET age = ?, name = ? WHERE id != ? OR NOT (age > ? AND name = ?)",
			[]interface{}{31, "Joe", 1, 18, "John"},
		},
		{
			"Update_where_any",
			func() *query.UpdateBuilder {
				return query.NewUpdateBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache(100)).
					Table("users").
					WhereAny([]string{"name", "note"}, "LIKE", "%test%").
					Update(map[string]interface{}{
						"name": "Joe",
						"age":  31,
					})
			},
			"UPDATE users SET age = ?, name = ? WHERE (name LIKE ? OR note LIKE ?)",
			[]interface{}{31, "Joe", "%test%", "%test%"},
		},
		{
			"Update_where_in",
			func() *query.UpdateBuilder {
				return query.NewUpdateBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache(100)).
					Table("users").
					WhereIn("id", []interface{}{1, 2, 3}).
					Update(map[string]interface{}{
						"name": "Joe",
						"age":  31,
					})
			},
			"UPDATE users SET age = ?, name = ? WHERE id IN (?, ?, ?)",
			[]interface{}{31, "Joe", 1, 2, 3},
		},
		{
			"Update_where_null",
			func() *query.UpdateBuilder {
				return query.NewUpdateBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache(100)).
					Table("users").
					WhereNull("name").
					Update(map[string]interface{}{
						"name": "Joe",
						"age":  31,
					})
			},
			"UPDATE users SET age = ?, name = ? WHERE name IS NULL",
			[]interface{}{31, "Joe"},
		},
		{
			"Update_where_column",
			func() *query.UpdateBuilder {
				return query.NewUpdateBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache(100)).
					Table("users").
					WhereColumn([]string{"name", "note"}, "name", "=", "note").
					Update(map[string]interface{}{
						"name": "Joe",
						"age":  31,
					})
			},
			"UPDATE users SET age = ?, name = ? WHERE name = note",
			[]interface{}{31, "Joe"},
		},
		{
			"Update_where_or_columns",
			func() *query.UpdateBuilder {
				return query.NewUpdateBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache(100)).
					Table("users").
					OrWhereColumns([]string{"name", "nick_name", "memo", "note"}, [][]string{{"name", "=", "nick_name"}, {"memo", "=", "note"}}).
					Update(map[string]interface{}{
						"name": "Joe",
						"age":  31,
					})
			},
			"UPDATE users SET age = ?, name = ? WHERE name = nick_name OR memo = note",
			[]interface{}{31, "Joe"},
		},
		{
			"Update_JOINS",
			func() *query.UpdateBuilder {
				return query.NewUpdateBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache(100)).
					Table("users").
					Join("profiles", "users.id", "=", "profiles.user_id").
					Where("age", ">", 18).
					Update(map[string]interface{}{
						"name": "Joe",
						"age":  31,
					})
			},
			"UPDATE users INNER JOIN profiles ON users.id = profiles.user_id SET age = ?, name = ? WHERE age > ?",
			[]interface{}{31, "Joe", 18},
		},
		{
			"Update_orderBy",
			func() *query.UpdateBuilder {
				return query.NewUpdateBuilder(&db.MySQLQueryBuilder{}, cache.NewAsyncQueryCache(100)).
					Table("users").
					OrderBy("name", "ASC").
					Update(map[string]interface{}{
						"name": "Joe",
						"age":  31,
					})
			},
			"UPDATE users SET age = ?, name = ? ORDER BY name ASC",
			[]interface{}{31, "Joe"},
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
