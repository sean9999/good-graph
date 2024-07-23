package transport

import "encoding/json"

type Msg struct {
	MsgType string          `json:"msgType"`
	Msg     json.RawMessage `json:"msg"`
	N       int             `json:"n"`
}

func (m Msg) String() string {
	return string(m.Msg)
}

func (m Msg) Set(str string) {
	m.Msg = json.RawMessage(str)
}

func NewMsg(typ string, msg string, n int) Msg {
	m := Msg{
		MsgType: typ,
		Msg:     json.RawMessage(msg),
		N:       n,
	}
	return m
}
