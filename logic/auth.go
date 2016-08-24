package main

import (
	"strconv"
)

// developer could implement "Auth" interface for decide how get userId, or roomId
type Auther interface {
	Auth(body []byte) (userId int64, roomId int32, err error)
}

type DefaultAuther struct {
}

func NewDefaultAuther() *DefaultAuther {
	return &DefaultAuther{}
}

func (a *DefaultAuther) Auth(body []byte) (userId int64, roomId int32, err error) {
	if userId, err = strconv.ParseInt(string(body), 10, 64); err != nil {

	} else {
		roomId = 1 // only for debug
	}
	return
}
