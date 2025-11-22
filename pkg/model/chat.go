package model

const (
	TypePrivateMessage int = iota // 私聊发送消息
	TypePrivateHistory            // 获取私聊历史记录
	TypePrivateUnread             // 获取私聊未读消息
	TypeGroupMessage              // 群聊发送消息
	TypeGroupHistory              // 获取群聊历史记录
)

type WSMessage struct {
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}

type PrivateMsg struct {
	Id       string `json:"Id"`
	UserId   string `json:"UserId"`
	ToUserID string `json:"toUserId"`
	Content  string `json:"content"`
	Status   int    `json:"status"`
	Time     int64  `json:"time"`
}
