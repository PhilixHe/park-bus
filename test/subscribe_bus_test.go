package test

import (
	"park-bus/config"
	"park-bus/pkg/subscribe_bus"
	"testing"
)

func TestGetMorningBus(t *testing.T) {
	config.LoadConfig("../config/config.yaml")
	token := config.Cfg.Passengers[0].Token
	morningTime := config.Cfg.Passengers[0].MorningBusTime
	afternoonTime := config.Cfg.Passengers[0].AfternoonBusTime
	parkBus := subscribe_bus.NewParkBus(token, morningTime, afternoonTime)

	parkBus.MorningBusSubscribe()

}

func TestGetAfternoonBus(t *testing.T) {
	config.LoadConfig("../config/config.yaml")
	token := config.Cfg.Passengers[0].Token
	morningTime := config.Cfg.Passengers[0].MorningBusTime
	afternoonTime := config.Cfg.Passengers[0].AfternoonBusTime
	parkBus := subscribe_bus.NewParkBus(token, morningTime, afternoonTime)

	parkBus.AfternoonBusSubscribe()
}
