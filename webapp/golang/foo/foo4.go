package main

import (
        "encoding/json"
        "log"
        "time"

        "github.com/jmoiron/sqlx"
        _ "github.com/go-sql-driver/mysql"
        "github.com/go-redis/redis/v8"
)

type Post struct {
	ID           int       `db:"id"`
	UserID       int       `db:"user_id"`
	Imgdata      []byte    `db:"imgdata"`
	Body         string    `db:"body"`
	Mime         string    `db:"mime"`
	CreatedAt    time.Time `db:"created_at"`
	CommentCount int
	//Comments     []Comment
	//User         User
	CSRFToken    string
}

type Config struct {
        RedisAddr     string
        RedisPassword string
        RedisDB       int
}

var (
    rdb  *redis.Client
    conf = &Config{
            RedisAddr:     "localhost:6379",
            RedisPassword: "",
            RedisDB:       0,
    }
	db    *sqlx.DB
)

func main() {
        getPosts()
}

func getPosts() ([]Post, error) {
	rdb = redis.NewClient(&redis.Options{
    	Addr:     conf.RedisAddr,
    	Password: conf.RedisPassword,
    	DB:       conf.RedisDB,
    })

    ctx := rdb.Context()

    const queryCacheKey = "SELECT `id`, `user_id`, `imgdata`, `body`, `mime`, `created_at` FROM `posts` ORDER BY `created_at` DESC"

	// Redisにキャッシュがある場合はそれを返す
	cachedData, err := rdb.Get(ctx, queryCacheKey).Result()
	if err == nil {
		var posts []Post
		err = json.Unmarshal([]byte(cachedData), &posts)
		if err == nil {
			return posts, nil
		}
	}

	// Redisにキャッシュがない場合はクエリを実行して結果をキャッシュする
	posts, err := executeQuery()
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(posts)
	if err != nil {
		return nil, err
	}
	err = rdb.Set(ctx, queryCacheKey, jsonData, 10*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func executeQuery() ([]Post, error) {
	var posts []Post

        db, err := sqlx.Open("mysql","isuconp:isuconp@(localhost:3306)/isuconp?parseTime=true")
        if err != nil {
                log.Fatalf("failed to open database: %v", err)
        }
        defer db.Close()

	err = db.Select(&posts, "SELECT `id`, `user_id`, `body`, `mime`, `created_at` FROM `posts` ORDER BY `created_at` DESC")
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return posts, nil
}