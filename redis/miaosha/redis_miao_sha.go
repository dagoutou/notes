package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

const (
	BeginTimestamp = 1704067200
	CountBit       = 32
)

var RedisClient = redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", Password: "", DB: 0})

func generatorID(keyPrefix string) (int64, error) {
	//生成时间戳
	timeNow := time.Now().Unix()
	timeStamp := timeNow - BeginTimestamp

	//生成自增序列
	ctx := context.Background()

	timeFormat := time.Now().Format("2006-01-02")
	key, err := RedisClient.Incr(ctx, "incr:"+keyPrefix+":"+timeFormat).Result()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	//拼接ID
	return timeStamp<<32 | key, nil
}
