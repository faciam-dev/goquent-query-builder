package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/pkg/api"
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
	driver, err := mysql.WithInstance(d, &mysql.Config{})
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
	iqb := api.NewInsertBuilder(dbStrategy, asyncCache).
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

	query, values := iqb.Build()
	_, err = QueryToMap(d, query, values...)
	if err != nil {
		m.Down()
		log.Fatal(err)
	}

	// INSERT INTO profiles (user_id, age) VALUES (?, ?)
	iqb = api.NewInsertBuilder(dbStrategy, asyncCache).
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
	query, values = iqb.Build()
	_, err = QueryToMap(d, query, values...)
	if err != nil {
		m.Down()
		log.Fatal(err)
	}

	// SELECT users.id, users.name AS name FROM users JOIN profiles ON users.id = profiles.user_id WHERE profiles.age > 18 ORDER BY users.name ASC
	qb := api.NewSelectBuilder(dbStrategy, asyncCache).
		Table("users").
		Select("users.id AS id", "users.name AS name").
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
