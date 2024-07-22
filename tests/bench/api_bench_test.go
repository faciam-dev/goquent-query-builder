package bench_test

import (
	"testing"

	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/pkg/api"
)

func BenchmarkSimpleSelectQuery(b *testing.B) {

	dbStrategy := db.NewMySQLQueryBuilder()

	blankCache := cache.NewBlankQueryCache()

	qb := api.NewQueryBuilder(dbStrategy, blankCache).
		Table("users").
		Select("id", "users.name AS name")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qb.Build()
		//blankCache.SetWithExpiry(query, query, 5*time.Minute)
	}

	// go test -benchmem -run=^$ -bench BenchmarkSimpleSelectQuery -benchtime=1s
	// before refactor
	// 2335863               509.9 ns/op           528 B/op        18 allocs/op
	// after refactor
	// 9414712               123.5 ns/op         416 B/op           3 allocs/op
}

func BenchmarkNormalSelectQuery(b *testing.B) {
	dbStrategy := db.NewMySQLQueryBuilder()

	blankCache := cache.NewBlankQueryCache()

	qb := api.NewQueryBuilder(dbStrategy, blankCache).
		Table("users").
		Select("id", "users.name AS name").
		Join("profiles", "users.id", "=", "profiles.user_id").
		Where("profiles.age", ">", 18).
		OrderBy("users.name", "ASC")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qb.Build()
		//blankCache.SetWithExpiry(query, query, 5*time.Minute)
	}

	// go test -benchmem -run=^$ -bench BenchmarkNormalSelectQuery -benchtime=1s
	// before refactor
	// 818690              1405 ns/op            1738 B/op        47 allocs/op
	// after refactor
	// 3166833               373.4 ns/op         816 B/op           6 allocs/op
}

func BenchmarkComplexSelectQuery(b *testing.B) {

	dbStrategy := db.NewMySQLQueryBuilder()

	blankCache := cache.NewBlankQueryCache()

	qb := api.NewQueryBuilder(dbStrategy, blankCache).
		Table("users").
		Select("id", "users.name AS name").
		Join("profiles", "users.id", "=", "profiles.user_id").
		Where("profiles.age", ">", 18).
		OrderBy("users.name", "ASC").
		OrderBy("profiles.age", "DESC").
		GroupBy("users.id").
		Having("COUNT(profiles.id)", ">", 1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qb.Build()
		//blankCache.SetWithExpiry(query, query, 5*time.Minute)
	}

	// go test -benchmem -run=^$ -bench BenchmarkComplexSelectQuery -benchtime=1s
	// before refactor
	// 675976              1747 ns/op           2323 B/op          57 allocs/op
	// after refactor
	// 2800032               431.2 ns/op          832 B/op          7 allocs/op
}

func BenchmarkComplexSelectQueryWithUsingSubQuery(b *testing.B) {
	dbStrategy := db.NewMySQLQueryBuilder()
	blankCache := cache.NewBlankQueryCache()

	qb := api.NewQueryBuilder(dbStrategy, blankCache).
		Table("users").
		Select("id", "users.name AS name").
		Join("profiles", "users.id", "=", "profiles.user_id").
		Where("profiles.age", ">", 18).
		OrderBy("users.name", "ASC").
		OrderBy("profiles.age", "DESC").
		GroupBy("users.id").
		Having("COUNT(profiles.id)", ">", 1).
		Where("users.id", "IN", func(qb *api.QueryBuilder) {
			qb.Table("users").
				Select("id").
				Join("profiles", "users.id", "=", "profiles.user_id").
				Where("profiles.age", ">", 18)
		})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qb.Build()
		//blankCache.SetWithExpiry(query, query, 5*time.Minute)
	}

	// go test -benchmem -run=^$ -bench BenchmarkComplexSelectQueryWithUsingSubQuery -benchtime=1s
	// before refactor
	// 640671              2157 ns/op            890 B/op          25 allocs/op
	// after refactor
	// 2204690               545.7 ns/op         1184 B/op          9 allocs/op
	// c.f.) use AsyncQueryCache
	// 2705778             446.4 ns/op          673 B/op          14 allocs/op
}

func BenchmarkSimpleInsert(b *testing.B) {
	dbStrategy := db.NewMySQLQueryBuilder()
	blankCache := cache.NewBlankQueryCache()

	qb := api.NewInsertQueryBuilder(dbStrategy, blankCache).
		Table("users").
		Insert(map[string]interface{}{
			"name": "John",
			"age":  30,
		})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qb.Build()
	}

	// go test -benchmem -run=^$ -bench BenchmarkSimpleInsert -benchtime=1s
	// before refactor
	//  5521004               220.6 ns/op           216 B/op         8 allocs/op
	// after refactor
	//  7707033               154.6 ns/op           576 B/op         3 allocs/op
}

func BenchmarkInsertBatch(b *testing.B) {
	dbStrategy := db.NewMySQLQueryBuilder()
	blankCache := cache.NewBlankQueryCache()

	qb := api.NewInsertQueryBuilder(dbStrategy, blankCache).
		Table("users").
		InsertBatch([]map[string]interface{}{
			{
				"name": "John",
				"age":  30,
			},
			{
				"name": "Mike",
				"age":  25,
			},
		})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qb.Build()
	}

	// go test -benchmem -run=^$ -bench BenchmarkInsertBatch -benchtime=1s
	// before refactor
	//  2263813               520.9 ns/op           472 B/op        17 allocs/op
	// after refactor
	//  3401209               352.6 ns/op          1184 B/op         5 allocs/op
}

func BenchmarkInsertUsing(b *testing.B) {
	dbStrategy := db.NewMySQLQueryBuilder()
	blankCache := cache.NewBlankQueryCache()

	qb := api.NewInsertQueryBuilder(dbStrategy, blankCache).
		Table("users").
		InsertBatch([]map[string]interface{}{
			{
				"name": "John",
				"age":  30,
			},
			{
				"name": "Mike",
				"age":  25,
			},
		}).
		InsertUsing([]string{"name", "age"}, api.NewQueryBuilder(dbStrategy, blankCache).
			Table("users").
			Select("id").
			Join("profiles", "users.id", "=", "profiles.user_id").
			Where("profiles.age", ">", 18).
			OrderBy("users.name", "ASC"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qb.Build()
	}

	// go test -benchmem -run=^$ -bench BenchmarkInsertUsing -benchtime=1s
	// before refactor (after select query's refactor)
	//  6140727               197.4 ns/op           336 B/op         8 allocs/op
	// after refactor
	//  6972561               169.6 ns/op           712 B/op         5 allocs/op
}

func BenchmarkSimpleUpdate(b *testing.B) {
	dbStrategy := db.NewMySQLQueryBuilder()
	blankCache := cache.NewBlankQueryCache()

	qb := api.NewUpdateBuilder(dbStrategy, blankCache).
		Table("users").
		Update(map[string]interface{}{
			"name": "Joe",
			"age":  31,
		})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qb.Build()
	}

	// go test -benchmem -run=^$ -bench BenchmarkSimpleUpdate -benchtime=1s
	// before refactor
	//  5506760               214.4 ns/op           336 B/op         8 allocs/op
	// after refactor
	//  7360680               164.7 ns/op           600 B/op         4 allocs/op
}

func BenchmarkUpdateWhere(b *testing.B) {
	dbStrategy := db.NewMySQLQueryBuilder()
	blankCache := cache.NewBlankQueryCache()

	qb := api.NewUpdateBuilder(dbStrategy, blankCache).
		Table("users").
		Where("id", "=", 1).
		Update(map[string]interface{}{
			"name": "Joe",
			"age":  31,
		})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qb.Build()
	}

	// go test -benchmem -run=^$ -bench BenchmarkUpdateWhere -benchtime=1s
	// before refactor
	//  3559383               339.8 ns/op           544 B/op        11 allocs/op
	// after refactor
	//  4128181               290.8 ns/op           808 B/op         7 allocs/op
}

func BenchmarkJoinUpdate(b *testing.B) {
	dbStrategy := db.NewMySQLQueryBuilder()
	blankCache := cache.NewBlankQueryCache()

	qb := api.NewUpdateBuilder(dbStrategy, blankCache).
		Table("users").
		Join("profiles", "users.id", "=", "profiles.user_id").
		Where("profiles.age", ">", 18).
		Update(map[string]interface{}{
			"name": "Joe",
			"age":  31,
		})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qb.Build()
	}

	// go test -benchmem -run=^$ -bench BenchmarkJoinUpdate -benchtime=1s
	// before refactor
	//  3003182               396.5 ns/op           544 B/op        11 allocs/op
	// after refactor
	//  3419262               349.0 ns/op           808 B/op         7 allocs/op
}
