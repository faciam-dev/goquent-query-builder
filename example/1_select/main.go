package main

import (
	"fmt"

	"github.com/faciam-dev/goquent-query-builder/api"
	"github.com/faciam-dev/goquent-query-builder/database/mysql"
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

	fmt.Println("Executing query:", query, "with values:", values)

}
