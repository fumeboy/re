package engine

import (
	"github.com/d5/tengo/v2"
	"github.com/pkg/errors"
)

/*
factorValue := factor(domain_id, "bridgeCode name")
*/
func factorGetInTengo(runtimeCtx domainStatus) func(args ...tengo.Object) (ret tengo.Object, err error) {
	return func(args ...tengo.Object) (ret tengo.Object, err error) {
		if len(args) != 2{
			return nil, errors.New("factor() 参数数量不为 2")
		}
		var ok bool
		var ds domainStatus
		factor_code, ok := tengo.ToString(args[1])
		if !ok {
			return nil, errors.New("tengo fn factor() 参数 factor_code 非 string 类型")
		}

		if tengo.ToInterface(args[0]) == nil{
			ds = runtimeCtx
		}else{
			path, ok := tengo.ToInterface(args[0]).([]interface{})
			if !ok {
				return nil, errors.New("tengo fn factor() 参数 path 非 Array 类型")
			}
			path2, err := utilArrTo(path)
			if err != nil {
				return nil, err
			}
			ds, err = runtimeCtx.getDomain(path2)
			if err != nil {
				return nil, err
			}
		}
		f,err := ds.getFactor(factor_code)
		if err != nil {
			return nil, err
		}
		return tengo.FromInterface(f.value)
	}
}

func utilArrTo(arr []interface{}) ([]string, error){
	if len(arr) == 0{
		return nil, nil
	}
	var s []string
	for _, a := range arr{
		if s2, ok := a.(string);ok {
			s = append(s, s2)
		}else{
			return nil, errors.New("utilArrTo failed")
		}
	}
	return s, nil
}
