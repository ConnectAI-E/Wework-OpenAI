package wework

import (
	"bytes"
	"time"
	"fmt"
	
	"io/ioutil"
	"net/http"
	"encoding/json"
)

type Wework struct {
	client					*http.Client

	ApiHost string

	CorpId					string
	AppId					int
	AppSecret				string
	Token					string
	EncodingAESKey			string
}

var WeworkApiHost = "https://qyapi.weixin.qq.com"

func NewWeworkClient(CorpId string, AppId int, AppSecret, Token, EncodingAESKey string) *Wework {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	return &Wework{
		client: client,
		ApiHost: WeworkApiHost,

		CorpId: CorpId,
		AppId: AppId,
		AppSecret: AppSecret,
		Token: Token,
		EncodingAESKey: EncodingAESKey,
	}
}

func (wework *Wework) SendRequest(method , uri string,requestBody interface{},responseBody interface{},) error {
	var jsonBody []byte
	var err error

	if requestBody != nil {
		jsonBody, err = json.Marshal(requestBody)
        if err != nil {
            return err
        }
	}

	api := wework.ApiHost + uri

	fmt.Println("request wework url: " + api)
	fmt.Println("payload is:", bytes.NewBuffer(jsonBody))

	req, err := http.NewRequest(method, api, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := wework.client
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	
	err = json.Unmarshal(body, &responseBody)

	if err != nil {
		return err
	}

	return nil
}