package main

import (
	"fmt"

	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/pkg/api"
)

func main() {
	// Initialize database strategy
	dbStrategy := db.NewMySQLQueryBuilder()

	asyncCache := cache.NewAsyncQueryCache(100)

	// Simple Insert query
	// INSERT INTO users (age, name) VALUES (30, 'John Doe')
	//
	// Executing query: INSERT INTO users (age, name) VALUES (?, ?) with values: [30 John Doe]
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

	// Insert query from SelectQueryBuilder
	// INSERT INTO users (age, name) SELECT age, name FROM profiles WHERE age > 18

	// Executing query: INSERT INTO users (age, name) SELECT age, name FROM profiles WHERE age > ? with values: [18]
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
