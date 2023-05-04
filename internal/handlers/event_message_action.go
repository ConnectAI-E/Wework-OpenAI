package handlers

import (
	"fmt"
	"encoding/xml"

	"wework-vkm/pkg/openai"
	"wework-vkm/pkg/wework"
)

type MessageAction struct {
}

func (*MessageAction) Execute(a *ActionInfo) bool {
	MsgType := a.msgInfo.MsgType
	if MsgType != "text" {
		return false
	}

	textMsg := &wework.TextMsgInfo{}
	err := xml.Unmarshal([]byte(a.msgInfo.msg), textMsg)
	if err != nil {
		fmt.Println(err)
		return false
	}

	msg := a.handler.sessionCache.GetMsg(a.msgInfo.MsgId)
	msg = append(msg, openai.Messages{
		Role: "user", Content: textMsg.Content,
	})
	completions, err := a.handler.gpt.Completions(msg)
	if err != nil {
		replyMsg(*a.ctx, fmt.Sprintf("🤖️：消息机器人摆烂了，请稍后再试～\n错误信息: %v", err), a.msgInfo.FromUserName)
		return false
	}
	msg = append(msg, completions)
	a.handler.sessionCache.SetMsg(a.msgInfo.MsgId, msg)
	//if new topic
	if len(msg) == 2 {
		//fmt.Println("new topic", msg[1].Content)
		// sendNewTopicCard(*a.ctx, a.msgInfo.MsgId, a.info.msgId,
		// 	completions.Content)
		// return false
	}
	err = replyMsg(*a.ctx, completions.Content, a.msgInfo.FromUserName)
	if err != nil {
		replyMsg(*a.ctx, fmt.Sprintf("🤖️：消息机器人摆烂了，请稍后再试～\n错误信息: %v", err), a.msgInfo.FromUserName)
		return false
	}
	return true
}