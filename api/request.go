package api

type WunderListCallbackRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

type RPCMessageWithTimeoutRequest struct {
	Token   string `json:"token"`
	Auth    string `json:"auth"`
	Message string `json:"message"`
	Timeout int    `json:"timeout"`
}
