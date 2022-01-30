package main

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
	"log"
	"math/rand"
	"time"
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
			TokenCalculateStrategy: flow.WarmUp, // 冷启动策略
			ControlBehavior:        flow.Reject, //超过直接拒绝
			Threshold:              1000,
			WarmUpPeriodSec:        30, //预热时长 1秒可以到1000并发，但是30秒才达到1000的请求
		},
	})
	if err != nil {
		log.Fatalf("配置限流规则失败 :%v", err)
	}
	var globalTotal int
	var passTotal int
	var blockTotal int
	stopChan := make(chan struct{})
	// 每一秒统计一次，这一秒只能通过 你通过了多少 总共多少 block 了多少
	for i := 0; i < 100; i++ {
		go func() {
			for {
				globalTotal++
				e, b := sentinel.Entry("vierhang-test", sentinel.WithTrafficType(base.Inbound))
				if b != nil {
					//违反规则了
					// fmt.Println("限流了")
					blockTotal++
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
				} else {
					passTotal++
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
					e.Exit()
				}
			}
		}()
	}
	go func() {
		var oldTotal int // 过去1s 总共多少个
		var oldPass int  // 过去1s pass
		var oldBlock int // 过去1s block多少个
		for {
			onSecondTotal := globalTotal - oldTotal
			oldTotal = globalTotal

			onSecondPass := passTotal - oldPass
			oldPass = passTotal

			oneSecondBlock := blockTotal - oldBlock
			oldBlock = blockTotal
			time.Sleep(time.Second)
			fmt.Printf("total :%d,pass:%d,block:%d \n", onSecondTotal, onSecondPass, oneSecondBlock)
		}
	}()
	<-stopChan
}
