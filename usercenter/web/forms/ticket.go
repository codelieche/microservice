package forms

type TicketValidateUser struct {
	ID       uint   `json:"id"`       // 用户的ID
	Username string `json:"username"` // 用户名
	Email    string `json:"email"`    // 用户邮箱
	Mobile   string `json:"mobile"`   // 用户手机号
}

type TicketValidateResponse struct {
	Errcode   int                 `json:"errcode"`              // 错误代码，无错误是0
	Errmsg    string              `json:"errmsg,omitempty"`     // 错误消息
	ReturnUrl string              `json:"return_url,omitempty"` // 跳转的URL
	Name      string              `json:"name"`                 // Ticket Name
	Session   string              `json:"session,omitempty"`    // Session ID
	User      *TicketValidateUser `json:"user"`                 // 用户基本信息
}
