package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"net/http"
)

func main() {
	//ctx := context.Background()
	//shutdown := opentelemetry.InitProvider("127.0.0.1:4318")
	//defer func() {
	//	shutdown(ctx)
	//}()
	//
	////resource.StringDetector("")
	//g := gin.Default()
	//g.Use(OpenTelemetryMiddleware())
	//g.Use(RequestInfoMiddleware())
	//g.GET("/", func(c *gin.Context) {
	//
	//	// 创建 Span
	//	_, span := otel.Tracer("my-service").Start(ctx, "handle-request")
	//	defer span.End()
	//	reqInfo, _ := c.Get("request_info")
	//	requestInfo := reqInfo.(string)
	//	span.SetAttributes(
	//		attribute.String("query_parameters", requestInfo),
	//		attribute.String("tracerID", span.SpanContext().TraceID().String()),
	//		attribute.String("spanID", span.SpanContext().SpanID().String()),
	//	)
	//	//注册baggage
	//	var mpc = make(propagation.MapCarrier)
	//	mpc.Set("gouzi", "test")
	//	var pt propagation.TraceContext
	//	pt.Inject(ctx, mpc)
	//	otel.GetTextMapPropagator().Inject(ctx, mpc)
	//	err := errors.New("error test!")
	//	span.AddEvent("error", trace.WithAttributes(attribute.String("error message", err.Error()),
	//		attribute.String("error.type", fmt.Sprintf("%T", err)),
	//	))
	//
	//	// 执行操作...
	//	c.JSON(200, gin.H{"message": "Hello, OpenTelemetry with Gin!"})
	//})
	//g.GET("/hello", func(c *gin.Context) {
	//	c.JSON(200, gin.H{"message": "hello!"})
	//})
	//// 启动服务器
	//if err := g.Run(":8080"); err != nil {
	//	log.Fatal(err)
	//}
	tp := trace.NewTracerProvider()
	otel.SetTextMapPropagator(propagation.Baggage{})

	otel.SetTracerProvider(tp)
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// 从 HTTP 请求中提取追踪上下文
		b := otel.GetTextMapPropagator()
		fmt.Println(b)
		ctx := otel.GetTextMapPropagator().Extract(req.Context(), propagation.HeaderCarrier(req.Header))

		// 创建一个新的 Span，使用提取的上下文作为父上下文
		tracer := otel.Tracer("example")
		ctx, span := tracer.Start(ctx, "server-span")

		var a = make(propagation.MapCarrier)
		// 使用注入方法将 Span 数据写入响应头部
		otel.GetTextMapPropagator().Inject(ctx, &a)

		span.End()

		fmt.Fprint(w, "Hello, OpenTelemetry!")
	})
	http.ListenAndServe(":8080", nil)

}

//	func OpenTelemetryMiddleware() gin.HandlerFunc {
//		return func(context *gin.Context) {
//			ctx := context.Request.Context()
//			//创建span
//			ctx, span := otel.Tracer("gin-demo").Start(ctx, "http-request")
//			defer span.End()
//			span.SetAttributes(
//				attribute.String("tracerID", span.SpanContext().TraceID().String()),
//				attribute.String("spanID", span.SpanContext().SpanID().String()),
//			)
//			ctx = trace.ContextWithSpan(ctx, span)
//			context.Request = context.Request.WithContext(ctx)
//			context.Next()
//		}
//	}
func RequestInfoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		queryParams := c.Request.URL.Query()
		headers := c.Request.Header
		requestInfo := fmt.Sprintf("Query Parameters: %+v, Headers:%+v", queryParams, headers)
		c.Set("request_info", requestInfo)
		c.Next()
	}
}
