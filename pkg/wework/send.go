package wework

import (
	"fmt"
	"log"
	"net/url"
)

type SendMessageResponse struct {
	*Response
	MsgId string `json:"msgid"`
}

func (w *Wework) SendMessage(message, toUserName string) (err error) {
	token, _ := w.GetAccessTokenValue()
	
	values := url.Values{}
	values.Set("access_token", token)
	queryString := values.Encode()

	uri := "/cgi-bin/message/send?" + queryString

	requestBody := &ApplicationTextMessage{
		ApplicationMessageOption: &ApplicationMessageOption{
			MsgType:  "text",
			ToUser: toUserName,
			AgentId: w.AppId,
			Safe: 0,
			EnableIDTrans: 0,
			EnableDuplicateCheck: 0,
			DuplicateCheckInterval: 1800,
		},
		Text: ApplicationTextContent {
			Content: message,
		},
	}
	responseBody := &SendMessageResponse{}
	err = w.SendRequest("POST", uri, requestBody, responseBody)

	if (err != nil) {
		log.Fatalf("failed to send message: %v", err)
		return err
	}

	if responseBody.ErrCode != 0 {
		return fmt.Errorf("send message error: %s", responseBody.ErrMsg)
	}

	return nil
}