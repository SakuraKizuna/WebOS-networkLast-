package e

type Code int

const (
	//通信型号
	WebsocketSuccessMessage = 50001
	WebsocketSuccess        = 50002
	WebsocketEnd            = 50003
	WebsocketOnlineReply    = 50004 //在线应答
	WebsocketOfflineReply   = 50005 //不在线应答
	WebsocketLimit          = 50006
	WebsocketHistoryMsg     = 50007 //历史消息
)

func (c Code) Msg() string {
	return codeMsg[c]
}
