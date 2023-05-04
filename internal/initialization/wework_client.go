package initialization

import (
	"wework-vkm/pkg/wework"
)

var weworkClient *wework.Wework

func LoadWeworkClient(config Config) {
	weworkClient = wework.NewWeworkClient(
		config.WeworkCorpId,
		config.WeworkAppId,
		config.WeworkAppSecret,
		config.WeworkToken,
		config.WeworkEncodingAESKey,
	)
}

func GetWeworkClient() *wework.Wework {
	return weworkClient
}

func GetWeworkAccessTokenValue() (string, error) {
	client := GetWeworkClient()
	token, err := client.GetAccessTokenValue()

	return token, err
}