package handlers

import (
	"fmt"
	"context"
	"github.com/gin-gonic/gin"

	"wework-vkm/internal/initialization"

	"wework-vkm/pkg/openai"
	"wework-vkm/pkg/wework"
)

type MessageHandlerInterface interface {
	msgReceivedHandler(ctx context.Context, msg string) error
}

// handlers 所有消息类型类型的处理器
var handlers MessageHandlerInterface

func InitHandlers(gpt *openai.ChatGPT, config initialization.Config) {
	handlers = NewMessageHandler(gpt, config)
}

func HandleMessage(c *gin.Context, encrpytMsgData *wework.EncrpytMsgData) (err error) {
	go func () {
		err = handlers.msgReceivedHandler(c, encrpytMsgData.Msg)
		if err != nil {
			fmt.Println("handle msg error: ", err)
		}
	}()

	c.String(200, "ok")
	return nil
}