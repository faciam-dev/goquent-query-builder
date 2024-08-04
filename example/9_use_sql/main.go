package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/golang-migrate/migrate/v4"
	my "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/faciam-dev/goquent-query-builder/api"
	"github.com/faciam-dev/goquent-query-builder/cache"
	"github.com/faciam-dev/goquent-query-builder/database/mysql"
)

// QueryToMap is a helper function to convert sql.Rows to []map[string]interface{}
// for MySQL database
func main() {
	//connStr := "username:password@tcp(addr:port)/database"
	connStr := "root:1234@tcp(172.30.249.68:13306)/test"
	d, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	// migrate
	driver, err := my.WithInstance(d, &my.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		log.Fatal(err)
	}

	// Initialize database strategy
	dbStrategy := mysql.NewMySQLQueryBuilder()
	asyncCache := cache.NewAsyncQueryCache(100)

	// INSERT INTO users (age, name) VALUES (?, ?)
	iqb := api.NewInsertQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		InsertBatch([]map[string]interface{}{
			{
				"name": "John Doe",
			},
			{
				"name": "Jane Doe",
			},
			{
				"name": "Alice",
			},
		})

	query, values, _ := iqb.Build()
	result, err := Exec(d, query, values...)
	if err != nil {
		m.Down()
		log.Fatal(err)
	}
	fmt.Println(result)

	// INSERT INTO profiles (user_id, age) VALUES (?, ?)
	iqb = api.NewInsertQueryBuilder(dbStrategy, asyncCache).
		Table("profiles").
		InsertBatch([]map[string]interface{}{
			{
				"user_id": 1,
				"age":     35,
			},
			{
				"user_id": 2,
				"age":     25,
			},
			{
				"user_id": 3,
				"age":     15,
			},
		})
	query, values, _ = iqb.Build()
	result, err = Exec(d, query, values...)
	if err != nil {
		m.Down()
		log.Fatal(err)
	}
	fmt.Println(result)

	// SELECT `users`.`id` as `id`, `users`.`name` as `name` FROM `users` INNER JOIN `profiles` ON `users`.`id` = `profiles`.`user_id` WHERE `profiles`.`age` > ? ORDER BY `users`.`name` ASC
	qb := api.NewSelectQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		Select("users.id as id", "users.name as name").
		Join("profiles", "users.id", "=", "profiles.user_id").
		Where("profiles.age", ">", 18).
		OrderBy("users.name", "ASC")

	query, values, err = qb.Build()

	log.Default().Println("Executing query:", query, "with values:", values)

	if err != nil {
		m.Down()
		fmt.Println("Error:", err)
		return
	}

	// Execute query
	results, err := QueryToMap(d, query, values...)
	if err != nil {
		m.Down()
		log.Fatal(err)
	}

	// Print results
	for _, row := range results {
		fmt.Println(row)
		// Output:
		// map[id:2 name:Jane Doe]
		// map[id:1 name:John Doe]
	}

	m.Down()
}
