package main

import (
	"fmt"

	"github.com/faciam-dev/goquent-query-builder/api"
	"github.com/faciam-dev/goquent-query-builder/database/mysql"
)

func main() {
	// Initialize database strategy
	dbStrategy := mysql.NewMySQLQueryBuilder()

	// Complex query with WhereGroup and having
	//
	// SELECT users.id, users.name as name
	// FROM users
	// JOIN profiles ON users.id = profiles.user_id
	// WHERE (profiles.age > 18)
	// GROUP BY users.id
	// HAVING COUNT(profiles.id) > 1
	// ORDER BY users.name ASC
	//

	// Executing query: SELECT `id`, `users`.`name` as `name` FROM `users` INNER JOIN `profiles` ON `users`.`id` = `profiles`.`user_id` WHERE (`profiles`.`age` > ?) GROUP BY `users`.`id` HAVING `COUNT(profiles`.`id)` > ? ORDER BY `users`.`name` ASC with values: [18 1]
	qb := api.NewSelectQueryBuilder(dbStrategy).
		Table("users").
		Select("id", "users.name as name").
		Join("profiles", "users.id", "=", "profiles.user_id").
		WhereGroup(func(qb *api.WhereSelectQueryBuilder) {
			qb.Where("profiles.age", ">", 18)
		}).
		GroupBy("users.id").
		Having("COUNT(profiles.id)", ">", 1).
		OrderBy("users.name", "ASC")

	query, values, err := qb.Build()

	if err != nil {
		panic(err)
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)

	// Complex query with WhereGroup and having and JoinQuery
	//
	// SELECT users.id, users.name as name
	// FROM users
	// JOIN profiles ON users.id = profiles.user_id
	// AND profiles.age > 18
	// WHERE (profiles.age > 18)
	// GROUP BY users.id
	// HAVING COUNT(profiles.id) > 1
	// ORDER BY users.name ASC

	// Executing query: SELECT `id`, `users`.`name` as `name` FROM `users` INNER JOIN `profiles` ON `users`.`id` = `profiles`.`user_id` AND `profiles`.`age` > ? WHERE (`profiles`.`age` > ?) GROUP BY `users`.`id` HAVING `COUNT(profiles`.`id)` > ? ORDER BY `users`.`name` ASC with values: [18 18 1]

	qb = api.NewSelectQueryBuilder(dbStrategy).
		Table("users").
		Select("id", "users.name as name").
		JoinQuery("profiles", func(b *api.JoinClauseQueryBuilder) {
			b.On("users.id", "=", "profiles.user_id").
				Where("profiles.age", ">", 18)
		}).
		WhereGroup(func(qb *api.WhereSelectQueryBuilder) {
			qb.Where("profiles.age", ">", 18)
		}).
		GroupBy("users.id").
		Having("COUNT(profiles.id)", ">", 1).
		OrderBy("users.name", "ASC")

	query, values, err = qb.Build()

	if err != nil {
		panic(err)
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)

	// Complex query with WhereGroup and having and JoinQuery and multiple conditions

	// SELECT users.id, users.name as name
	// FROM users
	// JOIN profiles ON users.id = profiles.user_id
	// AND profiles.age > 18
	// WHERE (profiles.age > 18)
	// GROUP BY users.id
	// HAVING COUNT(profiles.id) > 1
	// ORDER BY users.name ASC
	// LIMIT 1
	//

	// Executing query: SELECT `id`, `users`.`name` as `name` FROM `users` INNER JOIN `profiles` ON `users`.`id` = `profiles`.`user_id` AND `profiles`.`age` > ? WHERE (`profiles`.`age` > ?) GROUP BY `users`.`id` HAVING `COUNT(profiles`.`id)` > ? ORDER BY `users`.`name` ASC LIMIT 1 with values: [18 18 1]

	qb = api.NewSelectQueryBuilder(dbStrategy).
		Table("users").
		Select("id", "users.name as name").
		JoinQuery("profiles", func(b *api.JoinClauseQueryBuilder) {
			b.On("users.id", "=", "profiles.user_id").
				Where("profiles.age", ">", 18)
		}).
		WhereGroup(func(qb *api.WhereSelectQueryBuilder) {
			qb.Where("profiles.age", ">", 18)
		}).
		GroupBy("users.id").
		Having("COUNT(profiles.id)", ">", 1).
		OrderBy("users.name", "ASC").
		Limit(1)

	query, values, err = qb.Build()

	if err != nil {
		panic(err)
	}

	// Print the query and values
	fmt.Println("Executing query:", query, "with values:", values)
}
