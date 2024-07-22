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
	// 3636319               318.8 ns/op           488 B/op        14 allocs/op
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
	// 950131              1190 ns/op            1810 B/op        42 allocs/op
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
	// 785894              1486 ns/op           2139 B/op          51 allocs/op
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
	// 717132              1887 ns/op            738 B/op          55 allocs/op
	// c.f.) use AsyncQueryCache
	// 2705778             446.4 ns/op          673 B/op          14 allocs/op
}
