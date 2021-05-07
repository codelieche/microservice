package core

type Response struct {
	Code    int         `json:"code"`    // 返回的code，如果是0就表示正常
	Data    interface{} `json:"data"`    // 返回的数据
	Message string      `json:"message"` // 返回的消息
}

type ResponseList struct {
	CurrentPage int         `json:"current_page"`
	Count       int64       `json:"count"`   // 当前列表的数据
	Results     interface{} `json:"results"` // 返回的列表数据
}
