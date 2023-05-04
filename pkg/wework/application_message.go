package wework

type ApplicationMessageOption struct {
	ToUser string `json:"touser"`
	ToParty string `json:"toparty"`
	ToTag string `json:"totag"`
	MsgType string `json:"msgtype"`
	AgentId int `json:"agentid"`


	Safe int `json:"safe"`
	EnableIDTrans int `json:"enable_id_trans"`
	EnableDuplicateCheck int `json:"enable_duplicate_check"`
	DuplicateCheckInterval int `json:"duplicate_check_interval"`
}

type ApplicationTextContent struct {
	Content string `json:"content"`
}

type ApplicationTextMessage struct {
	*ApplicationMessageOption

	Text ApplicationTextContent `json:"text"`
}
