package main

import (
	"time"
    "fmt"
    "github.com/go-redis/redis"
    _ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
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
    db, err := sqlx.Open("mysql","isuconp:isuconp@(localhost:3306)/isuconp?parseTime=true")
    if err != nil {
        panic(err)
    }
    defer db.Close()

	const queryCacheKey = "SELECT `id`, `user_id`, `body`, `mime`, `created_at` FROM `posts` ORDER BY `created_at` DESC"

    rows, err := db.Query(queryCacheKey)
    if err != nil {
        return
    }
    defer rows.Close()

    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    // save
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.UserID, &post.Body, &post.Mime, &post.CreatedAt)
        if err != nil {
            panic(err)
        }

        // 構造体をRedisのハッシュ型として保存
        key := fmt.Sprintf("%d", post.ID)
        user_id := strconv.Itoa(post.UserID)
        body := string(post.Body)
        mime := string(post.Mime)
        created_at := timeToString(post.CreatedAt)

        err = rdb.HSet(key, "user_id", user_id).Err()
        if err != nil {
            panic(err)
        }

        err = rdb.HSet(key, "body", body).Err()
        if err != nil {
            panic(err)
        }

        err = rdb.HSet(key, "mime", mime).Err()
        if err != nil {
            panic(err)
        }

        err = rdb.HSet(key, "created_at", created_at).Err()
        if err != nil {
            panic(err)
        }
    }
    err = rows.Err()
    if err != nil {
        panic(err)
    }

    // get
    keys, err := rdb.Keys("*").Result()
    if err != nil {
        return
    }
    result := make(map[string]map[string]string)
    for _, key := range keys {
        fields, err := rdb.HGetAll(key).Result()
        if err != nil {
			return
        }
        result[key] = fields
    }

    var people []Post
    for id, data := range result {
        var person Post
        person.ID, _ = strconv.Atoi(id)
        person.UserID, _ = strconv.Atoi(data["user_id"])
        person.Body = data["body"]
        person.Mime = data["mime"]
        person.CreatedAt = stringToTime(data["created_at"])
        people = append(people, person)
    }
    fmt.Println(people)
}

var layout = "2006-01-02 15:04:05"

func timeToString(t time.Time) string {
    str := t.Format(layout)
    return str
}

func stringToTime(str string) time.Time {
    t, _ := time.Parse(layout, str)
    return t
}