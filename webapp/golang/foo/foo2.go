package main

import (
        "encoding/json"
        "fmt"
        "log"
        "net/http"
        "time"

        "github.com/jmoiron/sqlx"
        _ "github.com/go-sql-driver/mysql"
        "github.com/go-redis/redis/v8"
        _ "github.com/lib/pq"
)

type Config struct {
        RedisAddr     string
        RedisPassword string
        RedisDB       int
}

type User struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
}

var (
        rdb  *redis.Client
        conf = &Config{
                RedisAddr:     "localhost:6379",
                RedisPassword: "",
                RedisDB:       0,
        }
)

func main() {
        initRedis()

        db, err := sqlx.Open("mysql","isuconp:isuconp@(localhost:3306)/isuconp")
        if err != nil {
                log.Fatalf("failed to open database: %v", err)
        }
        defer db.Close()

        http.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
                id := r.URL.Path[len("/user/"):]
                user, err := getUser(id, db)
                if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }
                json.NewEncoder(w).Encode(user)
        })

        log.Fatal(http.ListenAndServe(":8090", nil))
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
			Addr:     conf.RedisAddr,
			Password: conf.RedisPassword,
			DB:       conf.RedisDB,
	})
}

func getUser(id string, db *sqlx.DB) (*User, error) {
	ctx := rdb.Context()
	cacheKey := fmt.Sprintf("user:%s", id)

	// Check if the data exists in the Redis Cache.
	data, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
			user := &User{}
			err = json.Unmarshal([]byte(data), user)
			if err == nil {
					return user, nil
			}
	}

	// If not in cache, get data from the database.
	row := db.QueryRow("SELECT id, account_name FROM users WHERE id = CONCAT(?, '%')", id)
	user := &User{}
	err = row.Scan(&user.ID, &user.Name)
	if err != nil {
			return nil, err
	}

	// Store the data in Redis Cache.
	dataBytes, err := json.Marshal(user)
	if err == nil {
			rdb.Set(ctx, cacheKey, dataBytes, time.Minute*5) // Set cache expiration to 5 minutes.
	}

	return user, nil
}