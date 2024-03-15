package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

type Result struct {
	PromQL   string `json:"promql"`
	ErrorMsg string `json:"error_msg"`
}

func main() {
	// 从命令行参数读取 PromQL 和 instance ID

	promQL := "db_stats_gauge{NAME=\"logons current\" }"
	instanceID := "2"

	result := Result{
		PromQL:   "",
		ErrorMsg: "",
	}

	// 解析 PromQL 表达式
	expr, err := parser.ParseExpr(promQL)
	if err != nil {
		result.ErrorMsg = fmt.Sprintf("Invalid PromQL: %s", err.Error())
		jsonResult, _ := json.Marshal(result)
		fmt.Println(string(jsonResult))
		os.Exit(1)
	}

	// 遍历表达式的节点，检查并添加 instance 标签
	parser.Inspect(expr, func(node parser.Node, path []parser.Node) error {
		if n, ok := node.(*parser.VectorSelector); ok {
			hasInstanceLabel := false
			for _, label := range n.LabelMatchers {
				if label.Name == "instance" {
					hasInstanceLabel = true
					break
				}
			}
			if !hasInstanceLabel {
				n.LabelMatchers = append(n.LabelMatchers, &labels.Matcher{
					Type:  labels.MatchEqual,
					Name:  "instance",
					Value: instanceID,
				})
			}
		}
		return nil
	})

	result.PromQL = expr.String()
	jsonResult, _ := json.Marshal(result)
	fmt.Println(string(jsonResult))
}
