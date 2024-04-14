package provider

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.17.0"
	"log"
	"time"
)

func InitTraceProvider(url string) func(context.Context) error {
	ctx := context.Background()
	res, err := resource.New(ctx, resource.WithOS(),
		resource.WithHostID(), resource.WithAttributes(
			semconv.ServiceName("test"),
			semconv.ServiceVersion("v1.0.0"),
			attribute.String("env", "dev"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	var httpExporter *otlptrace.Exporter
	{
		ctx, cancel := context.WithTimeout(ctx, 100*time.Second)
		defer cancel()
		if httpExporter, err = otlptracehttp.New(ctx, otlptracehttp.WithInsecure(), otlptracehttp.WithEndpoint(url)); err != nil {
			log.Fatal(err)
		}
		nsp := trace.NewBatchSpanProcessor(httpExporter)
		provider := trace.NewTracerProvider(
			trace.WithResource(res),
			trace.WithSampler(trace.AlwaysSample()),
			trace.WithSpanProcessor(nsp),
		)
		otel.SetTracerProvider(provider)
		otel.SetTextMapPropagator(propagation.Baggage{})
		return provider.Shutdown
	}
}
