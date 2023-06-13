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
		parkBus := subscribe_bus.NewParkBus(passenger.UserName, passenger.Password, passenger.MorningBusTime, passenger.AfternoonBusTime)
		if parkBus.Login() {
			if passenger.MorningBusTime != "" {
				// 每天上午6点整，开始预约早班
				if _, err := cronTab.AddFunc("0 6 * * *", parkBus.MorningBusSubscribe); err != nil {
					fmt.Printf("%s（%s）定时任务注册失败\n", passenger.UserName, passenger.MorningBusTime)
				}
				fmt.Printf("%s（%s）定时任务注册成功\n", passenger.UserName, passenger.MorningBusTime)
			}

			if passenger.AfternoonBusTime != "" {
				// 每天16点整，开始预约晚班
				if _, err := cronTab.AddFunc("0 16 * * *", parkBus.AfternoonBusSubscribe); err != nil {
					fmt.Printf("%s（%s）定时任务注册失败\n", passenger.UserName, passenger.AfternoonBusTime)
				}
				fmt.Printf("%s（%s）定时任务注册成功\n", passenger.UserName, passenger.AfternoonBusTime)
			}
		}

	}
	cronTab.Start()
	fmt.Println("定时任务启动中...")
}
