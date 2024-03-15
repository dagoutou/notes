package main

import (
	"context"
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"jaeger/bus"
	"net/http"

	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var exporterHttpEndpoint = flag.String("export-http-endpoint", "127.0.0.1:4318", "")
var exporterGRPCEndpoint = flag.String("export-grpc-endpoint", "127.0.0.1:4317", "")

func main() {
	//trace provider 追踪提供程序，启动程序初始化一次
	//trace 有 tracer provider 初始化
	//由trace开启span
	flag.Parse()
	shutdown := initProvider(*exporterHttpEndpoint, *exporterGRPCEndpoint)
	defer func() {
		shutdown(context.Background())
	}()

	tracer := otel.Tracer("main-trace")
	ctx := context.Background()

	ctx, span := tracer.Start(ctx, "main", trace.WithAttributes(attribute.String("package", "main")))
	defer span.End()

	span.RecordError(errors.New("异常信息"), trace.WithTimestamp(time.Now()))
	span.SetStatus(codes.Error, "报错了")

	r := gin.New()
	// 使用 gintrace 中间件进行集成
	r.Use(otelgin.Middleware("your-service-name"))
	// 定义路由
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, OpenTelemetry with Gin!")
	})
	r.Run(":8080")
	busTracer := otel.Tracer("bus-tracer")
	busObj := bus.NewBus(busTracer)
	for i := 0; i < 5; i++ {
		busObj.Sum(ctx, i, i+1)
		busObj.Product(ctx, i, i+1)
		<-time.After(time.Second * 3)
	}

}
func test(ctx context.Context) {

}
func initProvider(httpEndpoint, grpcEndpoint string) func(context.Context) error {
	ctx := context.Background()
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName("traces-basic"),
			semconv.ServiceVersion("1.0.0"),
			attribute.String("env", "dev"),
		), resource.WithOS(), resource.WithHost())
	if err != nil {
		log.Fatal(err)

	}
	//初始化export
	var traceExporter *otlptrace.Exporter
	{
		ctx, cancel := context.WithTimeout(ctx, 1000*time.Second)
		defer cancel()
		if grpcEndpoint != "" {
			conn, err := grpc.DialContext(ctx, grpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
			if err != nil {
				log.Fatal(err)

			}
			traceExporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithGRPCConn(conn))
			if err != nil {
				log.Fatal(err)
			}
		} else if httpEndpoint != "" {
			traceExporter, err = otlptracehttp.New(ctx, otlptracehttp.WithInsecure(), otlptracehttp.WithEndpoint(httpEndpoint))
			if err != nil {
				log.Fatal(err)

			} else {
				log.Fatal("没有任何exporter")
			}
		}
		bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
		traceProvider := sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithResource(res),
			sdktrace.WithSpanProcessor(bsp),
		)
		otel.SetTracerProvider(traceProvider)
		otel.SetTextMapPropagator(propagation.Baggage{})
		return traceProvider.Shutdown
	}
}
