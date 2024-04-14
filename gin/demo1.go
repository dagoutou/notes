package main

//import (
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"time"
//)
//
//func middleware() (HandlerFunc func(*gin.Context)) {
//	return func(c *gin.Context) {
//
//		t := time.Now()
//		fmt.Println("中间件开始执行了")
//		// 设置变量到Context的key中，可以通过Get()取
//		c.Set("request", "中间件")
//		c.Next()
//		status := c.Writer.Status()
//		fmt.Println("中间件执行完毕", status)
//		t2 := time.Since(t)
//		fmt.Println("time:", t2)
//	}
//}
//func main() {
//	g := gin.Default()
//	g.Use(middleware())
//	{
//		g.GET("/ce", func(c *gin.Context) {
//			req, _ := c.Get("request")
//			fmt.Println("request:", req)
//			c.JSON(200, gin.H{"request": req})
//		})
//	}
//	g.Run()
//}
