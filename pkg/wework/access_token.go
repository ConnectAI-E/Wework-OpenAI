package wework

import (
	"log"
	"os"
	"time"
	"io/ioutil"
	"net/url"
	"encoding/json"
)

type AccessTokenRequestBody struct {
	CorpId string `json:"corpid"`
	CorpSecret string `json:"corpsecret"`
}

type AccessTokenResponseBody struct {
	*Response

	AccessToken string `json:"access_token"`
	ExpiresIn int `json:"expires_in"`
	ExpireAt time.Time `json:"expired_at"`
}

const (
	AccessTokenFile = "wework_access_token.json"
)

func (t *AccessTokenResponseBody) setExpireAt() {
	// 当前时间加上expires_in的时间减少5分钟
	t.ExpireAt = time.Now().Add(time.Second * time.Duration(t.ExpiresIn-300))
}

func (t *AccessTokenResponseBody) IsExpired() bool {
    return t.ExpireAt.Before(time.Now())
}

func (t *AccessTokenResponseBody) getToken() string {
    return t.AccessToken
}

func (w *Wework) GetAccessTokenValue() (token string, err error) {
	AccessToken, err := w.GetAccessToken()
	if err != nil {
		return "", err
	}

	return AccessToken.getToken(), nil
}

func (w *Wework) GetAccessToken() (AccessToken *AccessTokenResponseBody, err error) {
	// 判断 token 缓存是否存在
	if !w.IsAccessTokenFileExist() {
		// 缓存文件不存在，获取新的 token
		AccessToken, err = w.GetNewAccessToken()
		if err != nil {
            log.Fatalf("failed to get wework access token: %v", err)
            return nil, err
        }

	} else {
		// 缓存文件存在，读取缓存文件
		AccessToken, err = w.ReadAccessTokenFromFile()
		if err != nil {
            log.Fatalf("failed to read access token from file: %v", err)
            return nil, err
        }

		// 判断 token 是否有效
        if AccessToken.IsExpired() {
            // 缓存文件中的 token 已经过期，获取新的 token
            AccessToken, err = w.GetNewAccessToken()
            if err != nil {
                log.Fatalf("failed to get wework access token: %v", err)
                return nil, err
            }
        }
	}
	
	return AccessToken, nil
}

func (w *Wework) IsAccessTokenFileExist() bool {
    _, err := os.Stat(AccessTokenFile)
    return !os.IsNotExist(err)
}

func (w *Wework) GetNewAccessToken() (AccessToken *AccessTokenResponseBody, err error) {
	requestBody := &AccessTokenRequestBody{
		CorpId: w.CorpId,
		CorpSecret: w.AppSecret,
	}
	responseBody := &AccessTokenResponseBody{}

	// 请求参数
	values := url.Values{}
	values.Set("corpid", requestBody.CorpId)
	values.Set("corpsecret", requestBody.CorpSecret)
	queryString := values.Encode()

	uri := "/cgi-bin/gettoken?" + queryString

	err = w.SendRequest("GET", uri, requestBody, responseBody)
	if (err != nil) {
		log.Fatalf("failed to get wework access token: %v", err)
		return nil, err
	}

	responseBody.setExpireAt()

	w.SaveAccessTokenToFile(responseBody)

	return responseBody, nil
}

func (w *Wework) SaveAccessTokenToFile(responseBody *AccessTokenResponseBody) error {
	data, err := json.Marshal(responseBody)
	if err != nil {
		return err
	}

	// 以 0600 权限打开文件
	f, err := os.OpenFile(AccessTokenFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (w *Wework) ReadAccessTokenFromFile() (*AccessTokenResponseBody, error) {
    // 打开缓存文件
    file, err := os.Open(AccessTokenFile)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // 读取文件内容
    bytes, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, err
    }

    // 反序列化为 AccessTokenResponseBody
    responseBody := &AccessTokenResponseBody{}
    err = json.Unmarshal(bytes, responseBody)
    if err != nil {
        return nil, err
    }

    return responseBody, nil
}