package engine

import (
	"fmt"
	"re/dal"

	"github.com/pkg/errors"
)

/*
	根据 combination 遍历 rules

	每个 rule
		根据 path 找到 domain，如果未存在，执行构造函数脚本
		得到 domain 后，获取 factor
		如果缓存中未存在，执行factor的构造函数
		获取 operation，将 factor 和 rule 的 args 传入
*/

type domainStatus struct {
	id          int
	bridgeCode  string
	factorCache map[string]interface{}
	domains map[string]domainStatus
}

func (ctx domainStatus) AddDomain(ds domainStatus) {
	ctx.domains[ds.bridgeCode] = ds
}

func (ctx domainStatus) getDomain1(code string) (domainStatus, error) {
	if d, ok := ctx.domains[code]; ok {
		return d, nil
	}
	return domainStatus{}, errors.New("getDomain1 failed")
}

func (ctx domainStatus) getDomain(path []string) (domainStatus, error) {
	var d = ctx
	var err error
	var i int
	for ;i<len(path);i++{
		if d2, err := d.getDomain1(path[i]); err == nil{
			d = d2
		}else{
			break
		}
	}
	for ; i < len(path); i++ {
		d, err = d.getDomain2(path[i])
		if err != nil{
			return domainStatus{}, err
		}
	}
	return d, nil
}

func (ds domainStatus) getDomain2(bridge_code string) (domainStatus, error) {
	// 获取 to domain, 和 bridge
	br, err := dal.DomainBridgeGet(ds.id, bridge_code)
	if err != nil {
		return domainStatus{}, err
	}
	// 执行 bridge scirpt
	exec, err := newBridgeScriptExecutor(br.Constructor)
	if err != nil {
		return domainStatus{}, err
	}
	// 执行 to domain constructor
	to_init, err := exec.exec(ds)
	if err != nil {
		return domainStatus{}, err
	}
	ds2 := domainStatus{
		bridgeCode:  bridge_code,
		id:          br.To,
		factorCache: to_init,
	}
	ds.AddDomain(ds2)
	return ds2, nil
}

func (ds domainStatus) getFactor(code string) (factor, error) {
	// 从缓存获取 factor value
	if v, ok := ds.factorCache[code]; ok {
		return factor{ v, ds}, nil
	}
	// 获取 db factor
	f, err := dal.FactorGet(ds.id, code)
	if err != nil {
		return factor{}, err
	}
	// 执行构造函数
	e, err := newFactorConstructorExecutor(f.Constructor)
	if err != nil {
		return factor{}, err
	}
	v, err := e.exec(ds)
	if err != nil {
		return factor{}, err
	}
	ds.factorCache[code] = v
	return factor{v, ds}, nil
}

type factor struct {
	value    tengoV
	domainStatus
}

func (f factor) execOperation(operationID int, args []Arg) (bool, error) {
	// 获取 operation
	var script string
	{
		op, err := dal.OperationGet(operationID)
		if err != nil {
			return false, err
		}
		script = op.Script
	}
	// TODO, 校验 args 数量
	op, err := newOperationExecutor(script)
	if err != nil {
		return false, err
	}
	return op.exec(f.domainStatus, f.value, ArgsTo(args))
}

func ruleExec(ctx domainStatus, r Rule) (bool, error) {
	domain, err := ctx.getDomain(r.Path)
	if err != nil {
		return false, err
	}
	f, err := domain.getFactor(r.Factor)
	if err != nil {
		return false, err
	}
	return f.execOperation(r.Operation, r.Args)
}

func strategyExec(s Strategy, eventData map[string]interface{}) error {
	// 获取 event domain
	rel, err := dal.EventDomainGet(s.EventID)
	if err != nil {
		return err
	}
	// 初始化 event domain
	exec, err := newEventDomainConstructExecutor(rel.Constructor)
	if err != nil {
		return err
	}
	ds := domainStatus{
		id: rel.DomainID,
		bridgeCode: "",
		domains: map[string]domainStatus{},
	}
	{
		o, err := exec.exec(eventData)
		if err != nil {
			return err
		}
		ds.factorCache = o
	}
	// 执行 ruleExec
	for _, r := range s.Rules {
		result, err := ruleExec(ds, r)
		fmt.Println(r, result, err)
	}
	return nil
}
