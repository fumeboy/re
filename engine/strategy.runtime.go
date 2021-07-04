package engine

import (
	"fmt"
	"github.com/pkg/errors"
	"re/dal"
)

/*
	根据 combination 遍历 rules

	每个 rule
		根据 path 找到 domain，如果未存在，执行构造函数脚本
		得到 domain 后，获取 factor
		如果缓存中未存在，执行factor的构造函数
		获取 operation，将 factor 和 rule 的 args 传入
*/

type runtimeContext struct {
	root    int // domain root
	domains map[int]domainStatus
}

func (ctx runtimeContext) AddDomain(ds domainStatus) {
	ctx.domains[ds.id] = ds
}

func (ctx runtimeContext) getDomainDirect(id int) (domainStatus, error) {
	if d, ok := ctx.domains[id]; ok {
		return d, nil
	}
	return domainStatus{}, errors.New("getDomainDirect failed")
}

func (ctx runtimeContext) getDomain(path []int) (domainStatus, error) {
	var err error
	var target int
	if len(path) == 0 {
		target = ctx.root
	} else {
		target = path[len(path)-1]
	}

	if d, ok := ctx.domains[target]; ok {
		return d, nil
	}

	var have = -1
	var d domainStatus
	for i := len(path) - 1; i >= 0; i-- {
		if d2, ok := ctx.domains[path[i]]; ok {
			have = i
			d = d2
			break
		}
	}
	if have == -1{
		d = ctx.domains[ctx.root]
	}

	for i := have+1; i < len(path); i++ {
		d, err = d.getDomain(ctx, path[i])
		if err != nil{
			return domainStatus{}, err
		}
	}
	return d, nil
}

type domainStatus struct {
	id          int
	factorCache map[string]interface{}
}

func (ds domainStatus) getDomain(ctx runtimeContext, to_domain_id int) (domainStatus, error){
	// 获取 to domain, 和 bridge
	bridge, err := dal.DomainBridgeGet(ds.id,to_domain_id)
	if err != nil{
		return domainStatus{}, err
	}
	// 执行 bridge scirpt
	exec,err := newBridgeScriptExecutor(bridge.Constructor)
	if err != nil{
		return domainStatus{}, err
	}
	// 执行 to domain constructor
	to_init, err := exec.exec(ctx, ds.id)
	if err != nil{
		return domainStatus{}, err
	}
	ds2 := domainStatus{
		id:          to_domain_id,
		factorCache: to_init,
	}
	ctx.AddDomain(ds2)
	return ds2, nil
}

func (ds domainStatus) getDomain2(ctx runtimeContext, bridge_code string) (domainStatus, error){
	// 获取 to domain, 和 bridge
	br, err := dal.DomainBridgeGet2(ds.id, bridge_code)
	if err != nil{
		return domainStatus{}, err
	}
	// 执行 bridge scirpt
	exec,err := newBridgeScriptExecutor(br.Constructor)
	if err != nil{
		return domainStatus{}, err
	}
	// 执行 to domain constructor
	to_init, err := exec.exec(ctx, ds.id)
	if err != nil{
		return domainStatus{}, err
	}
	ds2 := domainStatus{
		id:          br.To,
		factorCache: to_init,
	}
	ctx.AddDomain(ds2)
	return ds2, nil
}

func (ds domainStatus) getFactor(ctx runtimeContext, code string) (factor, error) {
	// 从缓存获取 factor value
	if v, ok := ds.factorCache[code]; ok{
		return factor{ds.id, v}, nil
	}
	// 获取 db factor
	f,err := dal.FactorGet(ds.id, code)
	if err != nil{
		return factor{}, err
	}
	// 执行构造函数
	e,err := newFactorConstructorExecutor(f.Constructor)
	if err != nil{
		return factor{}, err
	}
	v, err := e.exec(ctx, ds.id)
	if err != nil{
		return factor{}, err
	}
	ds.factorCache[code] = v
	return factor{ds.id, v}, nil
}

type factor struct {
	domainID int
	value    tengoV
}

func (f factor) execOperation(ctx runtimeContext, operationID int, args []Arg) (bool, error) {
	// 获取 operation
	var script string
	{
		op,err := dal.OperationGet(operationID)
		if err != nil{
			return false, err
		}
		script = op.Script
	}
	// TODO, 校验 args 数量
	op, err := newOperationExecutor(script)
	if err != nil{
		return false, err
	}
	return op.exec(ctx, f.domainID, f.value, ArgsTo(args))
}

func ruleExec(ctx runtimeContext, r Rule) (bool, error) {
	domain, err := ctx.getDomain(r.Path)
	if err != nil {
		return false, err
	}
	f, err := domain.getFactor(ctx, r.Factor)
	if err != nil {
		return false, err
	}
	return f.execOperation(ctx, r.Operation, r.Args)
}

func strategyExec(s Strategy, eventData map[string]interface{}) error {
	// 获取 event domain
	rel, err := dal.EventDomainGet(s.EventID)
	if err != nil {
		return err
	}
	// 创建 context
	ctx := runtimeContext{
		root:    rel.DomainID,
		domains: map[int]domainStatus{},
	}
	// 初始化 event domain
	exec, err := newEventDomainConstructExecutor(rel.Constructor)
	if err != nil {
		return err
	}
	var ds domainStatus
	{
		o, err := exec.exec(eventData)
		if err != nil {
			return err
		}
		ds.factorCache = o
		ds.id = rel.DomainID
	}
	ctx.AddDomain(ds)
	// 执行 ruleExec
	for _,r := range s.Rules{
		result,err := ruleExec(ctx, r)
		fmt.Println(r, result, err)
	}
	return nil
}
