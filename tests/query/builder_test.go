package query_test

import (
	"testing"

	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/common/sliceutils"
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
			"Count_Columns",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Count("id")
			},
			"SELECT COUNT(id) FROM ",
			nil,
		},
		{
			"Count_Distinct",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Count("id").Distinct("id")
			},
			"SELECT COUNT(DISTINCT id) FROM ",
			nil,
		},
		{
			"Count_Distinct_Columns",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Distinct("id", "name").Count("id", "name")
			},
			"SELECT COUNT(DISTINCT id), COUNT(DISTINCT name) FROM ",
			nil,
		},
		{
			"Distincts",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Distinct("id", "name")
			},
			"SELECT DISTINCT id, name FROM ",
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
			"Lateral_Join",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).JoinLateral(query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Select("id").Table("profiles").Where("age", ">", 18), "profiles")
			},
			"SELECT profiles.*, .* FROM  ,LATERAL(SELECT id FROM profiles WHERE age > ?) AS profiles",
			[]interface{}{18},
		},
		{
			"LeftLateral_Join",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).LeftJoinLateral(query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Select("id").Table("profiles").Where("age", ">", 18), "profiles")
			},
			"SELECT profiles.*, .* FROM  ,LEFT LATERAL(SELECT id FROM profiles WHERE age > ?) AS profiles",
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
					WhereGroup(func(b *query.WhereBuilder[query.Builder]) *query.WhereBuilder[query.Builder] {
						return &b.Where("age", ">", 18).Where("name", "=", "John").WhereBuilder
					})
			},
			"SELECT  FROM  WHERE (age > ? AND name = ?)",
			[]interface{}{18, "John"},
		},
		{
			"WhereGroup_And",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).
					WhereGroup(func(b *query.WhereBuilder[query.Builder]) *query.WhereBuilder[query.Builder] {
						return &b.Where("age", ">", 18).Where("name", "=", "John").WhereBuilder
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
					WhereGroup(func(b *query.WhereBuilder[query.Builder]) *query.WhereBuilder[query.Builder] {
						return &b.Where("age", ">", 18).Where("name", "=", "John").WhereBuilder
					})
			},
			"SELECT  FROM  WHERE city = ? AND (age > ? AND name = ?)",
			[]interface{}{"New York", 18, "John"},
		},
		{
			"WhereGroup_Or",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).
					WhereGroup(func(b *query.WhereBuilder[query.Builder]) *query.WhereBuilder[query.Builder] {
						return &b.Where("age", ">", 18).Where("name", "=", "John").WhereBuilder
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
					OrWhereGroup(func(b *query.WhereBuilder[query.Builder]) *query.WhereBuilder[query.Builder] {
						return &b.Where("age", ">", 18).Where("name", "=", "John").WhereBuilder
					})
			},
			"SELECT  FROM  WHERE city = ? OR (age > ? AND name = ?)",
			[]interface{}{"New York", 18, "John"},
		},
		{
			"WhereNot",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).
					WhereNot(func(b *query.WhereBuilder[query.Builder]) *query.WhereBuilder[query.Builder] {
						return &b.Where("age", ">", 18).Where("name", "=", "John").WhereBuilder
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
					OrWhereNot(func(b *query.WhereBuilder[query.Builder]) *query.WhereBuilder[query.Builder] {
						return &b.Where("age", ">", 18).Where("name", "=", "John").WhereBuilder
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
		{
			"WhereIn",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereIn("id", []int64{1, 2, 3})
			},
			"SELECT  FROM  WHERE id IN (?, ?, ?)",
			[]interface{}{1, 2, 3},
		},
		{
			"WhereIn (Subquery)",
			func() *query.Builder {
				sq := query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Select("id").Table("users").Where("name", "=", "John")
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereIn("id", sq)
			},
			"SELECT  FROM  WHERE id IN (SELECT id FROM users WHERE name = ?)",
			[]interface{}{"John"},
		},
		{
			"WhereNotIn",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereNotIn("id", []int64{1, 2, 3})
			},
			"SELECT  FROM  WHERE id NOT IN (?, ?, ?)",
			[]interface{}{1, 2, 3},
		},
		{
			"WhereNotIn (Subquery)",
			func() *query.Builder {
				sq := query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Select("id").Table("users").Where("name", "=", "John")
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereNotIn("id", sq)
			},
			"SELECT  FROM  WHERE id NOT IN (SELECT id FROM users WHERE name = ?)",
			[]interface{}{"John"},
		},
		{
			"OrWhereIn",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("age", ">", 19).OrWhereIn("id", []int64{1, 2, 3})
			},
			"SELECT  FROM  WHERE age > ? OR id IN (?, ?, ?)",
			[]interface{}{19, 1, 2, 3},
		},
		{
			"OrWhereIn (Subquery)",
			func() *query.Builder {
				sq := query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Select("id").Table("users").Where("name", "=", "John")
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("age", ">", 19).OrWhereIn("id", sq)
			},
			"SELECT  FROM  WHERE age > ? OR id IN (SELECT id FROM users WHERE name = ?)",
			[]interface{}{19, "John"},
		},
		{
			"OrWhereNotIn",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("age", ">", 19).OrWhereNotIn("id", []int64{1, 2, 3})
			},
			"SELECT  FROM  WHERE age > ? OR id NOT IN (?, ?, ?)",
			[]interface{}{19, 1, 2, 3},
		},
		{
			"OrWhereNotIn (Subquery)",
			func() *query.Builder {
				sq := query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Select("id").Table("users").Where("name", "=", "John")
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("age", ">", 19).OrWhereNotIn("id", sq)
			},
			"SELECT  FROM  WHERE age > ? OR id NOT IN (SELECT id FROM users WHERE name = ?)",
			[]interface{}{19, "John"},
		},
		{
			"WhereInSubquery",
			func() *query.Builder {
				sq := query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Select("id").Table("users").Where("name", "=", "John")
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereInSubQuery("id", sq)
			},
			"SELECT  FROM  WHERE id IN (SELECT id FROM users WHERE name = ?)",
			[]interface{}{"John"},
		},

		{
			"WhereNotInSubquery",
			func() *query.Builder {
				sq := query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Select("id").Table("users").Where("name", "=", "John")
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereNotInSubQuery("id", sq)
			},
			"SELECT  FROM  WHERE id NOT IN (SELECT id FROM users WHERE name = ?)",
			[]interface{}{"John"},
		},
		{
			"OrWhereInSubquery",
			func() *query.Builder {
				sq := query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Select("id").Table("users").Where("name", "=", "John")
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("age", ">", 19).OrWhereInSubQuery("id", sq)
			},
			"SELECT  FROM  WHERE age > ? OR id IN (SELECT id FROM users WHERE name = ?)",
			[]interface{}{19, "John"},
		},
		{
			"OrWhereNotInSubquery",
			func() *query.Builder {
				sq := query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Select("id").Table("users").Where("name", "=", "John")
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("age", ">", 19).OrWhereNotInSubQuery("id", sq)
			},
			"SELECT  FROM  WHERE age > ? OR id NOT IN (SELECT id FROM users WHERE name = ?)",
			[]interface{}{19, "John"},
		},
		{
			"WhereNull",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereNull("deleted_at")
			},
			"SELECT  FROM  WHERE deleted_at IS NULL",
			nil,
		},
		{
			"WhereNotNull",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereNotNull("deleted_at")
			},
			"SELECT  FROM  WHERE deleted_at IS NOT NULL",
			nil,
		},
		{
			"OrWhereNull",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("age", ">", 18).OrWhereNull("deleted_at")
			},
			"SELECT  FROM  WHERE age > ? OR deleted_at IS NULL",
			[]interface{}{18},
		},
		{
			"OrWhereNotNull",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("age", ">", 18).OrWhereNotNull("deleted_at")
			},
			"SELECT  FROM  WHERE age > ? OR deleted_at IS NOT NULL",
			[]interface{}{18},
		},
		{
			"WhereColumn",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereColumn([]string{"created_at", "updated_at", "deleted_at"}, "created_at", "=", "updated_at")
			},
			"SELECT  FROM  WHERE created_at = updated_at",
			nil,
		},
		{
			"WhereColumn_With_Operator",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereColumn([]string{"created_at", "updated_at", "deleted_at"}, "created_at", ">", "updated_at")
			},
			"SELECT  FROM  WHERE created_at > updated_at",
			nil,
		},
		{
			"OrWhereColumn",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("age", ">", 18).OrWhereColumn([]string{"created_at", "updated_at", "deleted_at"}, "created_at", "=", "updated_at")
			},
			"SELECT  FROM  WHERE age > ? OR created_at = updated_at",
			[]interface{}{18},
		},
		{
			"OrWhereColumn_With_Operator",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("age", ">", 18).OrWhereColumn([]string{"created_at", "updated_at", "deleted_at"}, "created_at", ">", "updated_at")
			},
			"SELECT  FROM  WHERE age > ? OR created_at > updated_at",
			[]interface{}{18},
		},
		{
			"WhereColumns",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereColumns([]string{"created_at", "updated_at", "deleted_at"}, [][]string{{"created_at", "=", "updated_at"}, {"deleted_at", "=", "updated_at"}})
			},
			"SELECT  FROM  WHERE created_at = updated_at AND deleted_at = updated_at",
			nil,
		},
		{
			"OrWhereColumns",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).OrWhereColumns([]string{"created_at", "updated_at", "deleted_at"}, [][]string{{"created_at", "=", "updated_at"}, {"deleted_at", "=", "updated_at"}})
			},
			"SELECT  FROM  WHERE created_at = updated_at OR deleted_at = updated_at",
			nil,
		},
		{
			"WhereBetween",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereBetween("age", 18, 30)
			},
			"SELECT  FROM  WHERE age BETWEEN ? AND ?",
			[]interface{}{18, 30},
		},
		{
			"OrWhereBetween",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("city", "=", "New York").OrWhereBetween("age", 18, 30)
			},
			"SELECT  FROM  WHERE city = ? OR age BETWEEN ? AND ?",
			[]interface{}{"New York", 18, 30},
		},
		{
			"WhereNotBetween",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereNotBetween("age", 18, 30)
			},
			"SELECT  FROM  WHERE age NOT BETWEEN ? AND ?",
			[]interface{}{18, 30},
		},
		{
			"OrWhereNotBetween",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("city", "=", "New York").OrWhereNotBetween("age", 18, 30)
			},
			"SELECT  FROM  WHERE city = ? OR age NOT BETWEEN ? AND ?",
			[]interface{}{"New York", 18, 30},
		},
		{
			"WhereBetweenColumns",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereBetweenColumns([]string{"created_at", "updated_at", "deleted_at"}, "created_at", "updated_at", "deleted_at")
			},
			"SELECT  FROM  WHERE created_at BETWEEN updated_at AND deleted_at",
			nil,
		},
		{
			"OrWhereBetweenColumns",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).OrWhereBetweenColumns([]string{"created_at", "updated_at", "deleted_at"}, "created_at", "updated_at", "deleted_at")
			},
			"SELECT  FROM  WHERE created_at BETWEEN updated_at AND deleted_at",
			nil,
		},
		{
			"WhereNotBetweenColumns",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereNotBetweenColumns([]string{"created_at", "updated_at", "deleted_at"}, "created_at", "updated_at", "deleted_at")
			},
			"SELECT  FROM  WHERE created_at NOT BETWEEN updated_at AND deleted_at",
			nil,
		},
		{
			"OrWhereNotBetweenColumns",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).OrWhereNotBetweenColumns([]string{"created_at", "updated_at", "deleted_at"}, "created_at", "updated_at", "deleted_at")
			},
			"SELECT  FROM  WHERE created_at NOT BETWEEN updated_at AND deleted_at",
			nil,
		},
		{
			"WhereExists",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereExists(func(b *query.Builder) *query.Builder {
					return b.Select("id").Table("users").Where("name", "=", "John")
				})
			},
			"SELECT  FROM  WHERE EXISTS (SELECT id FROM users WHERE name = ?)",
			[]interface{}{"John"},
		},
		{
			"OrWhereExists",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("city", "=", "New York").OrWhereExists(func(b *query.Builder) *query.Builder {
					return b.Select("id").Table("users").Where("name", "=", "John")
				})
			},
			"SELECT  FROM  WHERE city = ? OR EXISTS (SELECT id FROM users WHERE name = ?)",
			[]interface{}{"New York", "John"},
		},
		{
			"WhereNotExists",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereNotExists(func(b *query.Builder) *query.Builder {
					return b.Select("id").Table("users").Where("name", "=", "John")
				})
			},
			"SELECT  FROM  WHERE NOT EXISTS (SELECT id FROM users WHERE name = ?)",
			[]interface{}{"John"},
		},
		{
			"OrWhereNotExists",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("city", "=", "New York").OrWhereNotExists(func(b *query.Builder) *query.Builder {
					return b.Select("id").Table("users").Where("name", "=", "John")
				})
			},
			"SELECT  FROM  WHERE city = ? OR NOT EXISTS (SELECT id FROM users WHERE name = ?)",
			[]interface{}{"New York", "John"},
		},
		{
			"WhereExistsQuery",
			func() *query.Builder {
				q := query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Select("id").Table("users").Where("name", "=", "John")
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereExistsQuery(q)
			},
			"SELECT  FROM  WHERE EXISTS (SELECT id FROM users WHERE name = ?)",
			[]interface{}{"John"},
		},
		{
			"OrWhereExistsQuery",
			func() *query.Builder {
				q := query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Select("id").Table("users").Where("name", "=", "John")
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("city", "=", "New York").OrWhereExistsQuery(q)
			},
			"SELECT  FROM  WHERE city = ? OR EXISTS (SELECT id FROM users WHERE name = ?)",
			[]interface{}{"New York", "John"},
		},
		{
			"WhereNotExistsQuery",
			func() *query.Builder {
				q := query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Select("id").Table("users").Where("name", "=", "John")
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereNotExistsQuery(q)
			},
			"SELECT  FROM  WHERE NOT EXISTS (SELECT id FROM users WHERE name = ?)",
			[]interface{}{"John"},
		},
		{
			"OrWhereNotExistsQuery",
			func() *query.Builder {
				q := query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Select("id").Table("users").Where("name", "=", "John")
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("city", "=", "New York").OrWhereNotExistsQuery(q)
			},
			"SELECT  FROM  WHERE city = ? OR NOT EXISTS (SELECT id FROM users WHERE name = ?)",
			[]interface{}{"New York", "John"},
		},
		{
			"WhereFullText_MySQL",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereFullText([]string{"name", "note"}, "John Doe", map[string]interface{}{"mode": "boolean"})
			},
			"SELECT  FROM  WHERE MATCH (name, note) AGAINST (? IN BOOLEAN MODE)",
			[]interface{}{"John Doe"},
		},
		{
			"WhereFullText_PostgreSQL",
			func() *query.Builder {
				return query.NewBuilder(db.NewPostgreSQLQueryBuilder(), cache.NewAsyncQueryCache(100)).WhereFullText([]string{"name", "note"}, "John Doe", map[string]interface{}{"language": "english"})
			},
			"SELECT  FROM  WHERE (to_tsvector(?, name) || to_tsvector(?, note)) @@ plainto_tsquery(?, ?)",
			[]interface{}{"english", "english", "english", "John Doe"},
		},
		{
			"OrWhereFullText_MySQL",
			func() *query.Builder {
				return query.NewBuilder(db.NewMySQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("city", "=", "New York").OrWhereFullText([]string{"name", "note"}, "John Doe", map[string]interface{}{"mode": "boolean"})
			},
			"SELECT  FROM  WHERE city = ? OR MATCH (name, note) AGAINST (? IN BOOLEAN MODE)",
			[]interface{}{"New York", "John Doe"},
		},
		{

			"OrWhereFullText_PostgreSQL",
			func() *query.Builder {
				return query.NewBuilder(db.NewPostgreSQLQueryBuilder(), cache.NewAsyncQueryCache(100)).Where("city", "=", "New York").OrWhereFullText([]string{"name", "note"}, "John Doe", map[string]interface{}{"language": "english"})
			},
			"SELECT  FROM  WHERE city = ? OR (to_tsvector(?, name) || to_tsvector(?, note)) @@ plainto_tsquery(?, ?)",
			[]interface{}{"New York", "english", "english", "english", "John Doe"},
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

			convertedValues := sliceutils.ToInterfaceSlice(values)
			for i := range convertedValues {
				if values[i] != tt.expectedValues[i] {
					t.Errorf("expected value %v at index %d but got %v", tt.expectedValues[i], i, values[i])
				}
			}
		})
	}
}
