package main

import (
	"time"
    "database/sql"
    "fmt"
    "github.com/go-redis/redis"
    _ "github.com/go-sql-driver/mysql"
	"strconv"
)

type Post struct {
	ID           int       `db:"id"`
	UserID       int       `db:"user_id"`
	Body         string    `db:"body"`
	Mime         string    `db:"mime"`
	CreatedAt    time.Time `db:"created_at"`
}

func main() {
    // MySQLから構造体データを取得
    db, err := sql.Open("mysql","isuconp:isuconp@(localhost:3306)/isuconp?parseTime=true")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    rows, err := db.Query("SELECT id, user_id, body, mime, created_at FROM posts")
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    // Redisにハッシュ型として保存
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.UserID, &post.Body, &post.Mime, &post.CreatedAt)
        if err != nil {
            panic(err)
        }

        // 構造体をRedisのハッシュ型として保存
        key := fmt.Sprintf("post:%d", post.ID)
        user_id := strconv.Itoa(post.UserID)
        body := string(post.Body)
        mime := string(post.Mime)
        created_at := time.Time(post.CreatedAt)

        err = redisClient.HSet(key, "user_id", user_id).Err()
        if err != nil {
            panic(err)
        }

        err = redisClient.HSet(key, "body", body).Err()
        if err != nil {
            panic(err)
        }

        err = redisClient.HSet(key, "mime", mime).Err()
        if err != nil {
            panic(err)
        }

        err = redisClient.HSet(key, "created_at", created_at).Err()
        if err != nil {
            panic(err)
        }
    }
    err = rows.Err()
    if err != nil {
        panic(err)
    }

    key := "post:1903"
    results, err := redisClient.HGetAll(key).Result()
    if err != nil {
        panic(err)
    }

	// 取得したデータを表示
	fmt.Println("Hash data:")
	for key, value := range results {
		fmt.Printf("%s: %s\n", key, value)
	}
}
