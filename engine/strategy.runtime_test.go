// +build test

package engine

import (
	"fmt"
	"re/dal"
	"testing"
)

func initD(eventID int) runtimeContext{
	// 获取 event domain
	rel, err := dal.EventDomainGet(eventID)
	if err != nil {
		panic(err)
	}
	// 创建 context
	ctx := runtimeContext{
		root:    rel.DomainID,
		domains: map[int]domainStatus{},
	}
	// 初始化 event domain
	exec, err := newEventDomainConstructExecutor(rel.Constructor)
	if err != nil {
		panic(err)
	}
	var ds domainStatus
	{
		o, err := exec.exec(map[string]interface{}{
			"consumer_id":  1,
			"seller_id":    2,
			"amount":       100,
			"commodity_id": 3,
		})
		if err != nil {
			panic(err)
		}
		ds.factorCache = o
		ds.id = rel.DomainID
	}
	ctx.AddDomain(ds)
	fmt.Println("done")
	return ctx
}

func TestRuntime1(t *testing.T){
	fmt.Println(initD(1))
}