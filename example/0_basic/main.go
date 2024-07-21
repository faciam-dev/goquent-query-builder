package main

import (
	"fmt"
	"time"

	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/profiling"
	"github.com/faciam-dev/goquent-query-builder/pkg/api"
)

func main() {
	// データベースごとのクエリビルダーストラテジーを選択
	dbStrategy := &db.MySQLQueryBuilder{}

	asyncCache := cache.NewAsyncQueryCache(100)

	// SELECT users.id, users.name AS name FROM users JOIN profiles ON users.id = profiles.user_id WHERE profiles.age > 18 ORDER BY users.name ASC
	qb := api.NewQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		Select("id", "users.name AS name").
		Join("profiles", "users.id", "=", "profiles.user_id").
		Where("profiles.age", ">", 18).
		OrderBy("users.name", "ASC")

	query, values := qb.Build()

	profiling.Profile(query, func() {
		fmt.Println("Executing query:", query, "with values:", values)
		time.Sleep(2 * time.Second) // Simulate query execution
		asyncCache.SetWithExpiry(query, query, 5*time.Minute)
	})

	// use cache
	query, values = qb.Build()
	profiling.Profile(query, func() {
		fmt.Println("Executing query:", query, "with values:", values)
		time.Sleep(2 * time.Second) // Simulate query execution
		asyncCache.SetWithExpiry(query, query, 5*time.Minute)
	})

	// INSERT INTO users (age, name) VALUES (?, ?)
	iqb := api.NewInsertQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		Insert(map[string]interface{}{
			"name": "John Doe",
			"age":  30,
		})

	query, values = iqb.Build()

	profiling.Profile(query, func() {
		fmt.Println("Executing query:", query, "with values:", values)
		time.Sleep(2 * time.Second) // Simulate query execution
		asyncCache.SetWithExpiry(query, query, 5*time.Minute)
	})

	// UPDATE users SET age = ? WHERE id = ?
	uqb := api.NewUpdateQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		Update(map[string]interface{}{
			"age": 40,
		}).
		Where("id", "=", 1)

	query, values = uqb.Build()

	profiling.Profile(query, func() {
		fmt.Println("Executing query:", query, "with values:", values)
		time.Sleep(2 * time.Second) // Simulate query execution
		asyncCache.SetWithExpiry(query, query, 5*time.Minute)
	})

	// DELETE FROM users WHERE id = ?
	dqb := api.NewDeleteQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		Where("id", "=", 1).
		Delete()

	query, values = dqb.Build()

	profiling.Profile(query, func() {
		fmt.Println("Executing query:", query, "with values:", values)
		time.Sleep(2 * time.Second) // Simulate query execution
		asyncCache.SetWithExpiry(query, query, 5*time.Minute)
	})

}
