package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type NameInputPageData struct {
	Name string
}

// string 写入（增加指定值）、读取、删除、获取长度,追加，设置多个mset、位操作
// hash 写入、读取、写入多个、读取全部；获取全部的键、判断是否存在，只获取键，只获取值，获取字段的长度
// list 写入、读取、获取长度、遍历
// set 添加、删除、修改；获取全部,获取交集、并、独立
func main() {
	g := gin.Default()
	g.LoadHTMLGlob("/Users/wanglili/study/notes/redis/demo1_connection/templates/*")
	g.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	g.POST("/submit", func(c *gin.Context) {
		name := c.PostForm("name")
		ctx := context.Background()
		redisClient := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
		_, err := redisClient.Get(ctx, name).Result()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"name": "姓名未注册"})
			return
		}
		err = redisClient.Set(ctx, name, name, 0).Err()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"name": "姓名未注册"})
		}

		c.HTML(http.StatusOK, "submitted.html", gin.H{"Name": name})
	})
	g.GET("/redis", getConnection)
	g.Run(":8080")
}
func getConnection(c *gin.Context) {

}
