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

	// Simple Delete query
	// DELETE FROM users
	//
	// Executing query: DELETE FROM users with values: []
	qb := api.NewDeleteQueryBuilder(dbStrategy, asyncCache).
		Table("users")

	query, values, err := qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)

	// Delete query
	// DELETE FROM users
	// WHERE users.id = 1
	//
	// Executing query: DELETE FROM users WHERE users.id = ? with values: [1]
	qb = api.NewDeleteQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		Where("users.id", "=", 1)

	query, values, err = qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)

	// Complex Delete query with JoinQuery
	// DELETE FROM users
	// JOIN profiles ON users.id = profiles.user_id
	// WHERE users.id = 1
	//
	// Executing query: DELETE FROM users JOIN profiles ON users.id = profiles.user_id WHERE users.id = ? with values: [1]
	qb = api.NewDeleteQueryBuilder(dbStrategy, asyncCache).
		Table("users").
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

	// Complex Delete query with JoinQuery and multiple conditions
	// DELETE FROM users
	// JOIN profiles ON users.id = profiles.user_id
	// WHERE users.id = 1 AND profiles.age > 18
	//
	// Executing query: DELETE FROM users JOIN profiles ON users.id = profiles.user_id WHERE users.id = ? AND profiles.age > ? with values: [1 18]
	qb = api.NewDeleteQueryBuilder(dbStrategy, asyncCache).
		Table("users").
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

	// Complex Delete query with JoinQuery and multiple conditions
	// DELETE FROM users
	// JOIN profiles ON users.id = profiles.user_id
	// WHERE users.id = 1 AND profiles.age > 18
	// ORDER BY users.name ASC
	//

	// Executing query: DELETE FROM users JOIN profiles ON users.id = profiles.user_id WHERE users.id = ? AND profiles.age > ? ORDER BY users.name ASC with values: [1 18]
	qb = api.NewDeleteQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		JoinQuery("profiles", func(b *api.JoinClauseQueryBuilder) {

			b.On("users.id", "=", "profiles.user_id")
		}).
		Where("users.id", "=", 1).
		Where("profiles.age", ">", 18).
		OrderBy("users.name", "ASC")

	query, values, err = qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)

}
