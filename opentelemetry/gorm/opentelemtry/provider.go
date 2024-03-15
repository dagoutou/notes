package opentelemetry

import (
	"context"
	"fmt"
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
)

func InitProvider(httpEndpoint string) func(context.Context) error {
	fmt.Println("httpEndpoint-----------------", httpEndpoint)
	ctx := context.Background()
	res, err := resource.New(context.Background(), resource.WithAttributes(
		semconv.ServiceName("gorm-test"),
		semconv.ServiceVersion("1.0.0"),
		attribute.String("env", "dev"),
	),
		resource.WithOS(),
		resource.WithHost(),
	)
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
		otel.SetTextMapPropagator(propagation.Baggage{})
		return provider.Shutdown
	}

}
