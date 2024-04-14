package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	//data, err := json.Marshal(map[string]string{"key1": "val1", "key2": "val2"})
	//if err != nil {
	//	log.Fatal("json.Marshal error", data)
	//}
	//resp, err := http.Post(":8091", "application/json", bytes.NewReader(data))
	//if err != nil {
	//	log.Fatal("http.Post error", data)
	//}
	//defer resp.Body.Close()
	//respBody, _ := io.ReadAll(resp.Body)
	//fmt.Printf("resp: %s", respBody)

	router := gin.Default()
	router.GET("/cookie", func(c *gin.Context) {
		_, err := c.Cookie("gin_cookie") // 获取Cookie
		if err != nil {

			// 设置Cookie
			c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
		}
		a, err := c.Cookie("gin_cookie") // 获取Cookie
		fmt.Printf("Cookie value: %s \n", a)
	})

	router.Run(":8081")
}
