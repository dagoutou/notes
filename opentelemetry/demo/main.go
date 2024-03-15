package main

import (
	"context"
	"demo/provider"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type database struct {
	Username          string `json:"username"`
	Password          string `json:"password"`
	Version           string `json:"version"`
	Type              string `json:"type"`
	URL               string `json:"url"`
	Key               string `json:"key"`
	SecurityMechanism string `json:"securityMechanism"`
}
type params struct {
	SQL    string `json:"sql"`
	Schema string `json:"schema"`
}
type data struct {
	IP            string `json:"ip"`
	Name          string `json:"name"`
	MonitorStatus bool   `json:"monitor_status"`
	AuditStatus   bool   `json:"audit_status"`
}
type instance struct {
	ID       int           `json:"id"`
	Category []string      `json:"category"`
	Usefor   []string      `json:"usefor"`
	Account  []interface{} `json:"account"`
	Data     data          `json:"data"`
	Status   bool          `json:"status"`
	Network  string        `json:"network"`
}
type queryRequest struct {
	Database database `json:"database"`
	Params   params   `json:"params"`
}

func dbQuery(class string, data queryRequest, inst instance) {}
func main() {
	ctx := context.Background()
	shoutdown := provider.InitTraceProvider("127.0.0.1:4318")
	defer func() {
		shoutdown(ctx)
	}()
	//tracer := otel.Tracer("my_test")
	//instanceStruct := instance{
	//	ID:       123,
	//	Category: []string{"category1", "category2"},
	//	Usefor:   []string{"usefor1", "usefor2"},
	//	Account:  []interface{}{"account1", "account2"},
	//	Data: data{
	//		IP:            "192.168.1.1",
	//		Name:          "John Doe",
	//		MonitorStatus: true,
	//		AuditStatus:   false,
	//	},
	//	Status:  true,
	//	Network: "network1",
	//}
	//var buff bytes.Buffer
	//en := json.NewEncoder(&buff)
	//err := en.Encode(instanceStruct)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//str1 := buff.String()
	//size1 := unsafe.Sizeof(str1)
	//fmt.Println(size1)
	//marshal, err := json.Marshal(instanceStruct)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//str := string(marshal)
	//size2 := unsafe.Sizeof(str)
	//fmt.Println(size2)
	//ctx, span := tracer.Start(ctx, "example-span",
	//	trace.WithAttributes(
	//		attribute.String("icarus_request_parameters", str),
	//	),
	//)
	//defer span.End()

	tracer := otel.Tracer("example")

	ctx, parentSpan := tracer.Start(ctx, "parent-span")

	// 在不同的方法中传递 Trace 的上下文
	method1(ctx)
	method2(ctx)

	// 结束父 Span
	parentSpan.End()

	fmt.Println("Spans completed")
}
func method1(ctx context.Context) {
	ctx, span := otel.Tracer("example1").Start(ctx, "method1")
	defer span.End()
}
func method2(ctx context.Context) {
	ctx, span := otel.Tracer("example").Start(ctx, "method2", trace.WithAttributes(attribute.String("icarus_request_parameters", "aafe")))
	defer span.End()
}
