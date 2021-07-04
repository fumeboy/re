// +build test

package db

import (
	"re/db/client"
	"re/db/form"

	"gorm.io/gorm"
)

var db = client.DB

/*
	初始化业务数据用于演示

标记说明：
	<*> 表示对外（对运营）暴露，否则是中间变量
	(1) 圆括号后面接因子的操作符函数

三个域：

	下订单（事件初始域）
		顾客ID（用户）
		商家ID（用户）
		<*>消费金额
			(1) 和上一次消费金额相同
		购买物品ID（商品）

	用户
		ID
		<*>信用分
			(2) 高于百分之 {} 的用户
		<*>本月消费金额列表
			(3) 平均值小于 {}

	商品
		ID
		<*>商品月销量
			(4) 大于等于 {} 且小于 {}

*/

func init() {
	// event
	db.Create(&form.Event{gorm.Model{ID: 1}, "下订单"})
	// event domain
	db.Create(&form.Domain{gorm.Model{ID: 1}, "下订单"})
	{
		{ //factor
			db.Create(&form.Factor{
				Domain:      1,
				Code:        "consumer_id",
				Model:       gorm.Model{ID: 1},
				Name:        "顾客ID",
				Constructor: "",
			})
			db.Create(&form.Factor{
				Domain:      1,
				Code:        "seller_id",
				Model:       gorm.Model{ID: 2},
				Name:        "商家ID",
				Constructor: "",
			})
			db.Create(&form.Factor{
				Domain:      1,
				Code:        "amount",
				Model:       gorm.Model{ID: 3},
				Name:        "消费金额",
				Constructor: "",
			})
			{ // operation
				db.Create(&form.Operation{
					Model:  gorm.Model{ID: 1},
					View:   `和上一次消费金额相同`,
					Script: `
						subDomain := domain(this, "user domain1")
						list := factor(subDomain, "ConsumptionAmountListMonth")
						if len(list) > 0 {
							output = list[0] == self
						}else{
							output = false
						}
					`,
				})
				db.Create(&form.RelFactorOperation{OperationID: 1, FactorID: 3})
			}
			db.Create(&form.Factor{
				Domain:      1,
				Code:        "commodity_id",
				Model:       gorm.Model{ID: 4},
				Name:        "购买物品ID",
				Constructor: "",
			})
		}
		// rel event domain
		db.Create(&form.RelEventDomain{
			EventID:  1,
			DomainID: 1,
			Constructor: `
			output = {
				consumer_id: input.consumer_id,
				seller_id: input.seller_id,
				amount: input.amount,
				commodity_id: input.commodity_id
			}`,
		})
		{ //bridge
			db.Create(&form.DomainBridge{1, 2, "顾客指向的用户域", "user domain1", `
				output = {id: factor(this, "consumer_id")}
			`})
			db.Create(&form.DomainBridge{1, 2, "商家指向的用户域", "user domain2", `
				output = {id: factor(this, "seller_id")}
			`})
			db.Create(&form.DomainBridge{1, 3, "购买物品指向的商品域", "commodity domain", `
				output = {id: factor(this, "commodity_id"), monthly_sales: 188}
			`})
		}
	}

	db.Create(&form.Domain{gorm.Model{ID: 2}, "用户"})
	{ //factor
		db.Create(&form.Factor{
			Domain:      2,
			Code:        "id",
			Model:       gorm.Model{ID: 5},
			Name:        "ID",
			Constructor: "",
		})
		db.Create(&form.Factor{
			Domain:      2,
			Code:        "CreditScore",
			Model:       gorm.Model{ID: 6},
			Name:        "信用分",
			Constructor: "output = 100", // TODO, 这里应该是 rpc 调用结果
		})
		{ // operation
			db.Create(&form.Operation{
				Model:  gorm.Model{ID: 2},
				View:   `高于百分之 {type: "int32", placeholder: "占比"} 的用户`,
				Script: "output = true", // TODO, 这里应该是 rpc 调用结果
			})
			db.Create(&form.RelFactorOperation{OperationID: 2, FactorID: 6})
		}
		db.Create(&form.Factor{
			Domain: 2,
			Code:   "ConsumptionAmountListMonth",
			Model:  gorm.Model{ID: 7},
			Name:   "本月消费金额列表",
			Constructor: `
				// result := rpc.call(PSM, {params})
				// output = ...
				output = [200,100,150]`, // TODO, 这里应该是 rpc 调用结果
		})
		{ // operation
			db.Create(&form.Operation{
				Model: gorm.Model{ID: 3},
				View:  `平均值小于 {type: "int", placeholder: "数量"}`,
				Script: `
					n := 0
					for v in self{
						n += v
					}
					n /= len(self)
					output = n < args[0]
				`,
			})
			db.Create(&form.RelFactorOperation{OperationID: 3, FactorID: 7})
		}
	}

	db.Create(&form.Domain{gorm.Model{ID: 3}, "商品"})
	{ //factor
		db.Create(&form.Factor{
			Domain:      3,
			Code:        "id",
			Model:       gorm.Model{ID: 8},
			Name:        "ID",
			Constructor: "",
		})
		db.Create(&form.Factor{
			Domain:      3,
			Code:        "monthly_sales",
			Model:       gorm.Model{ID: 9},
			Name:        "商品月销量",
			Constructor: "",
		})
		{ // operation
			db.Create(&form.Operation{
				Model:  gorm.Model{ID: 4},
				View:   `大于等于 {type: "int32", placeholder: "数字1"} 且小于 {type: "int32", placeholder: "数字2"}`,
				Script: "output = (self >= args[0] && self < args[1])",
			})
			db.Create(&form.RelFactorOperation{OperationID: 4, FactorID: 9})
		}
	}
}
