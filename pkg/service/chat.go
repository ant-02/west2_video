package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
	"west2/database"
	"west2/pkg/model"
	"west2/util"
)

const privateMessageScript = `
	redis.call("ZADD", KEYS[1], ARGV[1], KEYS[2])
	redis.call("HSET", KEYS[2], unpack(ARGV, 2))
	return 1
`

type chatService struct{}

type ChatService interface {
	Chat(uid string, reqMsg []byte) ([]byte, error)
	handlePrivateMessage(uid string, msgType int, msg interface{}) ([]byte, error)
}

func NewChatService() ChatService {
	return &chatService{}
}

func (cs *chatService) Chat(uid string, reqMsg []byte) ([]byte, error) {
	var msg model.WSMessage
	if err := json.Unmarshal(reqMsg, &msg); err != nil {
		return nil, err
	}

	switch msg.Type {
	case model.TypePrivateMessage:
		return cs.handlePrivateMessage(uid, msg.Type, msg.Data)
	// case model.TypePrivateHistory:
	// 	return handlePrivateHistory(uid, msg)
	// case model.TypePrivateUnread:
	// 	return handlePrivateUnread(uid, msg)
	// case model.TypeGroupMessage:
	// 	return handleGroupMessage(uid, msg)
	// case model.TypeGroupHistory:
	// 	return handleGroupHistory(uid, msg)
	default:
		return sendError(msg.Type, "unknown message type", nil)
	}
}

func (cs *chatService) handlePrivateMessage(uid string, msgType int, msg interface{}) ([]byte, error) {
	dataBytes, err := json.Marshal(msg)
	if err != nil {
		return sendError(msgType, "failed to marshal message data", err)
	}

	var privateMsg model.PrivateMsg
	if err := json.Unmarshal(dataBytes, &privateMsg); err != nil {
		return sendError(msgType, "failed to unmarshal private message", err)
	}

	privateMsg.Id = util.GetID()
	privateMsg.Status = 0
	privateMsg.Time = time.Now().Unix()
	privateMsg.UserId = uid
	args := []interface{}{
		strconv.FormatFloat(float64(privateMsg.Time), 'f', -1, 64), // ARGV[1]: score
		"Id", privateMsg.Id,
		"Content", privateMsg.Content,
		"Time", strconv.FormatInt(privateMsg.Time, 10),
		"Status", strconv.Itoa(privateMsg.Status),
		"ToUserID", privateMsg.ToUserID,
	}
	instance := database.GetRedisInstance()
	ctx := context.Background()
	err = instance.Eval(ctx, privateMessageScript, []string{"message:" + privateMsg.UserId + ":" + privateMsg.UserId, privateMsg.Id}, args)
	if err != nil {
		return sendError(msgType, "failed to add private message to redis", err)
	}

	resMsg, _ := json.Marshal(&model.WSMessage{
		Type: msgType,
		Data: "success",
	})
	return resMsg, nil
}

// func handlePrivateHistory() {}

// func handlePrivateUnread() {}

// func handleGroupMessage() {}

// func handleGroupHistory() {}

func sendError(msgType int, msg string, err error) ([]byte, error) {
	log.Printf(msg+": msgType: %d, err: %v", msgType, err)
	resMsg, _ := json.Marshal(&model.WSMessage{
		Type: msgType,
		Data: "internal server error",
	})
	return resMsg, fmt.Errorf(msg, "")
}
