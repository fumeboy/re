package engine

import (
	"github.com/d5/tengo/v2"
	"github.com/pkg/errors"
)

type tengoMAP = map[string]tengoV
type tengoV = interface {}

/*
domain2 := domain("bridge_code")

domain 构造domain path

domain path 作为 tengo factor fn 的第一个入参
*/
func domainGetInTengo(args ...tengo.Object) (ret tengo.Object, err error) {
	var path []interface{}
	for _,a := range args{
		if s, ok := tengo.ToString(a); !ok{
			return nil, errors.New("arg not string")
		}else{
			path = append(path, s)
		}
	}
	return tengo.FromInterface(path)
}

func (ctx domainStatus) ImportTangoFn(c *tengo.Compiled){
	c.Set("factor", factorGetInTengo(ctx))
	c.Set("domain", domainGetInTengo)
}

func presetTangoFn(s *tengo.Script){
	s.Add("factor", nil)
	s.Add("domain", nil)
}
