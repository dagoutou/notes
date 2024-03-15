package opentelemetry

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.17.0"
	"log"
	"time"
	"yunqutech.gitlab.com/agilex/basis/plugin/logger"
)

func InitProvider(httpEndpoint string) func(context.Context) error {
	fmt.Println("httpEndpoint-----------------", httpEndpoint)
	ctx := context.Background()
	res, err := resource.New(context.Background(), resource.WithAttributes(
		semconv.ServiceName("gin-test"),
		semconv.ServiceVersion("1.0.0"),
		attribute.String("env", "dev"),
	),
		resource.WithOS(),
		resource.WithHost(),
	)
	a := res.String()
	fmt.Println(a)
	if err != nil {
		log.Fatal(err)
	}
	//初始化export
	var traceExporter *otlptrace.Exporter
	{
		ctx, cancel := context.WithTimeout(ctx, 100*time.Second)
		defer cancel()
		if httpEndpoint != "" {
			traceExporter, err = otlptracehttp.New(ctx, otlptracehttp.WithInsecure(), otlptracehttp.WithEndpoint(httpEndpoint))
			if err != nil {
				log.Fatal(err)
			}
		}
		bsp := sdk.NewBatchSpanProcessor(traceExporter)
		provider := sdk.NewTracerProvider(
			sdk.WithSampler(sdk.AlwaysSample()),
			sdk.WithResource(res),
			sdk.WithSpanProcessor(bsp),
		)
		otel.SetTracerProvider(provider)
		// 获取上下文
		bag := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
		otel.SetTextMapPropagator(bag)
		return provider.Shutdown
	}

}

type requestErr struct {
	error error
}

func (req *requestErr) Handle(err error) {
	logger.Errorf("TracerProvider error", logrus.Fields{"error": err})
}
