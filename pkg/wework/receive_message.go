package wework

type ReceiveMessage struct {
	ToUserName string `xml:"ToUserName"`
	Encrypt string `xml:"Encrypt"`
	AgentID string `xml:"AgentID"`
}

type EncrpytMsgData struct {
	Msg string `json:"msg"`
	ReceiveId string `json:"receiveid"`
}

// @see https://developer.work.weixin.qq.com/document/path/90239
type MsgFromInfo struct {
	ToUserName string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime string `xml:"CreateTime"`
	MsgType MsgType `xml:"MsgType"`

	MsgId string `xml:"MsgId"`
	AgentID string `xml:"AgentID"`
}


// @see https://developer.work.weixin.qq.com/document/path/90240
type EventMessage struct {
	*MsgFromInfo
	
	Event string `xml:"Event"`
	EventKey string `xml:"EventKey"`
}

type TextMsgInfo struct {
	*MsgFromInfo

	Content string `xml:"Content"`
}

type ImageMsgInfo struct {
	*MsgFromInfo

	PicUrl string `xml:"PicUrl"`
	MediaId string `xml:"MediaId"`
}

type VoiceMsgInfo struct {
	*MsgFromInfo

	MediaId string `xml:"MediaId"`
	Format string `xml:"Format"`
}

type VideoMsgInfo struct {
	*MsgFromInfo

	MediaId string `xml:"MediaId"`
	ThumbMediaId string `xml:"ThumbMediaId"`
}

type LocationMsgInfo struct {
	*MsgFromInfo

	Location_X string `xml:"Location_X"`
	Location_Y string `xml:"Location_Y"`
	Scale string `xml:"Scale"`
	Label string `xml:"Label"`

	AppType string `xml:"AppType"`
}

type LinkMsgInfo struct {
	*MsgFromInfo

	Title string `xml:"Title"`
	Description string `xml:"Description"`
	Url string `xml:"Url"`
	PicUrl string `xml:"PicUrl"`
}

type MsgType string
const (
	MsgTypeText MsgType = "text"
	MsgTypeImage MsgType = "image"
	MsgTypeVoice MsgType = "voice"
	MsgTypeVideo MsgType = "video"
	MsgTypeLocation MsgType = "location"
	MsgTypeLink MsgType = "link"
	MsgTypeEvent MsgType = "event"
)
