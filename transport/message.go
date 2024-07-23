package transport

type Msg struct {
	MsgType string `json:"msgType"`
	Msg     string `json:"msg,omitempty"`
	N       int    `json:"n,omitempty"`
}
