package engine

import (
	"github.com/d5/tengo/v2"
	"github.com/pkg/errors"
)

/*
构造函数脚本

输入
	input
输出
	初始化后的 domain 数据

example
	output = {
		consumer_id: input.consumer_id,
		seller_id: input.seller_id,
		amount: input.amount,
		commodity_id: input.commodity_id
	}
*/

type eventDomainConstructorExecutor struct {
	c *tengo.Compiled
}

func newEventDomainConstructExecutor(script string) (eventDomainConstructorExecutor, error) {
	s := tengo.NewScript([]byte(script))
	_ = s.Add("input", nil)
	_ = s.Add("output", nil)
	c, err := s.Compile()
	if err != nil {
		return eventDomainConstructorExecutor{}, err
	}
	return eventDomainConstructorExecutor{c}, nil
}

func (b eventDomainConstructorExecutor) exec(input tengoV) (tengoMAP, error) {
	b.c.Set("input", input)
	b.c.Set("output", nil)
	if err := b.c.Run(); err != nil {
		panic(err)
	}
	output := b.c.Get("output")
	m := output.Map()
	if m == nil{
		return nil, errors.New("eventDomainConstructorExecutor output not map")
	}
	return m, nil
}