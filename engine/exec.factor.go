package engine

import "github.com/d5/tengo/v2"

/*
构造函数脚本

输入
	this （this domain ID）
输出
	因子值

example
	output = [200,100，150]
	output = true
	output = 123
	output = {
		abc: factor(this, "abc")
	}
*/

type factorConstructorExecutor struct {
	c *tengo.Compiled
}

func newFactorConstructorExecutor(script string) (factorConstructorExecutor, error) {
	s := tengo.NewScript([]byte(script))
	presetTangoFn(s)
	_ = s.Add("input", nil)
	_ = s.Add("output", nil)
	c, err := s.Compile()
	if err != nil {
		return factorConstructorExecutor{}, err
	}
	return factorConstructorExecutor{c}, nil
}

func (b factorConstructorExecutor) exec(ctx runtimeContext, domain_id int) (tengo.Object, error) {
	ctx.ImportTangoFn(b.c)
	b.c.Set("this", domain_id)
	b.c.Set("output", nil)
	if err := b.c.Run(); err != nil {
		panic(err)
	}
	output := b.c.Get("output")
	return output.Object(), nil
}