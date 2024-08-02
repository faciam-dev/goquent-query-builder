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

	// SELECT users.id, users.name AS name FROM users JOIN profiles ON users.id = profiles.user_id WHERE profiles.age > 18 ORDER BY users.name ASC
	qb := api.NewSelectQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		Select("id", "users.name AS name").
		Join("profiles", "users.id", "=", "profiles.user_id").
		Where("profiles.age", ">", 18).
		OrderBy("users.name", "ASC")

	query, values, err := qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Executing query:", query, "with values:", values)

}
