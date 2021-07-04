package engine

import (
	"encoding/json"
	"fmt"
	"testing"
)

/*
	下订单（事件初始域）
		<*>消费金额
			(1) 和上一次消费金额相同

	用户
		<*>信用分
			(2) 高于百分之 {} 的用户
		<*>本月消费金额列表
			(3) 平均值小于 {}

	商品
		<*>商品月销量
			(4) 大于等于 {} 且小于 {}

	testData 调用上列的 1 2 3 4 四种操作符
 */
const testData = `
{
	"rules": [
		{"path":[], "factor":"amount", "operation": 1, "args": []},
		{"path":[2], "factor":"CreditScore", "operation": 2, "args": [{"typ":0, "value":20}]},
		{"path":[2], "factor":"ConsumptionAmountListMonth", "operation": 3, "args": [{"typ":0, "value":150}]},
		{"path":[3], "factor":"monthly_sales", "operation": 4, "args": [{"typ":0, "value":20}, {"typ":0, "value":280}]}
	],
	"combination": "1&2&3&4"
}
`

func Test1(t *testing.T) {
	var s Strategy
	json.Unmarshal([]byte(testData), &s)

	fmt.Println(strategyExec(s, map[string]interface{}{
		"consumer_id":  1,
		"seller_id":    2,
		"amount":       200, // 和上次消费金额相同；上次消费金额暂时写死在 db/test.go::本月消费金额列表 = [200,100,150]
		"commodity_id": 3,
	}))
}
