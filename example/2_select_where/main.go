package main

import (
	"fmt"

	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/pkg/api"
)

func main() {
	// Initialize database strategy
	dbStrategy := mysql.NewMySQLQueryBuilder()

	asyncCache := cache.NewAsyncQueryCache(100)

	// Complex query with WhereGroup
	//
	// SELECT users.id, users.name as name
	// FROM users
	// JOIN profiles ON users.id = profiles.user_id
	// WHERE (profiles.age > 18)
	// ORDER BY users.name ASC
	//
	// Executing query: SELECT users.id, users.name as name FROM users JOIN profiles ON users.id = profiles.user_id WHERE (profiles.age > 18) ORDER BY users.name ASC with values: [18]
	qb := api.NewSelectQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		Select("id", "users.name as name").
		Join("profiles", "users.id", "=", "profiles.user_id").
		WhereGroup(func(qb *api.WhereSelectQueryBuilder) {
			qb.Where("profiles.age", ">", 18)
		}).
		OrderBy("users.name", "ASC")

	query, values, err := qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)

	// Query with WhereNotBetween
	//
	// SELECT users.id, users.name as name
	// FROM users
	// WHERE users.age NOT BETWEEN 18 AND 30
	//
	// Executing query: SELECT users.id, users.name as name FROM users WHERE (users.age NOT BETWEEN 18 AND 30) with values: [18 30]
	qb = api.NewSelectQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		Select("id", "users.name as name").
		WhereNotBetween("users.age", 18, 30)

	query, values, err = qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)

	// Complex query with Union
	//
	// SELECT users.id, users.name as name
	// FROM users
	// WHERE users.age NOT BETWEEN 18 AND 30
	// UNION
	// SELECT users.id, users.name as name
	// FROM users
	// WHERE users.age BETWEEN 18 AND 30
	//
	// Executing query: SELECT users.id, users.name as name FROM users WHERE (users.age NOT BETWEEN 18 AND 30) UNION SELECT users.id, users.name as name FROM users WHERE (users.age BETWEEN 18 AND 30) with values: [18 30 18 30]
	qb = api.NewSelectQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		Select("id", "users.name as name").
		WhereNotBetween("users.age", 18, 30).
		Union(
			api.NewSelectQueryBuilder(dbStrategy, asyncCache).
				Table("users").
				Select("id", "users.name as name").
				WhereBetween("users.age", 18, 30),
		)

	query, values, err = qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)

	// Qery with WhereExists

	// SELECT users.id, users.name as name
	// FROM users
	// WHERE EXISTS (SELECT name, age FROM profiles WHERE age > 18)
	//

	// Executing query: SELECT users.id, users.name as name FROM users WHERE EXISTS (SELECT name, age FROM profiles WHERE age > 18) with values: [18]
	qb = api.NewSelectQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		Select("id", "users.name as name").
		WhereExists(func(q *api.SelectQueryBuilder) {
			q.Table("profiles").
				Select("name", "age").
				Where("age", ">", 18)
		})

	query, values, err = qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)

	// Query with WhereNotExists

	// SELECT users.id, users.name as name
	// FROM users
	// WHERE NOT EXISTS (SELECT name, age FROM profiles WHERE age > 18)
	//

	// Executing query: SELECT users.id, users.name as name FROM users WHERE NOT EXISTS (SELECT name, age FROM profiles WHERE age > 18) with values: [18]

	qb = api.NewSelectQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		Select("id", "users.name as name").
		WhereNotExists(func(q *api.SelectQueryBuilder) {
			q.Table("profiles").
				Select("name", "age").
				Where("age", ">", 18)
		})

	query, values, err = qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)

	// Query with WhereBetweenColumns

	// SELECT users.id, users.name as name
	// FROM users
	// WHERE users.age BETWEEN users.min_age AND users.max_age
	//

	// Executing query: SELECT users.id, users.name as name FROM users WHERE (users.age BETWEEN users.min_age AND users.max_age) with values: []

	qb = api.NewSelectQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		Select("id", "users.name as name").
		WhereBetweenColumns([]string{"users.age", "users.min_age", "users.max_age"}, "users.age", "users.min_age", "users.max_age")

	query, values, err = qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)

	// Query with WhereDate

	// SELECT users.id, users.name as name
	// FROM users
	// WHERE DATE(users.created_at) = '2021-01-01'
	//

	// Executing query: SELECT users.id, users.name as name FROM users WHERE DATE(users.created_at) = '2021-01-01' with values: [2021-01-01]

	qb = api.NewSelectQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		Select("id", "users.name as name").
		WhereDate("users.created_at", "=", "2021-01-01")

	query, values, err = qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)

}
