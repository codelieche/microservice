package forms

type TicketValidateResponse struct {
	Errcode   int    `json:"errcode"`              // 错误代码，无错误是0
	Errmsg    string `json:"errmsg,omitempty"`     // 错误消息
	ReturnUrl string `json:"return_url,omitempty"` // 跳转的URL
	Name      string `json:"name"`                 // Ticket Name
	Session   string `json:"session,omitempty"`    // Session ID
}
