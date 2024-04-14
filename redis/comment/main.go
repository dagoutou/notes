package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Comments struct {
	ID         int
	Content    string
	CreateTime time.Time
}
type session struct {
	data map[string]interface{}
}
type sessionStore struct {
	sessions map[string]*session
}

func main() {
	g := gin.New()
	var seStore = new(sessionStore)
	redsClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	dsn := "root:12345678@tcp(127.0.0.1:3306)/notes?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("db connection error", err)
	}
	type tbUser struct {
		Id       int    `json:"id"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	g.GET("/addComment", func(c *gin.Context) {
		for i := 1; i <= 10; i++ {
			var com = Comments{
				Content:    fmt.Sprintf("content"),
				CreateTime: time.Now(),
			}
			if err = db.Create(&com).Error; err != nil {
				log.Fatal("create error", err)
			}
			ctx := context.Background()
			a := redis.Z{
				Score:  60,
				Member: "a",
			}
			b := redis.Z{
				Score:  70,
				Member: "b",
			}
			cc := redis.Z{
				Score:  80,
				Member: "c",
			}
			var sr []redis.Z
			sr = append(sr, a, b, cc)
			err = redsClient.ZAdd(ctx, "scores", sr...).Err()
			if err != nil {
				log.Println(err)
				return
			}
			sw := redis.ZRangeBy{
				Min:    "60",
				Max:    "70",
				Offset: 0,
				Count:  0,
			}
			scores, _ := redsClient.ZRangeByScoreWithScores(ctx, "scores", &sw).Result()
			c.JSON(http.StatusOK, scores)
			fmt.Println(redsClient.ZScore(ctx, "scores", "a").Val())

			err = redsClient.LPush(ctx, "comments", fmt.Sprintf("%d:%s:%s", com.ID, com.Content, com.CreateTime)).Err()
			if err != nil {
				return
			}
			err = redsClient.LTrim(ctx, "comments", 0, 49).Err()
			if err != nil {
				return
			}
		}
	})
	g.GET("/getComment", func(c *gin.Context) {
		ctx := context.Background()
		length := redsClient.LLen(ctx, "comments").Val()
		result, err := redsClient.LRange(ctx, "comments", 0, length-1).Result()
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, result)
	})
	g.POST("/sendCode", func(c *gin.Context) {
		var mp = make(map[string]string)
		phone := c.Query("phone")
		rand.Seed(time.Now().UnixNano())
		code := generateVerificationCode()
		mp["phone"] = phone
		mp["code"] = code
		cookieName := fmt.Sprintf("code_%s", phone)
		_, err = c.Cookie(cookieName) // 获取Cookie
		if err != nil {
			// 设置Cookie
			c.SetCookie("1111", "33333", 3600, "/", "localhost", false, true)
		}
		codeCookie, err := c.Cookie("1111")
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.New("获取验证码错误！").Error())
			return
		}

		c.JSON(http.StatusOK, codeCookie)
	})
	g.GET("/getCode", func(c *gin.Context) {
		phone := c.Query("phone")
		code := c.Query("code")
		cookieName := fmt.Sprintf("code_%s", phone)
		codeCookie, err := c.Cookie(cookieName)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.New("获取验证码错误！").Error())
			return
		}
		if codeCookie == code {
			type user struct {
				name string
				age  int
			}
			var u = user{
				name: cookieName,
				age:  20,
			}
			var mp = make(map[string]interface{})
			var mp1 = make(map[string]*session)
			mp["users"] = u
			var se = session{data: mp}
			mp1["cookieName"] = &se
			seStore.sessions = mp1
			//存储信息到redis中
			c.JSON(http.StatusOK, seStore)
		}

	})
	g.Run(":8080")
}
func generateVerificationCode() string {
	const charset = "0123456789"
	code := make([]byte, 4)
	for i, _ := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
func getLatestComments(rdb *redis.Client) ([]string, error) {
	ctx := context.Background()
	result, err := rdb.LRange(ctx, "comments", 0, 100).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}
func a() {
	//查询数据库 1m
	//redis -->  go 查询数据库 --> redis

}
