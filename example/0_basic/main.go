package main

import (
	"fmt"
	"time"

	"github.com/faciam-dev/goquent-query-builder/api"
	"github.com/faciam-dev/goquent-query-builder/database/mysql"
	"github.com/faciam-dev/goquent-query-builder/internal/profiling"
)

func main() {
	// Initialize database strategy
	dbStrategy := mysql.NewMySQLQueryBuilder()

	// Executing query: SELECT `id`, `users`.`name` as `name` FROM `users` INNER JOIN `profiles` ON `users`.`id` = `profiles`.`user_id` WHERE `profiles`.`age` > ? ORDER BY `users`.`name` ASC with values: [18]
	qb := api.NewSelectQueryBuilder(dbStrategy).
		Table("users").
		Select("id", "users.name as name").
		Join("profiles", "users.id", "=", "profiles.user_id").
		Where("profiles.age", ">", 18).
		OrderBy("users.name", "ASC")

	query, values, err := qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	profiling.Profile(query, func() {
		fmt.Println("Executing query:", query, "with values:", values)
		time.Sleep(2 * time.Second) // Simulate query execution
	})

	// INSERT INTO users (age, name) VALUES (?, ?)
	iqb := api.NewInsertQueryBuilder(dbStrategy).
		Table("users").
		Insert(map[string]interface{}{
			"name": "John Doe",
			"age":  30,
		})

	query, values, _ = iqb.Build()

	profiling.Profile(query, func() {
		fmt.Println("Executing query:", query, "with values:", values)
		time.Sleep(2 * time.Second) // Simulate query execution
	})

	// UPDATE users SET age = ? WHERE id = ?
	uqb := api.NewUpdateQueryBuilder(dbStrategy).
		Table("users").
		Update(map[string]interface{}{
			"age": 40,
		}).
		Where("id", "=", 1)

	query, values, _ = uqb.Build()

	profiling.Profile(query, func() {
		fmt.Println("Executing query:", query, "with values:", values)
		time.Sleep(2 * time.Second) // Simulate query execution
	})

	// DELETE FROM users WHERE id = ?
	dqb := api.NewDeleteQueryBuilder(dbStrategy).
		Table("users").
		Where("id", "=", 1).
		Delete()

	query, values, _ = dqb.Build()

	profiling.Profile(query, func() {
		fmt.Println("Executing query:", query, "with values:", values)
		time.Sleep(2 * time.Second) // Simulate query execution
	})

}
