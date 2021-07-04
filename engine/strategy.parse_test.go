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
		<*>本月消费金额列表 = [200,100,150]
			(3) 平均值小于 {}

	商品
		<*>商品月销量 = 188
			(4) 大于等于 {} 且小于 {}

	testData 调用上列的 1 2 3 4 四种操作符
 */
const testData = `
{
	"rules": [
		{"path":[], "factor":"amount", "operation": 1, "args": []},
		{"path":["user domain1"], "factor":"CreditScore", "operation": 2, "args": [{"typ":0, "value":20}]},
		{"path":["user domain2"], "factor":"ConsumptionAmountListMonth", "operation": 3, "args": [{"typ":0, "value":151}]},
		{"path":["commodity domain"], "factor":"monthly_sales", "operation": 4, "args": [{"typ":0, "value":20}, {"typ":0, "value":180}]}
	],
	"combination": "1&2&3&4"
}
`

func TestParse(t *testing.T) {
	var s Strategy
	json.Unmarshal([]byte(testData), &s)
	fmt.Println(s)
}
