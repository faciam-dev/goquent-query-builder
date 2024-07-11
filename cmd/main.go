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
	dbStrategy := db.MySQLQueryBuilder{}

	qb := api.NewQueryBuilder(dbStrategy).
		Table("users").
		Select("id", "name").
		Where("age > ?", 18).
		OrderBy("name")

	asyncCache := cache.NewAsyncQueryCache()
	//txManager := transaction.NewTransactionManager()
	query, values := qb.Build()

	// プロファイリング
	profiling.Profile(query, func() {
		// キャッシュチェック
		if cachedQuery, found := asyncCache.Get(query); found {
			fmt.Println("Cache hit:", cachedQuery)
		} else {
			// トランザクション処理
			//txManager.Begin()
			fmt.Println("Executing query:", query, "with values:", values)
			time.Sleep(2 * time.Second) // Simulate query execution
			//txManager.Commit()
			asyncCache.SetWithExpiry(query, query, 5*time.Minute) // キャッシュに5分間保存
		}
	})

	// 非同期実行
	//async.ExecuteAsync(query, values)
}
