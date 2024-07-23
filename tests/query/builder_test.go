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
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Select("id", "name")
			},
			"SELECT id, name FROM ",
			nil,
		},
		{
			"SelectRaw",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).SelectRaw("COUNT(*) as total")
			},
			"SELECT COUNT(*) as total FROM ",
			nil,
		},
		{
			"Count",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Count()
			},
			"SELECT COUNT(*) FROM ",
			nil,
		},
		{
			"Max",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Max("price")
			},
			"SELECT MAX(price) FROM ",
			nil,
		},
		{
			"Min",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Min("price")
			},
			"SELECT MIN(price) FROM ",
			nil,
		},
		{
			"Sum",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Sum("price")
			},
			"SELECT SUM(price) FROM ",
			nil,
		},
		{
			"Avg",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Avg("price")
			},
			"SELECT AVG(price) FROM ",
			nil,
		},
		{
			"SelectRaw_With_Value",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).SelectRaw("price * ? as price_with_tax", 1.0825)
			},
			"SELECT price * ? as price_with_tax FROM ",
			[]interface{}{1.0825},
		},
		{
			"From",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Table("users")
			},
			"SELECT  FROM users",
			nil,
		},
		{
			"Inner_Join",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Join("orders", "users.id", "=", "orders.user_id")
			},
			"SELECT orders.*, .* FROM  INNER JOIN orders ON users.id = orders.user_id",
			nil,
		},
		{
			"Left_Join",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).LeftJoin("orders", "users.id", "=", "orders.user_id")
			},
			"SELECT orders.*, .* FROM  LEFT JOIN orders ON users.id = orders.user_id",
			nil,
		},
		{
			"Right_Join",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).RightJoin("orders", "users.id", "=", "orders.user_id")
			},
			"SELECT orders.*, .* FROM  RIGHT JOIN orders ON users.id = orders.user_id",
			nil,
		},
		{
			"Cross_Join",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).CrossJoin("orders")
			},
			"SELECT orders.*, .* FROM  CROSS JOIN orders",
			nil,
		},
		{
			"Join_and_Join",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Join("orders", "users.id", "=", "orders.user_id").Join("products", "orders.product_id", "=", "products.id")
			},
			"SELECT orders.*, .*, products.* FROM  INNER JOIN orders ON users.id = orders.user_id INNER JOIN products ON orders.product_id = products.id",
			nil,
		},
		{
			"JoinQuery",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).JoinQuery("users", func(b *query.JoinClauseBuilder) *query.JoinClauseBuilder {
					return b.On("users.id", "=", "profiles.user_id").OrOn("users.id", "=", "profiles.alter_user_id").Where("profiles.age", ">", 18)
				})
			},
			"SELECT users.* FROM  INNER JOIN users ON users.id = profiles.user_id OR users.id = profiles.alter_user_id AND profiles.age > ?",
			[]interface{}{18},
		},
		{
			"LeftJoinQuery",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).LeftJoinQuery("users", func(b *query.JoinClauseBuilder) *query.JoinClauseBuilder {
					return b.On("users.id", "=", "profiles.user_id").OrOn("users.id", "=", "profiles.alter_user_id").Where("profiles.age", ">", 18)
				})
			},
			"SELECT users.* FROM  LEFT JOIN users ON users.id = profiles.user_id OR users.id = profiles.alter_user_id AND profiles.age > ?",
			[]interface{}{18},
		},
		{
			"RightJoinQuery",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).RightJoinQuery("users", func(b *query.JoinClauseBuilder) *query.JoinClauseBuilder {
					return b.On("users.id", "=", "profiles.user_id").OrOn("users.id", "=", "profiles.alter_user_id").Where("profiles.age", ">", 18)
				})
			},
			"SELECT users.* FROM  RIGHT JOIN users ON users.id = profiles.user_id OR users.id = profiles.alter_user_id AND profiles.age > ?",
			[]interface{}{18},
		},
		{
			"JoinSub",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).JoinSub(query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Select("id").Table("profiles").Where("age", ">", 18), "profiles", "users.id", "=", "profiles.user_id")
			},
			"SELECT profiles.*, .* FROM  INNER JOIN (SELECT id FROM profiles WHERE age > ?) AS profiles ON users.id = profiles.user_id",
			[]interface{}{18},
		},
		{
			"OrderBy",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).OrderBy("name", "asc")
			},
			"SELECT  FROM  ORDER BY name ASC",
			nil,
		},
		{
			"OrderByDesc_And_ReOrder",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).OrderBy("name", "asc").ReOrder().OrderBy("name", "desc")
			},
			"SELECT  FROM  ORDER BY name DESC",
			nil,
		},
		{
			"OrderByRaw",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).OrderByRaw("RAND()")
			},
			"SELECT  FROM  ORDER BY RAND()",
			nil,
		},
		{
			"GroupBy",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).GroupBy("name", "age")
			},
			"SELECT  FROM  GROUP BY name, age",
			nil,
		},
		{
			"GroupBy_Having",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).GroupBy("name", "age").Having("age", ">", 18)
			},
			"SELECT  FROM  GROUP BY name, age HAVING age > ?",
			[]interface{}{18},
		},
		{
			"GroupBy_Having_OR",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).GroupBy("name", "age").Having("age", ">", 18).OrHaving("name", "=", "John")
			},
			"SELECT  FROM  GROUP BY name, age HAVING age > ? OR name = ?",
			[]interface{}{18, "John"},
		},
		{
			"GroupBy_Having_Raw",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).GroupBy("name", "age").HavingRaw("age > 18")
			},
			"SELECT  FROM  GROUP BY name, age HAVING age > 18",
			nil,
		},
		{
			"GroupBy_HavingRaw_OR",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).GroupBy("name", "age").HavingRaw("age > 18").OrHavingRaw("name = 'John'")
			},
			"SELECT  FROM  GROUP BY name, age HAVING age > 18 OR name = 'John'",
			nil,
		},
		{
			"Limit",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Limit(10)
			},
			"SELECT  FROM  LIMIT 10",
			nil,
		},
		{
			"Offset",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Offset(10)
			},
			"SELECT  FROM  OFFSET 10",
			nil,
		},
		{
			"Limit_And_Offset",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Limit(10).Offset(5)
			},
			"SELECT  FROM  LIMIT 10 OFFSET 5",
			nil,
		},
		{
			"Lock FOR UPDATE",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).LockForUpdate()
			},
			"SELECT  FROM  FOR UPDATE",
			nil,
		},
		{
			"Lock_Shared",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).SharedLock()
			},
			"SELECT  FROM  LOCK IN SHARE MODE",
			nil,
		},
		{
			"Complex_Query",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).
					Select("id", "name").
					Table("users").
					Join("profiles", "users.id", "=", "profiles.user_id").
					Where("age", ">", 18).
					OrderBy("name", "ASC")
			},
			"SELECT id, name FROM users INNER JOIN profiles ON users.id = profiles.user_id WHERE age > ? ORDER BY name ASC",
			[]interface{}{18},
		},
		{
			"Complex_Query_With_Subquery",
			func() *query.Builder {
				sq := query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Select("id").Table("users").Where("name", "=", "John")
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).
					SelectRaw("id, name, profiles.point * ? AS profiles_point", 1.05).
					Table("users").
					Join("profiles", "users.id", "=", "profiles.user_id").
					Where("status", "=", "active").
					WhereQuery("user_id", "IN", sq).
					Where("age", ">", 18).
					OrderBy("name", "ASC")
			},
			"SELECT id, name, profiles.point * ? AS profiles_point FROM users INNER JOIN profiles ON users.id = profiles.user_id WHERE status = ? AND user_id IN (SELECT id FROM users WHERE name = ?) AND age > ? ORDER BY name ASC",
			[]interface{}{1.05, "active", "John", 18},
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

func TestWhereSelectBuilder(t *testing.T) {
	tests := []struct {
		name           string
		setup          func() *query.Builder
		expectedQuery  string
		expectedValues []interface{}
	}{
		{
			"Where",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("age", ">", 18)
			},
			"SELECT  FROM  WHERE age > ?",
			[]interface{}{18},
		},
		{
			"OrWhere",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("email", "LIKE", "%@gmail.com%").OrWhere("email", "LIKE", "%@yahoo.com%").OrWhere("age", ">", 18)
			},
			"SELECT  FROM  WHERE email LIKE ? OR email LIKE ? OR age > ?",
			[]interface{}{"%@gmail.com%", "%@yahoo.com%", 18},
		},
		{
			"WhereRaw",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereRaw("age > ?", 18)
			},
			"SELECT  FROM  WHERE age > ?",
			[]interface{}{18},
		},
		{
			"OrWhereRaw",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereRaw("age > ?", 18).OrWhereRaw("name= ?", "John")
			},
			"SELECT  FROM  WHERE age > ? OR name= ?",
			[]interface{}{18, "John"},
		},
		{
			"WhereQuery",
			func() *query.Builder {
				sq := query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Select("id").Table("users").Where("name", "=", "John")
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereQuery("user_id", "IN", sq).Where("city", "=", "New York")
			},
			"SELECT  FROM  WHERE user_id IN (SELECT id FROM users WHERE name = ?) AND city = ?",
			[]interface{}{"John", "New York"},
		},
		{
			"OrWhereQuery",
			func() *query.Builder {
				sq := query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Select("id").Table("users").Where("name", "=", "John")
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("city", "=", "New York").OrWhereQuery("user_id", "IN", sq)
			},
			"SELECT  FROM  WHERE city = ? OR user_id IN (SELECT id FROM users WHERE name = ?)",
			[]interface{}{"New York", "John"},
		},
		{
			"WhereGroup",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).
					WhereGroup(func(b *query.WhereBuilder) *query.WhereBuilder {
						return b.Where("age", ">", 18).Where("name", "=", "John")
					})
			},
			"SELECT  FROM  WHERE (age > ? AND name = ?)",
			[]interface{}{18, "John"},
		},
		{
			"WhereGroup_And",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).
					WhereGroup(func(b *query.WhereBuilder) *query.WhereBuilder {
						return b.Where("age", ">", 18).Where("name", "=", "John")
					}).Where("city", "=", "New York")
			},
			"SELECT  FROM  WHERE (age > ? AND name = ?) AND city = ?",
			[]interface{}{18, "John", "New York"},
		},
		{
			"WhereGroup_And_2",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).
					Where("city", "=", "New York").
					WhereGroup(func(b *query.WhereBuilder) *query.WhereBuilder {
						return b.Where("age", ">", 18).Where("name", "=", "John")
					})
			},
			"SELECT  FROM  WHERE city = ? AND (age > ? AND name = ?)",
			[]interface{}{"New York", 18, "John"},
		},
		{
			"WhereGroup_Or",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).
					WhereGroup(func(b *query.WhereBuilder) *query.WhereBuilder {
						return b.Where("age", ">", 18).Where("name", "=", "John")
					}).OrWhere("city", "=", "New York")
			},
			"SELECT  FROM  WHERE (age > ? AND name = ?) OR city = ?",
			[]interface{}{18, "John", "New York"},
		},
		{
			"WhereGroup_Or_2",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).
					Where("city", "=", "New York").
					OrWhereGroup(func(b *query.WhereBuilder) *query.WhereBuilder {
						return b.Where("age", ">", 18).Where("name", "=", "John")
					})
			},
			"SELECT  FROM  WHERE city = ? OR (age > ? AND name = ?)",
			[]interface{}{"New York", 18, "John"},
		},
		{
			"WhereNot",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).
					WhereNot(func(b *query.WhereBuilder) *query.WhereBuilder {
						return b.Where("age", ">", 18).Where("name", "=", "John")
					})
			},
			"SELECT  FROM  WHERE NOT (age > ? AND name = ?)",
			[]interface{}{18, "John"},
		},
		{
			"OrWhereNot",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).
					Where("city", "=", "New York").
					OrWhereNot(func(b *query.WhereBuilder) *query.WhereBuilder {
						return b.Where("age", ">", 18).Where("name", "=", "John")
					})
			},
			"SELECT  FROM  WHERE city = ? OR NOT (age > ? AND name = ?)",
			[]interface{}{"New York", 18, "John"},
		},
		{
			"WhereAll",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).
					Where("age", ">", 18).
					WhereAll([]string{"name", "city"}, "LIKE", "%test%")
			},
			"SELECT  FROM  WHERE age > ? AND (name LIKE ? AND city LIKE ?)",
			[]interface{}{18, "%test%", "%test%"},
		},
		{
			"WhereAny",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).
					Where("age", ">", 18).
					WhereAny([]string{"name", "city"}, "LIKE", "%test%")
			},
			"SELECT  FROM  WHERE age > ? AND (name LIKE ? OR city LIKE ?)",
			[]interface{}{18, "%test%", "%test%"},
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
