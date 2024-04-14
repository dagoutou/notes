package main

import (
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
	"net/http"
)

func initTracer() {
	// 创建Jaeger导出器
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint())
	if err != nil {
		log.Fatalf("failed to create Jaeger exporter: %v", err)
	}

	// 创建简单的采样器
	// 这里使用 AlwaysSample，实际生产环境中可能需要更复杂的采样策略
	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(bsp))

	// 注册全局跟踪器
	otel.SetTracerProvider(tp)
}

func main() {
	// 初始化跟踪器
	initTracer()

	// 创建一个HTTP处理程序，添加跟踪代码
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		tracer := otel.Tracer("example-tracer")
		ctx, span := tracer.Start(req.Context(), "example-operation")
		defer span.End()

		// 模拟一些处理
		// 这里可以添加你的业务逻辑
		b := 10 + 30
		<-ctx.Done()
		fmt.Println(b)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	})

	// 启动HTTP服务器
	log.Fatal(http.ListenAndServe(":8080", nil))
}
