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

	asyncCache := cache.NewAsyncQueryCache()

	qb := api.NewQueryBuilder(dbStrategy, asyncCache).
		Table("users").
		Select("id", "name").
		//Join("profiles", "users.id", "=", "profiles.user_id").
		Where("age", ">", 18).
		OrderBy("name", "ASC")

	//txManager := transaction.NewTransactionManager()
	query, values := qb.Build()

	// プロファイリング
	profiling.Profile(query, func() {
		fmt.Println("Executing query:", query, "with values:", values)
		time.Sleep(2 * time.Second)                           // Simulate query execution
		asyncCache.SetWithExpiry(query, query, 5*time.Minute) // キャッシュに5分間保存
	})

	// 非同期実行
	//async.ExecuteAsync(query, values)
}
