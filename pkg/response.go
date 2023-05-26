package pkg

type Result struct {
	Code int    `json:"code"`
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
