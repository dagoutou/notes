package bus

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type Bus interface {
	Sum(ctx context.Context, a, b int) int
	Product(ctx context.Context, a, b int) int
}
type bus struct {
	trace trace.Tracer
}

func NewBus(tracer trace.Tracer) Bus {
	return &bus{
		trace: tracer,
	}
}
func (bus *bus) Sum(ctx context.Context, a, b int) int {
	c := a + b
	trace.WithAttributes(attribute.Int("sum c:", c))
	ctx, span := bus.trace.Start(ctx, "Sum", trace.WithAttributes(attribute.Int("c", c)))
	defer span.End()
	<-time.After(time.Millisecond * 100)
	return c

}
func (bus *bus) Product(ctx context.Context, a, b int) int {
	c := a * b
	trace.WithAttributes(attribute.Int("Product c:", c))
	//ctx, span := bus.trace.Start(ctx, "Product", trace.WithAttributes(attribute.Int("c", c)))
	//defer span.End()
	//<-time.After(time.Millisecond * 200)
	return c
}
