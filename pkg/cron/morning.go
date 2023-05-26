package cron

import (
	"fmt"
	"github.com/imroc/req/v3"
	"log"
	"park-bus/pkg"
	"time"
)

var (
	morningBusId     = 2
	getMorningBusAPI = "http://gl.yichengshidai.com/api/api-oa/park-bus-app/appointmentList/1"
	selectBusAPI     = "http://gl.yichengshidai.com/api/api-oa/park-bus-app/appointmentBus/2"
)

type MorningBus struct {
	HttpClient *req.Client
}

func NewMorningBus(Authorization string) *MorningBus {
	httpClient := req.NewClient()
	httpClient.SetTimeout(5 * time.Second)
	httpClient.SetCommonBearerAuthToken(Authorization)
	httpClient.SetCommonHeader("User-Agent", "User-Agent\tMozilla/5.0 (iPhone; CPU iPhone OS 15_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Html5Plus/1.0 (Immersed/20) uni-app")
	httpClient.SetCommonHeader("Content-Type", "application/json;charset=UTF-8")

	httpClient.SetCommonRetryCount(3).
		SetCommonRetryBackoffInterval(1*time.Second, 5*time.Second).
		SetCommonRetryFixedInterval(2 * time.Second)

	return &MorningBus{HttpClient: httpClient}
}

func noTicketRetry(resp *req.Response, err error) bool {
	if err != nil {
		return true
	}

	if !resp.IsSuccessState() {
		return true
	}

	return false
}

func (mb *MorningBus) GetMorningBus() {
	busList := pkg.BusList{}

	resp, err := mb.HttpClient.R().
		SetSuccessResult(&busList).
		AddRetryCondition(noTicketRetry). // 班车余票获取
		Get(getMorningBusAPI)

	if err != nil {
		log.Println("error:", err)
		return
	}

	if resp.IsSuccessState() {
		fmt.Printf("response : %v\n", busList)
		for _, busInfo := range busList.Data {
			if busInfo.Id == morningBusId && busInfo.TicketTotal > 0 {

			}
		}
	}
	return
}

func (mb *MorningBus) SelectBus() {

}
