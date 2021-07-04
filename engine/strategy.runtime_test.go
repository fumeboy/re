// +build test

package engine

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	var s Strategy
	json.Unmarshal([]byte(testData), &s)

	fmt.Println(strategyExec(s, map[string]interface{}{
		"consumer_id":  1,
		"seller_id":    2,
		"amount":       201, // 和上次消费金额相同；上次消费金额暂时写死在 db/test.go::本月消费金额列表 = [200,100,150]
		"commodity_id": 3,
	}))
}
