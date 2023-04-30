package main

import (
        "fmt"
        "os"
        "log"
        "strconv"
        _ "github.com/go-sql-driver/mysql"
        "github.com/jmoiron/sqlx"
        "context"
        "github.com/redis/go-redis/v9"
        "json"
)

//
var (
        db    *sqlx.DB
)

//
func main() {
        host := os.Getenv("ISUCONP_DB_HOST")
        if host == "" {
                host = "localhost"
        }
        port := os.Getenv("ISUCONP_DB_PORT")
        if port == "" {
                port = "3306"
        }
        _, err := strconv.Atoi(port)
        if err != nil {
                log.Fatalf("Failed to read DB port number from an environment variable ISUCONP_DB_PORT.\nError: %s", err.Error())
        }
        user := os.Getenv("ISUCONP_DB_USER")
        if user == "" {
                user = "root"
        }
        password := os.Getenv("ISUCONP_DB_PASSWORD")
        dbname := os.Getenv("ISUCONP_DB_NAME")
        if dbname == "" {
                dbname = "isuconp"
        }

        dsn := fmt.Sprintf(
                "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
                user,
                password,
                host,
                port,
                dbname,
        )

        db, err = sqlx.Open("mysql", dsn)
        if err != nil {
                log.Fatalf("Failed to connect to DB: %s.", err.Error())
        }
        defer db.Close()


        //
        searchId := 5
        redisClient := ExampleNewClient()

        nameInRedis, err := redisClient.Get(ctx, "name_"+strconv.Itoa(searchId)).Result()
        if err != nil {
                fmt.Println(err)
        } else if err == redis.Nil {
                db, err = sqlx.Open("mysql", dsn)
                if err != nil {
                        panic(err.Error())
                }
                defer db.Close()
                stmtOut, err := db.Prepare("SELECT `id`, `user_id`, `body`, `mime`, `created_at` FROM `posts` WHERE `id` = 1234;")
                if err != nil {
                        panic(err.Error())
                }
                defer stmtOut.Close()
                rows, err := stmtOut.Query(searchId)
                if err != nil {
                        panic(err.Error())
                }
                numRows = 0
                for rows.Next() {
                var nameInSQL string
                err = rows.Scan(&nameInSQL)
                if err != nil {
                        panic(err.Error())
                }
                fmt.Printf("name is %s\n", nameInSQL)
                numRows = numRows + 1
                }
                if numRows == 0 {
                        fmt.Printf("corresponding name is not found\n")
                } else {
                        fmt.Printf("name is %s\n", nameInRedis)
                }
        }
}

//
var (
        numRows int
        ctx = context.Background()
        rdb *redis.Client
)

//
func ExampleNewClient() *redis.Client  {
        rdb := redis.NewClient(&redis.Options{
                Addr:     "localhost:6379", // use default Addr
                Password: "",               // no password set
                DB:       0,                // use default DB
        })
        pong, err := rdb.Ping(ctx).Result()
        fmt.Println(pong, err)
        return rdb
}

//
def fetch(sql):
    res = db.Prepare(sql)
    if res:
        print("Cache Hit")
        return json.loads(res)
    res = Database.query(sql)
    ### setexはTTL付きでデータをstring型でsetする
    Cache.setex(sql, TTL, json.dumps(res))
    print("Cache Write")
    return res