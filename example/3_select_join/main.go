package main

import (
	"fmt"

	"github.com/faciam-dev/goquent-query-builder/api"
	"github.com/faciam-dev/goquent-query-builder/database/mysql"
)

func main() {
	// Initialize database strategy
	dbStrategy := mysql.NewMySQLQueryBuilder()

	// Simple query with JOIN
	//
	// SELECT users.id, users.name as name
	// FROM users
	// JOIN profiles ON users.id = profiles.user_id
	// ORDER BY users.name ASC
	//

	// Executing query: SELECT `id`, `users`.`name` as `name` FROM `users` INNER JOIN `profiles` ON `users`.`id` = `profiles`.`user_id` ORDER BY `users`.`name` ASC with values: []

	qb := api.NewSelectQueryBuilder(dbStrategy).
		Table("users").
		Select("id", "users.name as name").
		Join("profiles", "users.id", "=", "profiles.user_id").
		OrderBy("users.name", "ASC")

	query, values, err := qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)

	// Complex query with JoinQuery

	// SELECT users.id, users.name as name
	// FROM users
	// JOIN profiles ON users.id = profiles.user_id
	// AND profiles.age > 18
	// ORDER BY users.name ASC

	// Executing query: SELECT `id`, `users`.`name` as `name` FROM `users` INNER JOIN `profiles` ON `users`.`id` = `profiles`.`user_id` AND `profiles`.`age` > ? ORDER BY `users`.`name` ASC with values: [18]

	qb = api.NewSelectQueryBuilder(dbStrategy).
		Table("users").
		Select("id", "users.name as name").
		JoinQuery("profiles", func(b *api.JoinClauseQueryBuilder) {
			b.On("users.id", "=", "profiles.user_id").
				Where("profiles.age", ">", 18)
		}).
		OrderBy("users.name", "ASC")

	query, values, err = qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values

	fmt.Println("Executing query:", query, "with values:", values)

	// Query with multiple JOINs

	// SELECT users.id, users.name as name
	// FROM users
	// JOIN profiles ON users.id = profiles.user_id
	// JOIN addresses ON users.id = addresses.user_id
	// AND profiles.age > 18
	// ORDER BY users.name ASC

	// Executing query: SELECT `id`, `users`.`name` as `name` FROM `users` INNER JOIN `profiles` ON `users`.`id` = `profiles`.`user_id` AND `profiles`.`age` > ? INNER JOIN `addresses` ON `users`.`id` = `addresses`.`user_id` ORDER BY `users`.`name` ASC with values: [18]

	qb = api.NewSelectQueryBuilder(dbStrategy).
		Table("users").
		Select("id", "users.name as name").
		JoinQuery("profiles", func(b *api.JoinClauseQueryBuilder) {
			b.On("users.id", "=", "profiles.user_id").
				Where("profiles.age", ">", 18)
		}).
		Join("addresses", "users.id", "=", "addresses.user_id").
		OrderBy("users.name", "ASC")

	query, values, err = qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values

	fmt.Println("Executing query:", query, "with values:", values)

	// Query with multiple JOINs and multiple conditions

	// SELECT users.id, users.name as name
	// FROM users
	// JOIN profiles ON users.id = profiles.user_id
	// JOIN addresses ON users.id = addresses.user_id
	// AND profiles.age > 18
	// AND addresses.city = 'New York'
	// ORDER BY users.name ASC

	// Executing query: SELECT `id`, `users`.`name` as `name` FROM `users` INNER JOIN `profiles` ON `users`.`id` = `profiles`.`user_id` AND `profiles`.`age` > ? INNER JOIN `addresses` ON `users`.`id` = `addresses`.`user_id` AND `addresses`.`city` = ? ORDER BY `users`.`name` ASC with values: [18 New York]

	qb = api.NewSelectQueryBuilder(dbStrategy).
		Table("users").
		Select("id", "users.name as name").
		JoinQuery("profiles", func(b *api.JoinClauseQueryBuilder) {
			b.On("users.id", "=", "profiles.user_id").
				Where("profiles.age", ">", 18)
		}).
		JoinQuery("addresses", func(b *api.JoinClauseQueryBuilder) {
			b.On("users.id", "=", "addresses.user_id").
				Where("addresses.city", "=", "New York")
		}).
		OrderBy("users.name", "ASC")

	query, values, err = qb.Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the query and values

	fmt.Println("Executing query:", query, "with values:", values)

	// Query with Lateral Join

	// SELECT users.id, users.name as name
	// FROM users
	// ,LATERAL(SELECT id, name FROM profiles WHERE users.id = profiles.user_id AND profiles.age > ?) as profiles
	// AND profiles.age > 18
	// ORDER BY users.name ASC

	// Executing query: SELECT `id`, `users`.`name` as `name` FROM `users` ,LATERAL(SELECT `id`, `name` FROM `profiles` WHERE `users`.`id` = `profiles`.`user_id` AND `profiles`.`age` > ?) as `profiles` WHERE `profiles`.`age` > ? ORDER BY `users`.`name` ASC with values: [18 18]
	qb = api.NewSelectQueryBuilder(dbStrategy).
		Table("users").
		Select("id", "users.name as name").
		JoinLateral(api.NewSelectQueryBuilder(dbStrategy).
			Table("profiles").
			Select("id", "name").
			WhereColumn([]string{"users.id", "profiles.user_id"}, "users.id", "=", "profiles.user_id").
			Where("profiles.age", ">", 18),
			"profiles").
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
