package cron_job

import (
	"fmt"
	"github.com/robfig/cron/v3"
	conf "park-bus/config"
	"park-bus/pkg/subscribe_bus"
)

func Run() {
	// 定时任务
	fmt.Println("注册定时任务...")
	cronTab := cron.New()
	for _, passenger := range conf.Cfg.Passengers {
		parkBus := subscribe_bus.NewParkBus(passenger.Token, passenger.MorningBusTime, passenger.AfternoonBusTime)
		if passenger.MorningBusTime != "" {
			cronTab.AddFunc("0 6 * * *", parkBus.MorningBusSubscribe) // 每天上午6点整，开始预约早班
		}

		if passenger.AfternoonBusTime != "" {
			cronTab.AddFunc("0 16 * * *", parkBus.AfternoonBusSubscribe) // 每天16点整，开始预约晚班
		}
	}
	cronTab.Start()
	fmt.Println("定时任务启动中...")
}
