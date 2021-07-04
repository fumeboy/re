package engine

import (
	"fmt"
	"github.com/d5/tengo/v2"
	"github.com/pkg/errors"
)

type tengoMAP = map[string]tengoV
type tengoV = interface {}

/*
domain2 := bridge(this, "bridge_code")
*/
func domainGetInTengo(runtimeCtx runtimeContext) func(args ...tengo.Object) (ret tengo.Object, err error) {
	return func(args ...tengo.Object) (ret tengo.Object, err error) {
		if len(args) != 2 {
			return nil, errors.New("参数数量不为 2")
		}
		var ok bool
		domain_id, ok := tengo.ToInt(args[0])
		if !ok {
			return nil, errors.New("获取 domain id 失败")
		}
		bridge_code, ok := tengo.ToString(args[1])
		if !ok {
			return nil, errors.New("获取 bridge code 失败")
		}
		ds, err := runtimeCtx.getDomainDirect(domain_id)
		if err != nil {
			fmt.Println(runtimeCtx, domain_id)
			return nil, err
		}
		ds2,err := ds.getDomain2(runtimeCtx, bridge_code)
		if err != nil {
			return nil, err
		}
		return tengo.FromInterface(ds2.id)
	}
}

func (ctx runtimeContext) ImportTangoFn(c *tengo.Compiled){
	c.Set("factor", factorGetInTengo(ctx))
	c.Set("domain", domainGetInTengo(ctx))
}

func presetTangoFn(s *tengo.Script){
	s.Add("factor", nil)
	s.Add("domain", nil)
}
