package subscribe_bus

import (
	"fmt"
	"log"
	"time"

	"park-bus/pkg"

	"github.com/imroc/req/v3"
)

var (
	loginAPI           = "http://gl.yichengshidai.com/api/auth/login"
	getMorningBusAPI   = "http://gl.yichengshidai.com/api/api-oa/park-bus-app/appointmentList/1/1"
	getAfternoonBusAPI = "http://gl.yichengshidai.com/api/api-oa/park-bus-app/appointmentList/2/1"
	selectBusAPI       = "http://gl.yichengshidai.com/api/api-oa/park-bus-app/appointmentBus/%d"
)

type ParkBus struct {
	HttpClient       *req.Client
	UserName         string
	Password         string
	Token            string
	morningBusTime   string // 早班车时间 (eg: 08:30)
	afternoonBusTime string // 晚班车时间 (eg: 17:50)
}

// NewParkBus 初始化班车预约
func NewParkBus(token, morningBusTime, afternoonBusTime string) ParkBus {
	httpClient := req.NewClient()
	httpClient.SetTimeout(5 * time.Second)
	httpClient.SetCommonHeader("User-Agent", "User-Agent\tMozilla/5.0 (iPhone; CPU iPhone OS 15_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Html5Plus/1.0 (Immersed/20) uni-app")
	httpClient.SetCommonHeader("Content-Type", "application/json;charset=UTF-8")
	httpClient.SetCommonBearerAuthToken(token)
	httpClient.SetCommonRetryCount(10).
		SetCommonRetryBackoffInterval(1*time.Second, 5*time.Second).
		SetCommonRetryFixedInterval(2 * time.Second)

	return ParkBus{
		HttpClient:       httpClient,
		Token:            token,
		morningBusTime:   morningBusTime,
		afternoonBusTime: afternoonBusTime}
}

// Login abandon
func (pb *ParkBus) Login() bool {
	loginResp := pkg.LoginRest{}
	resp, err := pb.HttpClient.R().
		SetBody(fmt.Sprintf(`{"username":"%s","password":"%s"}`, pb.UserName, pb.Password)).
		SetSuccessResult(&loginResp).
		Post(loginAPI)

	if err != nil {
		fmt.Printf("登录请求错误: %s\n", err.Error())
		return false
	}

	if !resp.IsSuccessState() {
		fmt.Printf("登录错误: %s\n", resp.String())
		return false
	}

	if loginResp.Code != "0" {
		fmt.Println("登录失败")
		fmt.Println(resp.String())
		return false
	}
	// fmt.Printf("登录成功: %s\n", resp.String())
	pb.HttpClient.SetCommonBearerAuthToken(loginResp.Data.Token)
	return true
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
		fmt.Println("获取车次列表失败:", err.Error())
		return true
	}

	if busList.Code != "0" {
		fmt.Printf("获取车次列表失败: %s\n", resp.String())
		// 重新登录
		if busList.Code == "0004" {
			// fmt.Println("重新登录...")
			fmt.Println("认证信息失败。")
			return false
		}
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
		// SetBearerAuthToken(pb.Token).
		AddRetryCondition(pb.noTicketRetry). // 班车余票获取
		Get(getMorningBusAPI)

	if err != nil {
		log.Printf("获取早班(%s)车车次信息错误: %s\n", pb.morningBusTime, err.Error())
		return
	}

	if resp.IsSuccessState() {
		fmt.Printf("车次获取重试: %d次\n", resp.Request.RetryAttempt)
		for _, busInfo := range busList.Data {
			// 筛选出有效班车
			if busInfo.ShiftTime == pb.morningBusTime && busInfo.TicketTotal > 0 {
				// 锁定车票
				url := fmt.Sprintf(selectBusAPI, busInfo.Id)
				selectBusResp, err := pb.HttpClient.R().Get(url)
				if err != nil {
					log.Printf("早班车(%s)车票锁定错误: %s\n", pb.morningBusTime, err.Error())
					return
				}
				if !selectBusResp.IsSuccessState() {
					fmt.Printf("早班车(%s)车次锁定失败 : %s\n", pb.morningBusTime, resp.String())
				}
				fmt.Printf("早班车(%s)预约成功。\n", pb.morningBusTime)
				fmt.Println(selectBusResp.String())
				return
			}
		}
	}

	fmt.Printf("当前时间：%s, 早班车(%s)，抢票失败。\n", time.Now().String(), pb.morningBusTime)
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
		fmt.Printf("获取晚班车(%s)车次信息错误: %s\n", pb.afternoonBusTime, err.Error())
		return
	}

	if resp.IsSuccessState() {
		fmt.Printf("车次获取重试: %d次\n", resp.Request.RetryAttempt)
		for _, busInfo := range busList.Data {
			// 筛选出有效班车
			if busInfo.ShiftTime == pb.afternoonBusTime && busInfo.TicketTotal > 0 {
				// 锁定车票
				url := fmt.Sprintf(selectBusAPI, busInfo.Id)
				selectBusResp, err := pb.HttpClient.R().Get(url)
				if err != nil {
					fmt.Printf("晚班车(%s)车票锁定错误: %s \n", pb.afternoonBusTime, err.Error())
					return
				}
				if !selectBusResp.IsSuccessState() {
					fmt.Printf("晚班车(%s)车次锁定失败: %s\n", pb.afternoonBusTime, resp.String())
				}
				fmt.Printf("晚班车(%s)抢票成功。\n", pb.afternoonBusTime)
				fmt.Println(selectBusResp.String())
				return
			}
		}
	}

	fmt.Printf("当前时间：%s, 晚班车(%s)，抢票失败。\n", time.Now().String(), pb.afternoonBusTime)
	return
}
