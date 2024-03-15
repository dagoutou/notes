package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	opentelemetry "gorm-opentelemetry/opentelemtry"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/logging/logrus"
	"gorm.io/plugin/opentelemetry/tracing"
	"log"
	"time"
)

type AlertConfig struct {
	ID           int            `json:"id" gorm:"primary_key;not null"`
	Global       datatypes.JSON `json:"global" gorm:"type:jsonb;not null;default:'{}'" binding:"required"`        // 全局配置
	Route        datatypes.JSON `json:"route" gorm:"type:jsonb;not null;default:'{}'" binding:"required"`         // 路由
	Templates    datatypes.JSON `json:"templates" gorm:"type:jsonb;not null;default:'[]'" binding:"required"`     // 模版
	InhibitRules datatypes.JSON `json:"inhibit_rules" gorm:"type:jsonb;not null;default:'[]'" binding:"required"` // 内部规则
}

func (AlertConfig) TableName() string {
	return "athena_alert_configs"
}
func main() {
	ctx := context.Background()
	//注册全局provider
	shutdown := opentelemetry.InitProvider("127.0.0.1:4318")
	defer func() {
		shutdown(context.Background())
	}()

	logger := logger.New(
		logrus.NewWriter(),
		logger.Config{
			SlowThreshold: time.Millisecond,
			LogLevel:      logger.Warn,
			Colorful:      false,
		},
	)

	// 设置数据库连接配置
	dsn := "host=10.0.0.108 user=yunqu password=YunquTech01*@ dbname=devops port=5432 sslmode=disable"

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger})
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}
	tracer := otel.Tracer("gorm.io/plugin/opentelemetry")
	ctx, span := tracer.Start(context.Background(), "root")
	defer span.End()
	var res AlertConfig
	if err := db.WithContext(ctx).Model(&AlertConfig{}).Find(&res).Error; err != nil {
		log.Println(err)
		return
	}
	getConfig(ctx, db)
	fmt.Println(res)
}

type Config struct {
	ID       int            `json:"id" gorm:"primary_key;not null"`
	Name     string         `json:"name" gorm:"type:varchar(50);not null;default:''" validate:"required"`                               // 监控模板名称
	Type     string         `json:"type" gorm:"type:varchar(7);not null;default:''" enums:"alerts,metrics,notices" validate:"required"` // 监控模板类型 alerts告警规则 metrics监控指标 notices通知策略
	Category datatypes.JSON `json:"category" gorm:"type:jsonb;not null;default:'[]'" validate:"required"`                               // 资源类型数据库or系统
	Usefor   datatypes.JSON `json:"usefor" gorm:"type:jsonb;not null;default:'[]'" validate:"required"`                                 // 使用者
	Template bool           `json:"template" gorm:"type:boolean;not null"`                                                              // 是否为模板
	Status   bool           `json:"status" gorm:"type:boolean;not null"`                                                                // 当前使用状态
	Source   int            `json:"source" gorm:"not null;default:0"`                                                                   // 区别是sql审核还是鹰眼
}

func (Config) TableName() string {
	return "athena_configs"
}
func getConfig(ctx context.Context, db *gorm.DB) {
	var resq Config
	if err := db.WithContext(ctx).Model(&Config{}).Find(&resq).Error; err != nil {
		log.Println(err)
		return
	}
	sp := trace.SpanFromContext(ctx)
	sp.End()
	fmt.Println(resq)
}
