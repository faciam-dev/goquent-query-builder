package api_test

import (
	"strings"
	"testing"

	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/sliceutils"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/pkg/api"
)

func TestSelectDebugApiRawSqlTest(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() *api.SelectBuilder
		expectedQuery string
	}{
		{
			"Complex_Query_With_Union",
			func() *api.SelectBuilder {
				dbStrategy := db.NewMySQLQueryBuilder()

				blankCache := cache.NewBlankQueryCache()

				usq := api.NewSelectBuilder(dbStrategy, blankCache).
					Table("users").
					Select("id").
					Join("profiles", "users.id", "=", "profiles.user_id").
					Where("profiles.age", ">", 18)

				return api.NewSelectBuilder(dbStrategy, blankCache).
					Table("users").
					Select("id", "users.name AS name").
					Join("profiles", "users.id", "=", "profiles.user_id").
					Where("area", "=", "Jakarta").
					WhereBetween("profiles.age", 18, 30).
					OrderBy("users.name", "ASC").
					OrderBy("profiles.age", "DESC").
					GroupBy("users.id").
					Having("COUNT(profiles.id)", ">", 1).
					WhereIn("users.id", usq).
					Union(
						api.NewSelectBuilder(dbStrategy, blankCache).
							Table("users").
							Select("id", "users.name AS name").
							Join("profiles", "users.id", "=", "profiles.user_id").
							Where("area", "=", "Jakarta").
							WhereBetween("profiles.age", 18, 30).
							OrderBy("users.name", "ASC").
							OrderBy("profiles.age", "DESC").
							GroupBy("users.id").
							Having("COUNT(profiles.id)", ">", 1).
							WhereIn("users.id", usq),
					)

			},
			"SELECT id, users.name AS name FROM users INNER JOIN profiles ON users.id = profiles.user_id WHERE area = 'Jakarta' AND profiles.age BETWEEN 18 AND 30 AND users.id IN (SELECT id FROM users INNER JOIN profiles ON users.id = profiles.user_id WHERE profiles.age > 18) GROUP BY users.id HAVING COUNT(profiles.id) > 1 ORDER BY users.name ASC, profiles.age DESC UNION SELECT id, users.name AS name FROM users INNER JOIN profiles ON users.id = profiles.user_id WHERE area = 'Jakarta' AND profiles.age BETWEEN 18 AND 30 AND users.id IN (SELECT id FROM users INNER JOIN profiles ON users.id = profiles.user_id WHERE profiles.age > 18) GROUP BY users.id HAVING COUNT(profiles.id) > 1 ORDER BY users.name ASC, profiles.age DESC",
		},
		{
			"Complex_Query",
			func() *api.SelectBuilder {
				dbStrategy := db.NewMySQLQueryBuilder()

				blankCache := cache.NewBlankQueryCache()

				usq := api.NewSelectBuilder(dbStrategy, blankCache).
					Table("users").
					Select("id").
					Join("profiles", "users.id", "=", "profiles.user_id").
					Where("profiles.age", ">", 18)

				return api.NewSelectBuilder(dbStrategy, blankCache).
					Table("users").
					Select("id", "users.name AS name").
					Join("profiles", "users.id", "=", "profiles.user_id").
					Where("area", "=", "Jakarta").
					WhereBetween("profiles.age", 18, 30).
					OrderBy("users.name", "ASC").
					OrderBy("profiles.age", "DESC").
					GroupBy("users.id").
					Having("COUNT(profiles.id)", ">", 1).
					WhereIn("users.id", usq)

			},
			"SELECT id, users.name AS name FROM users INNER JOIN profiles ON users.id = profiles.user_id WHERE area = 'Jakarta' AND profiles.age BETWEEN 18 AND 30 AND users.id IN (SELECT id FROM users INNER JOIN profiles ON users.id = profiles.user_id WHERE profiles.age > 18) GROUP BY users.id HAVING COUNT(profiles.id) > 1 ORDER BY users.name ASC, profiles.age DESC",
		},
		{
			"Normal_Query",
			func() *api.SelectBuilder {
				dbStrategy := db.NewMySQLQueryBuilder()

				blankCache := cache.NewBlankQueryCache()

				return api.NewSelectBuilder(dbStrategy, blankCache).
					Table("users").
					Select("id", "users.name AS name").
					Join("profiles", "users.id", "=", "profiles.user_id").
					Where("profiles.age", ">", 18).
					OrderBy("users.name", "ASC")

			},
			"SELECT id, users.name AS name FROM users INNER JOIN profiles ON users.id = profiles.user_id WHERE profiles.age > 18 ORDER BY users.name ASC",
		},
		{
			"Normal_Query_With_WhereGroup",
			func() *api.SelectBuilder {
				dbStrategy := db.NewMySQLQueryBuilder()

				blankCache := cache.NewBlankQueryCache()

				return api.NewSelectBuilder(dbStrategy, blankCache).
					Table("users").
					Select("id", "users.name AS name").
					Join("profiles", "users.id", "=", "profiles.user_id").
					WhereGroup(func(w *api.WhereSelectBuilder) {
						w.Where("profiles.age", ">", 18).Where("profiles.age", "<", 30)
					}).OrderBy("users.name", "ASC")
			},
			"SELECT id, users.name AS name FROM users INNER JOIN profiles ON users.id = profiles.user_id WHERE (profiles.age > 18 AND profiles.age < 30) ORDER BY users.name ASC",
		},
		{
			"Simple_Query",
			func() *api.SelectBuilder {
				dbStrategy := db.NewMySQLQueryBuilder()

				blankCache := cache.NewBlankQueryCache()

				return api.NewSelectBuilder(dbStrategy, blankCache).
					Table("users").
					Select("id", "users.name AS name")

			},
			"SELECT id, users.name AS name FROM users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			builder := tt.setup()
			query, _ := builder.RawSql()

			if query != tt.expectedQuery {
				t.Errorf("expected '%s' but got '%s'", tt.expectedQuery, query)
			}

		})
	}
}

func TestInsertDebugApiRawSqlTest(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() *api.InsertBuilder
		expectedQuery string
	}{
		{
			"Complex_Query",
			func() *api.InsertBuilder {
				dbStrategy := db.NewMySQLQueryBuilder()

				blankCache := cache.NewBlankQueryCache()

				return api.NewInsertBuilder(dbStrategy, blankCache).
					Table("users").
					Insert(map[string]interface{}{
						"name": "Joe",
						"age":  31,
					})
			},
			"INSERT INTO users (age, name) VALUES (31, 'Joe')",
		},
		{
			"InsertUsing",
			func() *api.InsertBuilder {
				dbStrategy := db.NewMySQLQueryBuilder()

				blankCache := cache.NewBlankQueryCache()

				return api.NewInsertBuilder(dbStrategy, blankCache).
					Table("users").
					InsertUsing([]string{"name", "age"}, api.NewSelectBuilder(dbStrategy, blankCache).
						Table("profiles").
						Select("name", "age").
						Where("age", ">", 18)).
					Insert(map[string]interface{}{
						"name": "Joe",
						"age":  31,
					})
			},
			"INSERT INTO users (name, age) SELECT name, age FROM profiles WHERE age > 18",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			builder := tt.setup()
			query, _ := builder.RawSql()

			if query != tt.expectedQuery {
				t.Errorf("expected '%s' but got '%s'", tt.expectedQuery, query)
			}

		})
	}
}

func TestUpdateDebugApiRawSqlTest(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() *api.UpdateBuilder
		expectedQuery string
	}{
		{
			"Complex_Query",
			func() *api.UpdateBuilder {
				dbStrategy := db.NewMySQLQueryBuilder()

				blankCache := cache.NewBlankQueryCache()

				return api.NewUpdateBuilder(dbStrategy, blankCache).
					Table("users").
					Join("profiles", "users.id", "=", "profiles.user_id").
					Where("age", ">", 18).
					Update(map[string]interface{}{
						"name": "Joe",
						"age":  31,
					})
			},
			"UPDATE users INNER JOIN profiles ON users.id = profiles.user_id SET age = 31, name = 'Joe' WHERE age > 18",
		},
		{
			"Update_ORDER_BY",
			func() *api.UpdateBuilder {
				dbStrategy := db.NewMySQLQueryBuilder()

				blankCache := cache.NewBlankQueryCache()

				return api.NewUpdateBuilder(dbStrategy, blankCache).
					Table("users").
					OrderBy("name", "ASC").
					Update(map[string]interface{}{
						"name": "Joe",
						"age":  31,
					})
			},
			"UPDATE users SET age = 31, name = 'Joe' ORDER BY name ASC",
		},
		{
			"Update_Where_Date",
			func() *api.UpdateBuilder {
				dbStrategy := db.NewMySQLQueryBuilder()

				blankCache := cache.NewBlankQueryCache()

				return api.NewUpdateBuilder(dbStrategy, blankCache).
					Table("users").
					WhereDate("created_at", "=", "2021-01-01").
					Update(map[string]interface{}{
						"name": "Joe",
						"age":  31,
					})
			},
			"UPDATE users SET age = 31, name = 'Joe' WHERE DATE(created_at) = '2021-01-01'",
		},
		{
			"Update_Where_Between_Columns",
			func() *api.UpdateBuilder {
				dbStrategy := db.NewMySQLQueryBuilder()

				blankCache := cache.NewBlankQueryCache()

				return api.NewUpdateBuilder(dbStrategy, blankCache).
					Table("users").
					WhereBetweenColumns([]string{"age", "min_age", "max_age"}, "age", "min_age", "max_age").
					Update(map[string]interface{}{
						"name": "Joe",
						"age":  31,
					})
			},
			"UPDATE users SET age = 31, name = 'Joe' WHERE age BETWEEN min_age AND max_age",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			builder := tt.setup()
			query, _ := builder.RawSql()

			if query != tt.expectedQuery {
				t.Errorf("expected '%s' but got '%s'", tt.expectedQuery, query)
			}

		})
	}
}

func TestDeleteDebugApiRawSqlTest(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() *api.DeleteBuilder
		expectedQuery string
	}{
		{
			"Complex_Query",
			func() *api.DeleteBuilder {
				dbStrategy := db.NewMySQLQueryBuilder()

				blankCache := cache.NewBlankQueryCache()

				return api.NewDeleteBuilder(dbStrategy, blankCache).
					Table("users").
					Join("profiles", "users.id", "=", "profiles.user_id").
					Where("age", ">", 18).
					Delete()
			},
			"DELETE users FROM users INNER JOIN profiles ON users.id = profiles.user_id WHERE age > 18",
		},
		{
			"Delete_Where_Between",
			func() *api.DeleteBuilder {
				dbStrategy := db.NewMySQLQueryBuilder()

				blankCache := cache.NewBlankQueryCache()

				return api.NewDeleteBuilder(dbStrategy, blankCache).
					Table("users").
					WhereNotBetween("age", 18, 30).
					Delete()
			},
			"DELETE FROM users WHERE age NOT BETWEEN 18 AND 30",
		},
		{
			"Delete_Where_Between_Columns",
			func() *api.DeleteBuilder {
				dbStrategy := db.NewMySQLQueryBuilder()

				blankCache := cache.NewBlankQueryCache()

				return api.NewDeleteBuilder(dbStrategy, blankCache).
					Table("users").
					WhereBetweenColumns([]string{"created_at", "updated_at", "deleted_at"}, "created_at", "updated_at", "deleted_at").
					Delete()
			},
			"DELETE FROM users WHERE created_at BETWEEN updated_at AND deleted_at",
		},
		{
			"Delete_Where_Columns",
			func() *api.DeleteBuilder {
				dbStrategy := db.NewMySQLQueryBuilder()

				blankCache := cache.NewBlankQueryCache()

				return api.NewDeleteBuilder(dbStrategy, blankCache).
					Table("users").
					WhereColumns([]string{"name", "age"},
						[][]string{
							{"name", "age"},
						}).
					Delete()
			},
			"DELETE FROM users WHERE name = age",
		},
		{
			"Delete_Where_Columns_With_WhereGroup",
			func() *api.DeleteBuilder {
				dbStrategy := db.NewMySQLQueryBuilder()

				blankCache := cache.NewBlankQueryCache()

				return api.NewDeleteBuilder(dbStrategy, blankCache).
					Table("users").
					WhereGroup(func(w *api.WhereDeleteBuilder) {
						w.Where("name", "=", "Joe").Where("age", "=", 31)
					}).
					Delete()
			},
			"DELETE FROM users WHERE (name = 'Joe' AND age = 31)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			builder := tt.setup()
			query, _ := builder.RawSql()

			if query != tt.expectedQuery {
				t.Errorf("expected '%s' but got '%s'", tt.expectedQuery, query)
			}

		})
	}
}

func TestDebugDumpTest(t *testing.T) {
	tests := []struct {
		name           string
		setup          func() *api.SelectBuilder
		expectedQuery  string
		expectedValues []interface{}
	}{
		{
			"Complex_Query_With_Union",
			func() *api.SelectBuilder {
				dbStrategy := db.NewMySQLQueryBuilder()

				blankCache := cache.NewBlankQueryCache()

				usq := api.NewSelectBuilder(dbStrategy, blankCache).
					Table("users").
					Select("id").
					Join("profiles", "users.id", "=", "profiles.user_id").
					Where("profiles.age", ">", 18)

				return api.NewSelectBuilder(dbStrategy, blankCache).
					Table("users").
					Select("id", "users.name AS name").
					Join("profiles", "users.id", "=", "profiles.user_id").
					Where("area", "=", "Jakarta").
					WhereBetween("profiles.age", 18, 30).
					OrderBy("users.name", "ASC").
					OrderBy("profiles.age", "DESC").
					GroupBy("users.id").
					Having("COUNT(profiles.id)", ">", 1).
					WhereIn("users.id", usq).
					Union(
						api.NewSelectBuilder(dbStrategy, blankCache).
							Table("users").
							Select("id", "users.name AS name").
							Join("profiles", "users.id", "=", "profiles.user_id").
							Where("area", "=", "Jakarta").
							WhereBetween("profiles.age", 18, 30).
							OrderBy("users.name", "ASC").
							OrderBy("profiles.age", "DESC").
							GroupBy("users.id").
							Having("COUNT(profiles.id)", ">", 1).
							WhereIn("users.id", usq),
					)

			},
			"SELECT id, users.name AS name FROM users INNER JOIN profiles ON users.id = profiles.user_id WHERE area = ? AND profiles.age BETWEEN ? AND ? AND users.id IN (SELECT id FROM users INNER JOIN profiles ON users.id = profiles.user_id WHERE profiles.age > ?) GROUP BY users.id HAVING COUNT(profiles.id) > ? ORDER BY users.name ASC, profiles.age DESC UNION SELECT id, users.name AS name FROM users INNER JOIN profiles ON users.id = profiles.user_id WHERE area = ? AND profiles.age BETWEEN ? AND ? AND users.id IN (SELECT id FROM users INNER JOIN profiles ON users.id = profiles.user_id WHERE profiles.age > ?) GROUP BY users.id HAVING COUNT(profiles.id) > ? ORDER BY users.name ASC, profiles.age DESC",
			[]interface{}{
				"Jakarta", 18, 30, 18, 1, "Jakarta", 18, 30, 18, 1,
			},
		},
		{
			"Complex_Query",
			func() *api.SelectBuilder {
				dbStrategy := db.NewMySQLQueryBuilder()

				blankCache := cache.NewBlankQueryCache()

				usq := api.NewSelectBuilder(dbStrategy, blankCache).
					Table("users").
					Select("id").
					Join("profiles", "users.id", "=", "profiles.user_id").
					Where("profiles.age", ">", 18)

				return api.NewSelectBuilder(dbStrategy, blankCache).
					Table("users").
					Select("id", "users.name AS name").
					Join("profiles", "users.id", "=", "profiles.user_id").
					Where("area", "=", "Jakarta").
					WhereBetween("profiles.age", 18, 30).
					OrderBy("users.name", "ASC").
					OrderBy("profiles.age", "DESC").
					GroupBy("users.id").
					Having("COUNT(profiles.id)", ">", 1).
					WhereIn("users.id", usq)

			},
			"SELECT id, users.name AS name FROM users INNER JOIN profiles ON users.id = profiles.user_id WHERE area = ? AND profiles.age BETWEEN ? AND ? AND users.id IN (SELECT id FROM users INNER JOIN profiles ON users.id = profiles.user_id WHERE profiles.age > ?) GROUP BY users.id HAVING COUNT(profiles.id) > ? ORDER BY users.name ASC, profiles.age DESC",
			[]interface{}{
				"Jakarta", 18, 30, 18, 1,
			},
		},
		{
			"Simple_Query",
			func() *api.SelectBuilder {
				dbStrategy := db.NewMySQLQueryBuilder()

				blankCache := cache.NewBlankQueryCache()

				return api.NewSelectBuilder(dbStrategy, blankCache).
					Table("users").
					Select("id", "users.name AS name")

			},
			"SELECT id, users.name AS name FROM users",
			[]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			builder := tt.setup()
			query, values, _ := builder.Dump()

			if query != tt.expectedQuery {
				t.Errorf("expected '%s' but got '%s'", tt.expectedQuery, query)
			}

			if len(values) != len(tt.expectedValues) {
				t.Errorf("expected values %v but got %v", tt.expectedValues, values)
			}

			convertedValues := sliceutils.ToInterfaceSlice(values)
			for i := range convertedValues {
				if values[i] != tt.expectedValues[i] {
					t.Errorf("expected value %v at index %d but got %v", tt.expectedValues[i], i, values[i])
				}
			}

			builder.Where("debug", "=", 1)

			query, values, _ = builder.Build()
			t.Log(query, values)

			tt.expectedValues = append(tt.expectedValues, 1)

			if !strings.Contains(query, "debug = ?") {
				t.Errorf("expected '%s' but got '%s'", tt.expectedQuery, query)
			}

			if len(values) != len(tt.expectedValues) {
				t.Errorf("expected values %v but got %v", tt.expectedValues, values)
			}

			convertedValues = sliceutils.ToInterfaceSlice(values)
			for i := range convertedValues {
				if values[i] != tt.expectedValues[i] {
					t.Errorf("expected value %v at index %d but got %v", tt.expectedValues[i], i, values[i])
				}
			}

		})
	}
}
