package engine

import (
	"github.com/d5/tengo/v2"
	"github.com/pkg/errors"
)

/*
输入
	input （from domain）
输出
	to domain 构造函数的入参

example:
	output = {user_id: factor(this, "consumer_id")}
*/

type bridgeScriptExecutor struct {
	c *tengo.Compiled
}

func newBridgeScriptExecutor(script string) (bridgeScriptExecutor, error) {
	s := tengo.NewScript([]byte(script))
	presetTangoFn(s)
	_ = s.Add("this", nil)
	_ = s.Add("output", nil)
	c, err := s.Compile()
	if err != nil {
		return bridgeScriptExecutor{}, err
	}
	return bridgeScriptExecutor{c}, nil
}

func (b bridgeScriptExecutor) exec(ctx runtimeContext, domain_id int) (tengoMAP, error) {
	ctx.ImportTangoFn(b.c)
	b.c.Set("this", domain_id)
	b.c.Set("output", nil)
	if err := b.c.Run(); err != nil {
		panic(err)
	}
	output := b.c.Get("output")
	m := output.Map()
	if m == nil{
		return nil, errors.New("bridgeScriptExecutor output not map")
	}
	return m, nil
}
