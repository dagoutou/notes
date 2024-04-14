package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	lock2 "notes/redis/miaosha/lock"
	"strconv"
	"sync"
	"time"
)

var DB *gorm.DB
var mu sync.Mutex

func main() {
	dsn := "root:12345678@tcp(127.0.0.1:3306)/notes?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DB = db
	g := gin.Default()
	g.POST("/voucher-order/seckill/:id", seckillVoucherOrder)
	g.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, "success")
	})
	g.Run(":8081")
}

type TbVoucher struct {
	Id          int       `json:"id"`
	ShopId      int       `json:"shopId"`
	Title       string    `json:"title"`
	SubTitle    string    `json:"subTitle"`
	Rules       string    `json:"rules"`
	PayValue    int       `json:"payValue"`
	ActualValue int       `json:"actualValue"`
	Status      int       `json:"status"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}
type TbVoucherOrder struct {
	Id        int64 `json:"id"`
	UserId    int   `json:"userId"`
	VoucherId int   `json:"voucherId"`
}
type seckillVoucherInfo struct {
	VoucherId int       `json:"voucherId"`
	BeginTime time.Time `json:"beginTime"`
	EndTime   time.Time `json:"endTime"`
	Stock     int       `json:"stock"`
}

func seckillVoucherOrder(c *gin.Context) {
	var tv seckillVoucherInfo
	id := c.Param("id")
	atoi, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.New("ID转换失败！"))
		return
	}
	if err := DB.Table("tb_seckill_voucher").Select(` begin_time, end_time, stock`, id).Where("voucher_id = ?", atoi).Find(&tv).Error; err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	now := time.Now().Unix()
	if tv.BeginTime.Unix() > now {
		log.Println("秒杀未开始！")
		c.JSON(http.StatusBadRequest, errors.New("秒杀未开始！"))
		return
	}
	if tv.EndTime.Unix() < now {
		log.Println("秒杀已经结束了！")
		c.JSON(http.StatusBadRequest, errors.New("秒杀已经结束了！"))
		return
	}

	orderID, err := generatorID("order")
	if err != nil {
		log.Println("生成订单ID失败！")
		c.JSON(http.StatusBadRequest, errors.New("生成订单ID失败！"))
		return
	}

	tx := DB.Begin()
	if tv.Stock < 1 {
		log.Println("库存不足！")
		c.JSON(http.StatusBadRequest, errors.New("库存不足！"))
		return
	}
	var cnt int64
	ctx := context.Background()
	//bol := RedisClient.SetNX(ctx, "order:user:1", 4, 10*time.Second).Val()
	//if !bol {
	//	log.Println("SetNX error")
	//	return
	//}
	//pid := os.Getpid()
	//pidStr := strconv.Itoa(pid)
	key := "lock:1"
	var lock = lock2.Lock{
		Ctx:      ctx,
		Client:   RedisClient,
		LockName: key,
	}
	if err = lock.TryLock(10 * time.Second); err != nil {
		log.Println("lock error")
		return
	}
	if err = tx.Table("tb_voucher_order").Where("voucher_id = ? and user_id = ? ", atoi, 1).Count(&cnt).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if cnt > 0 {
		log.Println("该用户已经抢购过了不能重复抢购！")
		c.JSON(http.StatusBadRequest, errors.New("该用户已经抢购过了不能重复抢购！"))
		return
	}
	if err = tx.Table("tb_seckill_voucher").Debug().Where("voucher_id = ? and stock > 0", atoi).Update("stock", gorm.Expr("stock -1")).Error; err != nil {
		tx.Rollback()
		log.Println("库存不足！")
		c.JSON(http.StatusBadRequest, errors.New("库存不足！"))
		return
	}

	var order TbVoucherOrder

	order.UserId = 1
	order.VoucherId = atoi
	order.Id = orderID
	if err = tx.Table("tb_voucher_order").Create(&order).Error; err != nil {
		tx.Rollback()
		log.Println("保存订单信息失败！")
		c.JSON(http.StatusBadRequest, errors.New("保存订单信息失败！"))
		return
	}
	tx.Commit()
	if err = lock.UnLocke(); err != nil {
		log.Println("unlock error")
		return
	}
	c.JSON(http.StatusOK, "order")
}
