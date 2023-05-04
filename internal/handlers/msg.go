package handlers

import (
	"context"
	"wework-vkm/internal/initialization"
)

func replyMsg(ctx context.Context, msg string, fromUserName string) error {
	weworkClient := initialization.GetWeworkClient()
	err := weworkClient.SendMessage(msg, fromUserName)
	return err
}
