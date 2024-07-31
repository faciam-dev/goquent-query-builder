package api_test

import (
	"strings"
	"testing"

	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/sliceutils"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/pkg/api"
)

func TestDebugApiRawSqlTest(t *testing.T) {
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
