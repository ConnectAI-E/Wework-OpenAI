package handlers

import (
	"context"
	"encoding/xml"

	"wework-vkm/internal/initialization"

	"wework-vkm/pkg/openai"
	"wework-vkm/pkg/services"
	"wework-vkm/pkg/wework"
)

type MessageHandler struct {
	sessionCache services.SessionServiceCacheInterface
	msgCache     services.MsgCacheInterface
	gpt          *openai.ChatGPT
	config       initialization.Config
}

func (m MessageHandler) msgReceivedHandler(ctx context.Context, msg string) error {
	MsgFromInfo := &wework.MsgFromInfo{}
	err := xml.Unmarshal([]byte(msg), MsgFromInfo)
	if err != nil {
		return err
	}

	msgInfo := &MsgInfo{
		MsgFromInfo: MsgFromInfo,
		msg: msg,
	}

	data := &ActionInfo{
		ctx:     &ctx,
		handler: &m,
		msgInfo:    msgInfo,
	}

	actions := []Action{
		// &ProcessedUniqueAction{}, //避免重复处理
		// &ProcessMentionAction{},  //判断机器人是否应该被调用
		// &AudioAction{},           //语音处理
		// &EmptyAction{},           //空消息处理
		// &ClearAction{},           //清除消息处理
		// &PicAction{},             //图片处理
		// &RoleListAction{},        //角色列表处理
		// &HelpAction{},            //帮助处理
		// &BalanceAction{},         //余额处理
		// &RolePlayAction{},        //角色扮演处理
		// &VkmAction{},             //知识库处理
		&MessageAction{},         //消息处理

	}
	chain(data, actions...)
	return nil
}

var _ MessageHandlerInterface = (*MessageHandler)(nil)

func NewMessageHandler(gpt *openai.ChatGPT, config initialization.Config) MessageHandlerInterface {
	return &MessageHandler{
		sessionCache: services.GetSessionCache(),
		msgCache:     services.GetMsgCache(),
		gpt:          gpt,
		config:       config,
	}
}
