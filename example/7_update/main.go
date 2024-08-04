package main

import (
	"fmt"

	"github.com/faciam-dev/goquent-query-builder/api"
	"github.com/faciam-dev/goquent-query-builder/internal/cache"
)

func main() {
	// Initialize database strategy
	dbStrategy := mysql.NewMySQLQueryBuilder()

	asyncCache := cache.NewAsyncQueryCache(100)

	// Simple Update query
	// UPDATE users SET name = 'John Doe'
	//

	// Executing query: UPDATE users SET name = ? with values: [John Doe]

	qb := api.NewUpdateQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		Update(map[string]interface{}{"name": "John Doe"})

	query, values, err := qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)

	// Update query
	// UPDATE users SET name = 'John Doe'
	// WHERE users.id = 1
	//

	// Executing query: UPDATE users SET name = ? WHERE users.id = ? with values: [John Doe 1]

	qb = api.NewUpdateQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		Update(map[string]interface{}{"name": "John Doe"}).
		Where("users.id", "=", 1)

	query, values, err = qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)

	// Complex Update query with JoinQuery
	// UPDATE users
	// JOIN profiles ON users.id = profiles.user_id
	// SET users.name = 'John Doe'
	// WHERE users.id = 1
	//

	// Executing query: UPDATE users JOIN profiles ON users.id = profiles.user_id SET users.name = ? WHERE users.id = ? with values: [John Doe 1]

	qb = api.NewUpdateQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		Update(map[string]interface{}{"name": "John Doe"}).
		JoinQuery("profiles", func(b *api.JoinClauseQueryBuilder) {
			b.On("users.id", "=", "profiles.user_id")

		}).
		Where("users.id", "=", 1)

	query, values, err = qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)

	// Complex Update query with JoinQuery and multiple conditions
	// UPDATE users
	// JOIN profiles ON users.id = profiles.user_id
	// SET users.name = 'John Doe'
	// WHERE users.id = 1 AND profiles.age > 18
	//

	// Executing query: UPDATE users JOIN profiles ON users.id = profiles.user_id SET users.name = ? WHERE users.id = ? AND profiles.age > ? with values: [John Doe 1 18]

	qb = api.NewUpdateQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		Update(map[string]interface{}{"name": "John Doe"}).
		JoinQuery("profiles", func(b *api.JoinClauseQueryBuilder) {
			b.On("users.id", "=", "profiles.user_id")
		}).
		Where("users.id", "=", 1).
		Where("profiles.age", ">", 18)

	query, values, err = qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)

}
