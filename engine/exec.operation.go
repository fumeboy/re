package engine

import (
	"fmt"
	"github.com/d5/tengo/v2"
	"github.com/pkg/errors"
)

/*
构造函数脚本

输入
	因子数据
	参数列表 args
输出
	bool 值

example
	output = true
*/


type operationExecutor struct {
	c *tengo.Compiled
}

func newOperationExecutor(script string) (operationExecutor, error) {
	s := tengo.NewScript([]byte(script))
	presetTangoFn(s)
	_ = s.Add("this", nil)
	_ = s.Add("output", nil)
	_ = s.Add("self", nil)
	_ = s.Add("args", nil)
	c, err := s.Compile()
	if err != nil {
		return operationExecutor{}, err
	}
	return operationExecutor{c}, nil
}

func (b operationExecutor) exec(ctx domainStatus, self tengoV, args []tengoV) (bool, error) {
	fmt.Println(123,self,args)
	ctx.ImportTangoFn(b.c)
	b.c.Set("self", self)
	b.c.Set("args", args)
	b.c.Set("output", nil)

	if err := b.c.Run(); err != nil {
		return false, err
	}
	output := b.c.Get("output")
	if bo, ok := output.Value().(bool); ok{
		return bo, nil
	}
	return false, errors.New("operationExecutor output not bool")
}
