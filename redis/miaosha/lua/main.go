package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	result, err := rdb.Eval(ctx, "return 1+1", nil).Result()
	if err != nil {
		return
	}
	fmt.Println(result)
	script := `local count = redis.call("get", KEYS[1])
if count == false then
    count = 0
else
    count = tonumber(count)
end
count = count + tonumber(ARGV[1])
redis.call("set", KEYS[1], count)
return count`
	re, err := rdb.Eval(ctx, script, []string{"mykey"}, 1).Result()
	if err != nil {
		return
	}
	fmt.Println(re)
	select {}
}
