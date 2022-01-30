package main

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
	"log"
)

func main() {
	// 初始化sentinel
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("初始化失败 :%v", err)
	}

	//配置限流规则
	_, err = flow.LoadRules([]*flow.Rule{
		// 一秒钟 10 个流量
		{
			Resource:               "vierhang-test",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject, //超过直接拒绝
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "vierhang-test2",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject, //超过直接拒绝
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		log.Fatalf("配置限流规则失败 :%v", err)
	}
	// 调用12次
	for i := 0; i < 12; i++ {
		//Inbound 流控入口点，规则是vierhang-test
		//Outbound 出口流量控制
		e, b := sentinel.Entry("vierhang-test", sentinel.WithTrafficType(base.Inbound))
		if b != nil {
			//违反规则了
			fmt.Println("限流了")
		} else {
			fmt.Println("检查通过")
			e.Exit()
		}

	}
}
