package main
import (
	"context"
	"github.com/redis/go-redis/v9"
)

func main() {
    hoge()
}

// redis
var (
	ctx = context.Background()
	rdb *redis.Client
)

type Data struct {
	key   string
	value string
}

func hoge() {
	c := ExampleNewClient()
	d := Data{
		key:   "key1",
		value: "value1",
	} 
	var ctx = context.Background()
	if err := c.Set(ctx, d.key, d.value, 0).Err(); err != nil {
		panic(err)
	}
	val, err := c.Get(ctx, d.key).Result()
	switch {
	case err == redis.Nil:
		panic("key does not exist")
	case err != nil:
		panic(err)
	case val == "":
		panic("value is empty")
	}
	fmt.Println(d.key, val)
}

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