package main

import (
	"encoding/json"
	"fmt"
	log "github.com/thinkboy/log4go"
	"strconv"
	"yf-im/yfgoim"
)

const (
	ERR_ARGS_INVALID  string = "ERR_ARGS_INVALID"
	ERR_SERVER_ERROR  string = "ERR_SERVER_ERROR"
	ERR_SIGN_ERROR    string = "ERR_SIGN_ERROR"
	ERR_TOKEN_INVALID string = "ERR_TOKEN_INVALID"
)

const (
	TYPE_USER      = iota // 0:个人消息
	TYPE_ROOM             // 1：房间消息
	TYPE_BROADCAST        // 2:广播消息
)

type MsgInfo struct {
	Msg    string `json:"msg"`
	Type   int    `json:"type"`
	UserId int64  `json:"user_id"`
	RoomId string `json:"roomid"`
	Token  string `json:"token"`
	Sign   string `json:"sign"`
}

func SendMsg2User(msgInfo *MsgInfo) (err error) {
	if msgInfo.UserId < 1 {
		log.Error("invalid user_id: %d", msgInfo.UserId)
		err = fmt.Errorf(ERR_ARGS_INVALID)
		return
	}
	if len(msgInfo.Msg) < 1 {
		log.Error("invalid msg: %s", string(msgInfo.Msg))
		err = fmt.Errorf(ERR_ARGS_INVALID)
		return
	}

	btBody := []byte(msgInfo.Msg)
	subKeys := genSubKey(msgInfo.UserId)
	for serverId, keys := range subKeys {
		if err = mpushKafka(serverId, keys, btBody); err != nil {
			log.Error("mpushKafka failed! err: %v", err)
			err = fmt.Errorf(ERR_SERVER_ERROR)
			return
		}
	}

	return
}

func SendMsg2Room(msgInfo *MsgInfo) (err error) {
	if len(msgInfo.RoomId) < 1 {
		log.Error("invalid roomid: %s", msgInfo.RoomId)
		err = fmt.Errorf(ERR_ARGS_INVALID)
		return
	}
	if len(msgInfo.Msg) < 1 {
		log.Error("invalid msg: %s", string(msgInfo.Msg))
		err = fmt.Errorf(ERR_ARGS_INVALID)
		return
	}
	// btBody := []byte(msgInfo.Msg)
	// {\"test\": 1}
	sender, err := yfgoim.CheckToken(msgInfo.Token)
	if err != nil {
		log.Error("CheckToken(%s) failed! err: %v", msgInfo.Token, err)
		return err
	}
	btBody := []byte(fmt.Sprintf(`{"msg":"%s", "type": %d, "sender": %d}`, msgInfo.Msg, msgInfo.Type, sender))
	// TODO: 解密RoomId
	var rid int
	if rid, err = strconv.Atoi(msgInfo.RoomId); err != nil {
		log.Error("strconv.Atoi(\"%s\") error(%v)", msgInfo.RoomId, err)
		err = fmt.Errorf(ERR_ARGS_INVALID)
		return
	}
	ensure := false
	if err = broadcastRoomKafka(int32(rid), btBody, ensure); err != nil {
		log.Error("broadcastRoomKafka(\"%s\",\"%s\",\"%d\") error(%s)", rid, string(btBody), ensure, err)
		err = fmt.Errorf(ERR_SERVER_ERROR)
		return
	}
	return
}

func sendmsg(body []byte) (err error) {
	var msgInfo MsgInfo
	log.Info("yf-auth recv a Send msg: %v", string(body))
	err = json.Unmarshal(body, &msgInfo)
	if err != nil {
		log.Error("Recv a invalid Msg: %v", string(body))
		err = fmt.Errorf(ERR_ARGS_INVALID)
		return
	}

	// TODO: Check Sign.

	switch msgInfo.Type {
	case TYPE_USER:
		return SendMsg2User(&msgInfo)
	case TYPE_ROOM:
		return SendMsg2Room(&msgInfo)
	default:
		return fmt.Errorf(ERR_ARGS_INVALID)
	}
	return nil
}
