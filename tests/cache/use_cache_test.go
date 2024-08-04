package cache_test

import (
	"testing"
	"time"

	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/db/mysql"
	"github.com/faciam-dev/goquent-query-builder/internal/query"
)

func TestUseCacheTest(t *testing.T) {
	tests := []struct {
		name           string
		setup          func() *query.Builder
		expectedQuery  string
		expectedValues []interface{}
	}{
		{
			"Complex_Query",
			func() *query.Builder {
				return query.NewBuilder(mysql.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).
					Select("id", "name").
					Table("users").
					Join("profiles", "users.id", "=", "profiles.user_id").
					Where("age", ">", 18).
					OrderBy("name", "ASC")
			},
			"SELECT `id`, `name` FROM `users` INNER JOIN `profiles` ON `users`.`id` = `profiles`.`user_id` WHERE `age` > ? ORDER BY `name` ASC",
			[]interface{}{18},
		},
		{
			"Complex_Query_With_Subquery",
			func() *query.Builder {
				sq := query.NewBuilder(mysql.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Select("id").Table("users").Where("name", "=", "John")
				return query.NewBuilder(mysql.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).
					SelectRaw("`id`, `name`, `profiles`.`point` * ? as `profiles_point`", 1.05).
					Table("users").
					Join("profiles", "users.id", "=", "profiles.user_id").
					Where("status", "=", "active").
					WhereSubQuery("user_id", "IN", sq).
					Where("age", ">", 18).
					OrderBy("name", "ASC")
			},
			"SELECT `id`, `name`, `profiles`.`point` * ? as `profiles_point` FROM `users` INNER JOIN `profiles` ON `users`.`id` = `profiles`.`user_id` WHERE `status` = ? AND `user_id` IN (SELECT `id` FROM `users` WHERE `name` = ?) AND `age` > ? ORDER BY `name` ASC",
			[]interface{}{1.05, "active", "John", 18},
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
			time.Sleep(500 * time.Millisecond)

			query, values, _ = builder.Build()
			if query != tt.expectedQuery {
				t.Errorf("(cache error) expected '%s' but got '%s'", tt.expectedQuery, query)
			}

			if len(values) != len(tt.expectedValues) {
				t.Errorf("(cache error) expected values %v but got %v", tt.expectedValues, values)
			}

			for i := range values {
				if values[i] != tt.expectedValues[i] {
					t.Errorf("(cache error) expected value %v at index %d but got %v", tt.expectedValues[i], i, values[i])
				}
			}

		})
	}
}
