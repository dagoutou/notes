package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"runtime"
)

var (
	myMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "my_metric",
			Help: "This is my custom metric",
		},
		[]string{"label1", "label2"}, // 根据您的需求添加标签
	)
)

func init() {
	prometheus.MustRegister(myMetric)
}

func main() {
	//http.Handle("/metrics", promhttp.Handler())
	//err := http.ListenAndServe(":8000", nil)
	//if err != nil {
	//	panic(err)
	//}
	fmt.Println(runtime.GOARCH)

}
