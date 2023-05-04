package handlers

import (
	"context"
	"wework-vkm/pkg/wework"
)

type MsgInfo struct {
	*wework.MsgFromInfo

	msg string
}

type ActionInfo struct {
	handler *MessageHandler
	ctx *context.Context
	msgInfo *MsgInfo
}

type Action interface {
	Execute(data *ActionInfo) bool
}

func chain(data *ActionInfo, actions ...Action) bool {
	for _, v := range actions {
		if !v.Execute(data) {
			return false
		}
	}
	return true
}