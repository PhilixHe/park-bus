package test

import (
	"park-bus/config"
	"park-bus/pkg/subscribe_bus"
	"testing"
)

func TestGetMorningBus(t *testing.T) {
	config.LoadConfig("../config/config.yaml")
	username := config.Cfg.Passengers[0].UserName
	password := config.Cfg.Passengers[0].Password
	morningTime := config.Cfg.Passengers[0].MorningBusTime
	afternoonTime := config.Cfg.Passengers[0].AfternoonBusTime
	parkBus := subscribe_bus.NewParkBus(username, password, morningTime, afternoonTime)
	parkBus.Login()
	parkBus.MorningBusSubscribe()

}

func TestGetAfternoonBus(t *testing.T) {
	config.LoadConfig("../config/config.yaml")
	username := config.Cfg.Passengers[0].UserName
	password := config.Cfg.Passengers[0].Password
	morningTime := config.Cfg.Passengers[0].MorningBusTime
	afternoonTime := config.Cfg.Passengers[0].AfternoonBusTime
	parkBus := subscribe_bus.NewParkBus(username, password, morningTime, afternoonTime)
	parkBus.Login()
	parkBus.AfternoonBusSubscribe()
}
