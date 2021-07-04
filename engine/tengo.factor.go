package engine

import (
	"github.com/d5/tengo/v2"
	"github.com/pkg/errors"
)

/*
factorValue := factor(domain_id, "code name")
*/
func factorGetInTengo(runtimeCtx runtimeContext) func(args ...tengo.Object) (ret tengo.Object, err error) {
	return func(args ...tengo.Object) (ret tengo.Object, err error) {
		if len(args) != 2 {
			return nil, errors.New("tengo fn factor() 参数数量不为 2")
		}
		var ok bool
		domain_id, ok := tengo.ToInt(args[0])
		if !ok {
			return nil, errors.New("tengo fn factor() 参数 domain_id 非 int 类型")
		}
		factor_code, ok := tengo.ToString(args[1])
		if !ok {
			return nil, errors.New("tengo fn factor() 参数 factor_code 非 string 类型")
		}
		ds, err := runtimeCtx.getDomainDirect(domain_id)
		if err != nil {
			return nil, err
		}
		f,err := ds.getFactor(runtimeCtx, factor_code)
		if err != nil {
			return nil, err
		}
		return tengo.FromInterface(f.value)
	}
}
