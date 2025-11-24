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
	Id       string `json:"id"`
	UserId   string `json:"userId"`
	ToUserId string `json:"toUserId"`
	Content  string `json:"content"`
	Status   int    `json:"status"`
	Time     int64  `json:"time"`
}

type HistoryRequest struct {
	TargetUserId string `json:"targetUserId"`
	PageNum      int64  `json:"pageNum"`  // 从 1 开始
	PageSize     int64  `json:"pageSize"` // 默认 20
}

type UnreadHistoryRequest struct {
	TargetUserId string `json:"targetUserId"`
}

type GroupMessage struct {
	Id      string `json:"id"`
	GroupId string `json:"groupId"`
	UserId  string `json:"userId"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
}

type Group struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type GroupHistoryRequest struct {
	GroupId  string `json:"groupId"`
	PageNum  int64  `json:"pageNum"`  // 从 1 开始
	PageSize int64  `json:"pageSize"` // 默认 20
}
