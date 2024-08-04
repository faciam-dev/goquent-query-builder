package main

import (
	"fmt"

	"github.com/faciam-dev/goquent-query-builder/api"
	"github.com/faciam-dev/goquent-query-builder/cache"
	"github.com/faciam-dev/goquent-query-builder/database/mysql"
)

func main() {
	// Initialize database strategy
	dbStrategy := mysql.NewMySQLQueryBuilder()

	asyncCache := cache.NewAsyncQueryCache(100)

	// Simple Insert query
	// INSERT INTO users (age, name) VALUES (30, 'John Doe')
	//
	// INSERT INTO `users` (`age`, `name`) VALUES (?, ?) with values: [30 John Doe]
	qb := api.NewInsertQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		Insert(map[string]interface{}{"name": "John Doe", "age": 30})

	query, values, err := qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)

	// Insert query with multiple values
	// INSERT INTO users (age, name) VALUES (30, 'John Doe'), (25, 'Jane Doe'), (20, 'Alice')
	//
	// Executing query: INSERT INTO `users` (`age`, `name`) VALUES (?, ?), (?, ?), (?, ?) with values: [30 John Doe 25 Jane Doe 20 Alice]
	qb = api.NewInsertQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		InsertBatch([]map[string]interface{}{
			{"name": "John Doe", "age": 30},
			{"name": "Jane Doe", "age": 25},
			{"name": "Alice", "age": 20},
		})

	query, values, err = qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)

	// Insert query from SelectQueryBuilder
	// INSERT INTO users (age, name) SELECT age, name FROM profiles WHERE age > 18

	// Executing query: INSERT INTO `users` (`age`, `name`) SELECT `age`, `name` FROM `profiles` WHERE `age` > ? with values: [18]
	qb = api.NewInsertQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		InsertUsing([]string{"age", "name"},
			api.NewSelectQueryBuilder(dbStrategy, asyncCache).
				Table("profiles").
				Select("age", "name").
				Where("age", ">", 18),
		)

	query, values, err = qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)

}
