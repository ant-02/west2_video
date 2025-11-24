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

const setPrivateMessageScript = `
	redis.call("ZADD", KEYS[1], ARGV[1], KEYS[3])
	redis.call("ZADD", KEYS[2], ARGV[1], KEYS[3])
	redis.call("HSET", KEYS[3], unpack(ARGV, 2))
	return 1
`

const getUnreadHistoryMessageScript = `
	local ids = redis.call("ZREVRANGE", KEYS[1], ARGV[1], ARGV[2])
	local messages = {}
	for i, id in ipairs(ids) do 
		local message = redis.call('HGETALL', id)
		messages[i] = message
		redis.call('HSET', id, 'Status', 1)
	end
	redis.call('DEL', KEYS[1])
	return messages
`

const setGroupMessageScript = `
	redis.call("ZADD", KEYS[1], ARGV[1], KEYS[2])
	redis.call("HSET", KEYS[2], unpack(ARGV, 2))
	return 1
`

type chatService struct{}

type ChatService interface {
	Chat(uid string, reqMsg []byte) ([]byte, error)
	handlePrivateMessage(uid string, msg *model.WSMessage) ([]byte, error)
	handlePrivateHistory(uid string, msg *model.WSMessage) ([]byte, error)
	handlePrivateUnread(uid string, msg *model.WSMessage) ([]byte, error)
	handleGroupMessage(uid string, msg *model.WSMessage) ([]byte, error)
	handleGroupHistory(uid string, msg *model.WSMessage) ([]byte, error)
}

func NewChatService() ChatService {
	return &chatService{}
}

func (cs *chatService) Chat(uid string, reqMsg []byte) ([]byte, error) {
	var msg model.WSMessage
	if err := json.Unmarshal(reqMsg, &msg); err != nil {
		return cs.sendError(msg.Type, "failed to unmarshal message", nil)
	}

	switch msg.Type {
	case model.TypePrivateMessage:
		return cs.handlePrivateMessage(uid, &msg)
	case model.TypePrivateHistory:
		return cs.handlePrivateHistory(uid, &msg)
	case model.TypePrivateUnread:
		return cs.handlePrivateUnread(uid, &msg)
	case model.TypeGroupMessage:
		return cs.handleGroupMessage(uid, &msg)
	case model.TypeGroupHistory:
		return cs.handleGroupHistory(uid, &msg)
	default:
		return cs.sendError(msg.Type, "unknown message type", nil)
	}
}

func (cs *chatService) handlePrivateMessage(uid string, msg *model.WSMessage) ([]byte, error) {
	dataBytes, err := json.Marshal(msg.Data)
	if err != nil {
		return cs.sendError(msg.Type, "failed to marshal message data", err)
	}

	var privateMsg model.PrivateMsg
	if err := json.Unmarshal(dataBytes, &privateMsg); err != nil {
		return cs.sendError(msg.Type, "failed to unmarshal private message", err)
	}

	privateMsg.Id = util.GetID()
	privateMsg.Status = 0
	privateMsg.Time = time.Now().Unix()
	privateMsg.UserId = uid
	args := []interface{}{
		strconv.FormatFloat(float64(privateMsg.Time), 'f', -1, 64), // ARGV[1]: score
		"Id", privateMsg.Id,
		"UserId", privateMsg.UserId,
		"ToUserId", privateMsg.ToUserId,
		"Content", privateMsg.Content,
		"Time", strconv.FormatInt(privateMsg.Time, 10),
		"Status", strconv.Itoa(privateMsg.Status),
	}
	instance := database.GetRedisInstance()
	ctx := context.Background()
	_, err = instance.Eval(ctx, setPrivateMessageScript, []string{"message:all:" + privateMsg.UserId + ":" + privateMsg.ToUserId, "message:unread:" + privateMsg.UserId + ":" + privateMsg.ToUserId, privateMsg.Id}, args)
	if err != nil {
		return cs.sendError(msg.Type, "failed to add private message to redis", err)
	}

	resMsg, err := json.Marshal(&model.WSMessage{
		Type: msg.Type,
		Data: "success",
	})
	if err != nil {
		return cs.sendError(msg.Type, "failed to marshal message", err)
	}
	return resMsg, nil
}

func (cs *chatService) handlePrivateHistory(uid string, msg *model.WSMessage) ([]byte, error) {
	dataBytes, err := json.Marshal(msg.Data)
	if err != nil {
		return cs.sendError(msg.Type, "failed to marshal message data", err)
	}

	var historyReq model.HistoryRequest
	if err := json.Unmarshal(dataBytes, &historyReq); err != nil {
		return cs.sendError(msg.Type, "failed to unmarshal history request", err)
	}

	instance := database.GetRedisInstance()
	ctx := context.Background()
	ids, err := instance.ZRevRange(ctx, "message:all:"+uid+":"+historyReq.TargetUserId, (historyReq.PageNum-1)*historyReq.PageSize-1, historyReq.PageNum*historyReq.PageSize)
	if err != nil {
		return cs.sendError(msg.Type, "failed to get history message ids from redis", err)
	}
	privateMsgs := make([]*model.PrivateMsg, len(ids))
	for i, id := range ids {
		hash, err := instance.HGetAll(ctx, id)
		if err != nil {
			return cs.sendError(msg.Type, "failed to get history message from redis by id", err)
		}
		var privateMsg model.PrivateMsg
		privateMsg.Id = id
		userId, ok := hash["UserId"]
		if !ok {
			return cs.sendError(msg.Type, "failed to parse message userId", err)
		}
		privateMsg.UserId = userId
		toUserId, ok := hash["ToUserId"]
		if !ok {
			return cs.sendError(msg.Type, "failed to parse message toUserId", err)
		}
		privateMsg.ToUserId = toUserId
		content, ok := hash["Content"]
		if !ok {
			return cs.sendError(msg.Type, "failed to parse message content", err)
		}
		privateMsg.Content = content
		t1, ok := hash["Time"]
		if !ok {
			return cs.sendError(msg.Type, "failed to parse message time", err)
		}
		t2, err := strconv.ParseInt(t1, 10, 64)
		if err != nil {
			return cs.sendError(msg.Type, "failed to parse string to int64", err)
		}
		privateMsg.Time = t2
		sts, ok := hash["Status"]
		if !ok {
			return cs.sendError(msg.Type, "failed to parse message status", err)
		}
		status, err := strconv.ParseInt(sts, 10, 64)
		privateMsg.Status = int(status)

		privateMsgs[i] = &privateMsg
	}

	resMsg, err := json.Marshal(&model.WSMessage{
		Type: msg.Type,
		Data: privateMsgs,
	})
	if err != nil {
		return cs.sendError(msg.Type, "failed to marshal message", err)
	}
	return resMsg, nil
}

func (cs *chatService) handlePrivateUnread(uid string, msg *model.WSMessage) ([]byte, error) {
	dataBytes, err := json.Marshal(msg.Data)
	if err != nil {
		return cs.sendError(msg.Type, "failed to marshal message data", err)
	}

	var unreadHistoryReq model.UnreadHistoryRequest
	if err := json.Unmarshal(dataBytes, &unreadHistoryReq); err != nil {
		return cs.sendError(msg.Type, "failed to unmarshal history request", err)
	}

	instance := database.GetRedisInstance()
	ctx := context.Background()
	keys := []string{
		"message:unread:" + uid + ":" + unreadHistoryReq.TargetUserId,
	}
	args := []interface{}{
		0,
		-1,
	}
	result, err := instance.Eval(ctx, getUnreadHistoryMessageScript, keys, args)
	if err != nil {
		return cs.sendError(msg.Type, "failed to get unread history message ids from redis", err)
	}
	list, ok := result.([]interface{})
	if !ok {
		return cs.sendError(msg.Type, "failed to parse unread history message list", err)
	}
	privateMsgs := make([]*model.PrivateMsg, len(list))
	for i, v := range list {
		msgSlice, ok := v.([]interface{})
		if !ok {
			return cs.sendError(msg.Type, "failed to parse message list", err)
		}
		var privateMsg model.PrivateMsg
		for j := 0; j < len(msgSlice); j++ {
			switch j {
			case 1:
				id, ok := msgSlice[j].(string)
				if !ok {
					return cs.sendError(msg.Type, "failed to parse message id", err)
				}
				privateMsg.Id = id
			case 3:
				userId, ok := msgSlice[j].(string)
				if !ok {
					return cs.sendError(msg.Type, "failed to parse message userId", err)
				}
				privateMsg.UserId = userId
			case 5:
				toUserId, ok := msgSlice[j].(string)
				if !ok {
					cs.sendError(msg.Type, "failed to parse message toUserId", err)
				}
				privateMsg.ToUserId = toUserId
			case 7:
				content, ok := msgSlice[j].(string)
				if !ok {
					cs.sendError(msg.Type, "failed to parse message content", err)
				}
				privateMsg.Content = content
			case 9:
				t1, ok := msgSlice[j].(string)
				if !ok {
					cs.sendError(msg.Type, "failed to parse message time", err)
				}
				t2, err := strconv.ParseInt(t1, 10, 64)
				if err != nil {
					cs.sendError(msg.Type, "failed to parse time", err)
				}
				privateMsg.Time = t2
			case 11:
				sts, ok := msgSlice[j].(string)
				if !ok {
					cs.sendError(msg.Type, "failed to parse message status", err)
				}
				status, err := strconv.ParseInt(sts, 10, 64)
				if err != nil {
					cs.sendError(msg.Type, "failed to parse status", err)
				}
				privateMsg.Status = int(status)
			}
		}
		privateMsgs[i] = &privateMsg
	}

	resMsg, err := json.Marshal(&model.WSMessage{
		Type: msg.Type,
		Data: privateMsgs,
	})
	if err != nil {
		return cs.sendError(msg.Type, "failed to marshal message", err)
	}
	return resMsg, nil
}

func (cs *chatService) handleGroupMessage(uid string, msg *model.WSMessage) ([]byte, error) {
	dataBytes, err := json.Marshal(msg.Data)
	if err != nil {
		return cs.sendError(msg.Type, "failed to marshal message data", err)
	}

	var groupMsg model.GroupMessage
	if err := json.Unmarshal(dataBytes, &groupMsg); err != nil {
		return cs.sendError(msg.Type, "failed to unmarshal group message", err)
	}

	groupMsg.Id = util.GetID()
	groupMsg.UserId = uid
	groupMsg.Time = time.Now().Unix()

	instance := database.GetRedisInstance()
	ctx := context.Background()
	keys := []string{
		"group:" + groupMsg.GroupId,
		"group:message:" + groupMsg.Id,
	}
	args := []interface{}{
		strconv.FormatFloat(float64(groupMsg.Time), 'f', -1, 64), // ARGV[1]: score
		"Id", groupMsg.Id,
		"GroupId", groupMsg.GroupId,
		"UserId", groupMsg.UserId,
		"Content", groupMsg.Content,
		"Time", strconv.FormatInt(groupMsg.Time, 10),
	}
	_, err = instance.Eval(ctx, setGroupMessageScript, keys, args)
	if err != nil {
		return cs.sendError(msg.Type, "failed to add group message to redis", err)
	}

	resMsg, err := json.Marshal(&model.WSMessage{
		Type: msg.Type,
		Data: "success",
	})
	if err != nil {
		return cs.sendError(msg.Type, "failed to marshal message", err)
	}
	return resMsg, nil
}

func (cs *chatService) handleGroupHistory(uid string, msg *model.WSMessage) ([]byte, error) {
	dataBytes, err := json.Marshal(msg.Data)
	if err != nil {
		return cs.sendError(msg.Type, "failed to marshal message data", err)
	}

	var groupHistoryMsg model.GroupHistoryRequest
	if err := json.Unmarshal(dataBytes, &groupHistoryMsg); err != nil {
		return cs.sendError(msg.Type, "failed to unmarshal group history message", err)
	}

	instance := database.GetRedisInstance()
	ctx := context.Background()
	keys, err := instance.ZRevRange(ctx, "group:"+groupHistoryMsg.GroupId, (groupHistoryMsg.PageNum-1)*groupHistoryMsg.PageSize, groupHistoryMsg.PageNum*groupHistoryMsg.PageSize-1)
	if err != nil {
		return cs.sendError(msg.Type, "failed to get group message keys from redis", err)
	}

	groupMsgs := make([]*model.GroupMessage, len(keys))
	for i, id := range keys {
		hash, err := instance.HGetAll(ctx, id)
		if err != nil {
			return cs.sendError(msg.Type, "failed to get group message by id from redis", err)
		}
		var groupMsg model.GroupMessage
		id, ok := hash["Id"]
		if !ok {
			return cs.sendError(msg.Type, "failed to parse group history message id", err)
		}
		groupMsg.Id = id
		groupId, ok := hash["GroupId"]
		if !ok {
			return cs.sendError(msg.Type, "failed to parse group history message groupId", err)
		}
		groupMsg.GroupId = groupId
		userId, ok := hash["UserId"]
		if !ok {
			return cs.sendError(msg.Type, "failed to parse group history message userId", err)
		}
		groupMsg.UserId = userId
		content, ok := hash["Content"]
		if !ok {
			return cs.sendError(msg.Type, "failed to parse group history message content", err)
		}
		groupMsg.Content = content
		t1, ok := hash["Time"]
		if !ok {
			return cs.sendError(msg.Type, "failed to parse group history message time", err)
		}
		t2, err := strconv.ParseInt(t1, 10, 64)
		if err != nil {
			return cs.sendError(msg.Type, "failed to parse string to int64", err)
		}
		groupMsg.Time = t2
		groupMsgs[i] = &groupMsg
	}

	resMsg, err := json.Marshal(&model.WSMessage{
		Type: msg.Type,
		Data: groupMsgs,
	})
	if err != nil {
		return cs.sendError(msg.Type, "failed to marshal message", err)
	}
	return resMsg, nil
}

func (cs *chatService) sendError(msgType int, msg string, err error) ([]byte, error) {
	log.Printf(msg+": msgType: %d, err: %v", msgType, err)
	resMsg, _ := json.Marshal(&model.WSMessage{
		Type: msgType,
		Data: "internal server error",
	})
	return resMsg, fmt.Errorf(msg, "")
}
