package subscribe_bus

import (
	"fmt"
	"github.com/imroc/req/v3"
	"log"
	"park-bus/pkg"
	"time"
)

var (
	getMorningBusAPI   = "http://gl.yichengshidai.com/api/api-oa/park-bus-app/appointmentList/1"
	getAfternoonBusAPI = "http://gl.yichengshidai.com/api/api-oa/park-bus-app/appointmentList/2"
	selectBusAPI       = "http://gl.yichengshidai.com/api/api-oa/park-bus-app/appointmentBus/%d"
)

type ParkBus struct {
	HttpClient       *req.Client
	morningBusTime   string // 早班车时间 (eg: 08:30)
	afternoonBusTime string // 晚班车时间 (eg: 17:50)
}

// NewParkBus 初始化班车预约
func NewParkBus(token, morningBusTime, afternoonBusTime string) ParkBus {
	httpClient := req.NewClient()
	httpClient.SetTimeout(5 * time.Second)
	httpClient.SetCommonBearerAuthToken(token)
	httpClient.SetCommonHeader("User-Agent", "User-Agent\tMozilla/5.0 (iPhone; CPU iPhone OS 15_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Html5Plus/1.0 (Immersed/20) uni-app")
	httpClient.SetCommonHeader("Content-Type", "application/json;charset=UTF-8")

	httpClient.SetCommonRetryCount(20).
		SetCommonRetryBackoffInterval(1*time.Second, 5*time.Second).
		SetCommonRetryFixedInterval(2 * time.Second)

	return ParkBus{HttpClient: httpClient, morningBusTime: morningBusTime, afternoonBusTime: afternoonBusTime}
}

// noTicketRetry 班车余票获取重试
func (pb *ParkBus) noTicketRetry(resp *req.Response, err error) bool {
	if resp.Request.RetryAttempt == 0 {
		fmt.Println("获取车次列表...")
	} else {
		fmt.Println("尝试重新获取车次列表...")
	}

	if err != nil {
		fmt.Println("获取车次列表错误:", err)
		return true
	}

	if !resp.IsSuccessState() {
		fmt.Println("获取车次列表失败:", resp.String())
		return true
	}

	busList := pkg.BusList{}
	err = resp.UnmarshalJson(&busList)
	if err != nil {
		return true
	}

	for _, busInfo := range busList.Data {
		if busInfo.ShiftTime == pb.morningBusTime {
			fmt.Println("早班车余票：", busInfo.TicketTotal)
			// 检查余票
			if busInfo.TicketTotal == 0 {
				return true
			} else {
				return false
			}
		}
		if busInfo.ShiftTime == pb.afternoonBusTime {
			fmt.Println("晚班车余票：", busInfo.TicketTotal)
			// 检查余票
			if busInfo.TicketTotal == 0 {
				return true
			} else {
				return false
			}
		}
	}

	return false
}

// MorningBusSubscribe 预约早班车
func (pb *ParkBus) MorningBusSubscribe() {
	busList := pkg.BusList{}
	// 获取班车列表
	resp, err := pb.HttpClient.R().
		SetSuccessResult(&busList).
		AddRetryCondition(pb.noTicketRetry). // 班车余票获取
		Get(getMorningBusAPI)

	if err != nil {
		log.Println("获取早班车车次信息错误:", err)
		return
	}

	if resp.IsSuccessState() {
		fmt.Printf("车次获取重试: %d次\n", resp.Request.RetryAttempt)
		for _, busInfo := range busList.Data {
			// 筛选出有效班车
			if busInfo.ShiftTime == pb.morningBusTime && busInfo.TicketTotal > 0 {
				url := fmt.Sprintf(selectBusAPI, busInfo.Id)
				selectBusResp, err := pb.HttpClient.R().Get(url)
				if err != nil {
					log.Println("早班车车票锁定错误: ", err)
					return
				}
				if !selectBusResp.IsSuccessState() {
					fmt.Printf("早班车车次锁定失败 : %v\n", resp.String())
				}
				fmt.Printf("早班车(%s)预约成功。\n", pb.morningBusTime)
				fmt.Println(selectBusResp.String())
				return
			}
		}
	}

	fmt.Printf("早班车(%s)时间已经过期，抢票失败。\n", pb.morningBusTime)
	return
}

// AfternoonBusSubscribe 预约晚班车
func (pb *ParkBus) AfternoonBusSubscribe() {
	busList := pkg.BusList{}
	// 获取班车列表
	resp, err := pb.HttpClient.R().
		SetSuccessResult(&busList).
		AddRetryCondition(pb.noTicketRetry). // 班车余票获取
		Get(getAfternoonBusAPI)

	if err != nil {
		fmt.Println("获取晚班车车次信息错误: ", err)
		return
	}

	if resp.IsSuccessState() {
		fmt.Printf("车次获取重试: %d次\n", resp.Request.RetryAttempt)
		for _, busInfo := range busList.Data {
			// 筛选出有效班车
			if busInfo.ShiftTime == pb.afternoonBusTime && busInfo.TicketTotal > 0 {
				url := fmt.Sprintf(selectBusAPI, busInfo.Id)
				selectBusResp, err := pb.HttpClient.R().Get(url)
				if err != nil {
					log.Println("晚班车车票锁定错误: ", err)
					return
				}
				if !selectBusResp.IsSuccessState() {
					fmt.Println("晚班车车次锁定失败。")
				}
				fmt.Printf("晚班车(%s)抢票成功。\n", pb.afternoonBusTime)
				fmt.Println(selectBusResp.String())
				return
			}
		}
	}

	fmt.Printf("晚班车(%s)时间已经过期，抢票失败。\n", pb.afternoonBusTime)
	return
}
