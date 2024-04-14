package lock

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strconv"
	"time"
)

type ILock interface {
	TryLock(timeoutSec time.Duration) error
	UnLocke() error
}
type Lock struct {
	Ctx      context.Context
	Client   *redis.Client
	LockName string
}

func (l *Lock) TryLock(timeoutSec time.Duration) error {
	pid := os.Getpid()
	pidStr := strconv.Itoa(pid)
	bol := l.Client.SetNX(l.Ctx, l.LockName, pidStr, timeoutSec).Val()
	if !bol {
		err := errors.New("重复抢购商品！")
		return err
	}
	return nil
}
func (l *Lock) UnLocke() error {
	pid := os.Getpid()
	pidStr := strconv.Itoa(pid)
	result, err := l.Client.Get(l.Ctx, l.LockName).Result()
	if err != nil {
		log.Println("获取key信息失败")
		return err
	}
	if result == pidStr {
		if err := l.Client.Del(l.Ctx, l.LockName).Err(); err != nil {
			return err
		}
	}
	return nil
}
