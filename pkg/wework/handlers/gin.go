package handlers

import (
	"net/http"
	"strconv"
	"errors"
	"io/ioutil"
	"encoding/xml"

	"github.com/gin-gonic/gin"

	internalHandlers "wework-vkm/internal/handlers"
	"wework-vkm/internal/initialization"
	"wework-vkm/pkg/wework"
)

func LoadConfig(c *gin.Context) *initialization.Config {
	// 读取 config 配置
	config := c.Keys["config"].(*initialization.Config)

	return config
}

// @see https://developer.work.weixin.qq.com/document/path/90930
// GET 请求，验证回调参数
// POST 请求，接收消息处理
func HandleGinRequest(c *gin.Context) {
	config := LoadConfig(c)

	data, err := parseBody(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var encryptMsg string
	var ReceiveMessage wework.ReceiveMessage

	if (c.Request.Method == "GET") {
		encryptMsg = data
	} else {
		err := xml.Unmarshal([]byte(data), &ReceiveMessage)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		encryptMsg = ReceiveMessage.Encrypt
	}

	// 校验消息的合法性
	if !checkSignature(c, encryptMsg) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 解密消息内容
	encrpytMsgData, err := wework.DecryptMsg(config.WeworkEncodingAESKey, encryptMsg)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	handleGinResponse(c, encrpytMsgData)
}

func handleGinResponse(c *gin.Context, encrpytMsgData *wework.EncrpytMsgData) {
	method := c.Request.Method

	if method == "GET" {
		responseVerifyURL(c, encrpytMsgData)
	} else if method == "POST" {
		handleMessage(c, encrpytMsgData)
	} else {
		c.AbortWithError(http.StatusInternalServerError, errors.New("unknown request body"))
	}
}

func parseBody(c *gin.Context) (data string, err error) {
	method := c.Request.Method

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return "", err
	}

	if method == "GET" {
		data = c.Query("echostr")
	} else {
		data = string(body)
	}

	return data, nil
}

func checkSignature(c *gin.Context, encryptMsg string) bool {
	// 读取 config 配置
	config := LoadConfig(c)

	// 读取 URL 参数
	msgSignature := c.Query("msg_signature")
	nonce := c.Query("nonce")
	timestamp, _ := strconv.ParseInt(c.Query("timestamp"), 10, 64)

	// 校验消息的合法性
	if !wework.CheckSignature(config.WeworkToken, msgSignature, timestamp, nonce, encryptMsg) {
		return false
	}

	return true
}

func responseVerifyURL(c *gin.Context, encrpytMsgData *wework.EncrpytMsgData) {
	c.String(200, encrpytMsgData.Msg)
}

func handleMessage(c *gin.Context, encrpytMsgData *wework.EncrpytMsgData) (err error) {
	internalHandlers.HandleMessage(c, encrpytMsgData)
	return nil
}
