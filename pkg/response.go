package pkg

type Result struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type BusList struct {
	Result
	Data []BusInfo `json:"data"`
}

type BusInfo struct {
	Id            int    `json:"id"`
	TicketTotal   int    `json:"ticketTotal"`
	ShiftDate     string `json:"shiftDate"`
	ShiftTime     string `json:"shiftTime"`
	BusAddress    string `json:"busAddress"`
	Announcements string `json:"announcements"`
}

type LoginRest struct {
	Result
	Data LoginData `json:"data"`
}

type LoginData struct {
	RefreshToken string `json:"refresh_token"`
	YdPwd        string `json:"yd_pwd"`
	OrgType      string `json:"org_type"`
	YdName       string `json:"yd_name"`
	Token        string `json:"token"`
	Username     string `json:"username"`
}
