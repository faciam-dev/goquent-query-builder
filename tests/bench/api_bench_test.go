package bench_test

import (
	"testing"
	"time"

	"github.com/faciam-dev/goquent-query-builder/internal/cache"
	"github.com/faciam-dev/goquent-query-builder/internal/db"
	"github.com/faciam-dev/goquent-query-builder/internal/profiling"
	"github.com/faciam-dev/goquent-query-builder/pkg/api"
)

func BenchmarkSimpleSelectQuery(b *testing.B) {
	dbStrategy := &db.MySQLQueryBuilder{}

	blankCache := cache.NewBlankQueryCache()

	qb := api.NewQueryBuilder(dbStrategy, blankCache).
		Table("users").
		Select("id", "users.name AS name")

	query, _ := qb.Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		profiling.Profile(query, func() {
			blankCache.SetWithExpiry(query, query, 5*time.Minute)
		})
	}

	// BenchmarkSimpleSelectQuery-32
	// before refactor
	// 153862             13493 ns/op        20 B/op          2 allocs/op
	// after refactor
	//
}

func BenchmarkNormalSelectQuery(b *testing.B) {
	dbStrategy := &db.MySQLQueryBuilder{}

	blankCache := cache.NewBlankQueryCache()

	qb := api.NewQueryBuilder(dbStrategy, blankCache).
		Table("users").
		Select("id", "users.name AS name").
		Join("profiles", "users.id", "=", "profiles.user_id").
		Where("profiles.age", ">", 18).
		OrderBy("users.name", "ASC")

	query, _ := qb.Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		profiling.Profile(query, func() {
			blankCache.SetWithExpiry(query, query, 5*time.Minute)
		})
	}

	// BenchmarkNormalSelectQuery-32
	// before refactor
	// 122570             16933 ns/op        20 B/op          2 allocs/op
	// after refactor
	//
}

func BenchmarkComplexSelectQuery(b *testing.B) {

	dbStrategy := &db.MySQLQueryBuilder{}

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

	query, _ := qb.Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		profiling.Profile(query, func() {
			blankCache.SetWithExpiry(query, query, 5*time.Minute)
		})
	}

	// BenchmarkComplexSelectQuery-32
	// before refactor
	// 135646             16405 ns/op        20 B/op          2 allocs/op
	// after refactor
	//
}

func BenchmarkComplexSelectQueryWithUsingSubQuery(b *testing.B) {
	dbStrategy := &db.MySQLQueryBuilder{}

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

	query, _ := qb.Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		profiling.Profile(query, func() {
			blankCache.SetWithExpiry(query, query, 5*time.Minute)
		})
	}

	// BenchmarkComplexSelectQueryWithUsingSubQuery-32
	// before refactor
	// 111592       18156 ns/op              20 B/op          2 allocs/op
	// after refactor
	//
}
