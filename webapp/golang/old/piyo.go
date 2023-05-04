package main

import (
	"fmt"
    "github.com/go-redis/redis/v8"
    "context"
)

func main() {
    // Redisに接続
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379", // Redisのアドレス
        Password: "",               // Redisに設定したパスワード（ない場合は空文字列）
        DB:       0,                // 使用するRedisのDB番号
    })

    ctx := context.Background()
    keys, err := rdb.Keys(ctx, "*").Result()
    if err != nil {
        return
    }
    result := make(map[string]map[string]string)
    for _, key := range keys {
        fields, err := rdb.HGetAll(ctx, key).Result()
        if err != nil {
			return
        }
        result[key] = fields
    }

	fmt.Println(result)
}