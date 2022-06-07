package models

import (
	"drop_os_back/util/e"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"sync"
	"time"
)

//用户类
type Client struct {
	ID           string
	SendID       string
	FromUsername string
	ToUsername   string
	Socket       *websocket.Conn
	Send         chan []byte
}

//管理用户登录登出回复广告等
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

//广播类
type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

//管理通项管理
var Manager = ClientManager{
	Clients:    make(map[string]*Client),
	Broadcast:  make(chan *Broadcast, 100),
	Register:   make(chan *Client),
	Reply:      make(chan *Client),
	Unregister: make(chan *Client),
}

//广播写入conn.send结构体
type BroadcastMsg struct {
	Message     string `json:"message"`
	MessageType int    `json:"message_type"`
}

//回复消息
type ReplyMsg struct {
	From        string `json:"from"`
	Code        int    `json:"code"`
	Content     string `json:"content"`
	MessageType int    `json:"message_type"`
	PicShow     bool   `json:"pic_show"`
}

//发送消息
type SendMsg struct {
	Type        int    `json:"type"`
	Content     string `json:"content"`
	MessageType int    `json:"message_type"`
}

//回复消息2(历史消息)
type ReplyMsg2 struct {
	From string `json:"from"`
	//Msg         interface{} `json:"msg"`
	Code        int    `json:"code"`
	Content     string `json:"content"`
	Read        int    `json:"read"`
	CreatAt     string `json:"creat_at"`
	MessageType int    `json:"message_type"`
	PicShow     bool   `json:"pic_show"`
}

type Result struct {
	Start_time int64       `json:"start_time"`
	Msg        SendSortMsg `json:"msg"`
	From       string      `json:"from"`
	Code       int         `json:"code"`
}

type SendSortMsg struct {
	Content     string `json:"content"`
	Read        int    `json:"read"`
	CreatAt     string `json:"creat_at"`
	MessageType int    `json:"message_type"`
}

//管理消息推送
var ReceiveMsgManager = ClientRecMsgManager{
	Clients:     make(map[string]*ClientRecMsg),
	ClientCount: make(map[string]int),
	Broadcast:   make(chan *Broadcast, 100),
	Unregister:  make(chan *ClientRecMsg),
}

//管理消息推送长连接
type ClientRecMsgManager struct {
	Clients map[string]*ClientRecMsg
	//用户计数器，用来缓存websocket延迟关闭删除用户导致消息推送失败
	ClientCount map[string]int
	Broadcast   chan *Broadcast
	Unregister  chan *ClientRecMsg
}

//用户登陆后的长连接结构体
type ClientRecMsg struct {
	ID     string
	Socket *websocket.Conn
	Send   chan []byte
}

//消息推送具体内容
type PublishMsg struct {
	FromUser       string `json:"from_user"`
	MessageContent string `json:"message_content"`
	MessageType    int    `json:"message_type"`
	HeartBeat      int    `json:"HeartBeat"`
}

const month = 60 * 60 * 24 * 30 //一个月30天

//读写锁
var RWMux sync.RWMutex

//互斥锁
var mux sync.Mutex

//websocket向用户写入数据
func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			var message2 BroadcastMsg
			_ = json.Unmarshal(message, &message2)
			replyMsg := &ReplyMsg{
				From:        "you",
				Code:        e.WebsocketSuccessMessage,
				Content:     message2.Message,
				MessageType: message2.MessageType,
			}
			msg, _ := json.Marshal(replyMsg)
			RWMux.Lock()
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			RWMux.Unlock()
		}
	}
}

//websocket用户写入数据
func (c *Client) Read() {
	defer func() {
		mux.Lock()
		Manager.Unregister <- c
		mux.Unlock()
		_ = c.Socket.Close()
	}()

	for {
		c.Socket.PongHandler()
		sendMSg := new(SendMsg)
		err := c.Socket.ReadJSON(&sendMSg)
		if err != nil {
			fmt.Println("数据格式不正确", err)
			//Manager.Unregister <- c
			//_ = c.Socket.Close()
			break
		}
		if sendMSg.Type == 1 {
			mux.Lock()
			Manager.Broadcast <- &Broadcast{
				Client:  c,
				Message: []byte(sendMSg.Content), //发送过来的消息
				Type:    0,
			}
			mux.Unlock()
		} else if sendMSg.Type == 2 {
			//这里将content内获取的数字当作页码
			pageNum, err := strconv.Atoi(sendMSg.Content)
			if err != nil {
				//如果获取错误默认收到的额content的值是1
				pageNum = 1
			}
			results, _ := FindMany(c.SendID, c.ID, pageNum)
			if len(results) == 0 {
				replyMsg := ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "到底了",
				}
				msg, _ := json.Marshal(replyMsg) //序列化
				RWMux.Lock()
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				RWMux.Unlock()
				continue
			}
			if pageNum == 1 {
				results = ReverseResults(results)
			}
			for _, result := range results {
				//history msg
				replyMsg := ReplyMsg2{
					From:        result.From,
					Code:        e.WebsocketHistoryMsg,
					Content:     result.Msg.Content,
					Read:        result.Msg.Read,
					CreatAt:     result.Msg.CreatAt,
					MessageType: result.Msg.MessageType,
					//Msg:  result.Msg,
				}
				msg, _ := json.Marshal(replyMsg) //序列化
				RWMux.Lock()
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				RWMux.Unlock()
			}
		}
	}
}

//将结果倒序
func ReverseResults(s []Result) []Result {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

//websocket管道通信监听
func (manager *ClientManager) Start() {
	for {
		fmt.Println("---监听管道通信---")
		select {
		case conn := <-Manager.Register:
			fmt.Printf("有新连接：%v\n", conn.ID)
			//fmt.Println(&Manager.Register)
			Manager.Clients[conn.ID] = conn //将该连接放到用户管理上
			replyMsg := &ReplyMsg{
				Code:    e.WebsocketSuccess,
				Content: "服务器连接成功",
			}
			msg, _ := json.Marshal(replyMsg)
			RWMux.Lock()
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
			RWMux.Unlock()
		case conn := <-Manager.Unregister:
			fmt.Printf("连接中断%s\n", conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok {
				replyMsg := &ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "连接中断",
				}
				msg, _ := json.Marshal(replyMsg)
				RWMux.Lock()
				_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
				RWMux.Unlock()
				close(conn.Send)
				delete(Manager.Clients, conn.ID)
			}
		case broadcast := <-Manager.Broadcast: //1->2
			//start := time.Now()
			broadcastMessage := string(broadcast.Message)
			message := &BroadcastMsg{
				Message:     broadcastMessage,
				MessageType: broadcast.Type,
			}
			message2, _ := json.Marshal(message)
			SendId := broadcast.Client.SendID //2->1
			flag := false                     //默认对方是不在线的
			//去用户管理里寻找sendid，如果有则证明是该被发送者是在线的，如果没有则不在线
			for id, conn := range Manager.Clients {
				if id != SendId {
					continue
				}
				select {
				case conn.Send <- message2:
					flag = true
				default:
					close(conn.Send)
					delete(Manager.Clients, conn.ID)
				}
			}
			id := broadcast.Client.ID //1->2
			if flag {
				fmt.Println("对方在线")
				replyMsg := &ReplyMsg{
					Code:    e.WebsocketOnlineReply,
					Content: "对方在线应答",
				}
				msg, _ := json.Marshal(replyMsg)
				RWMux.Lock()
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				RWMux.Unlock()
				newInsert := &InsertMysql{
					Id:           id,
					Content:      message.Message,
					Read:         1,
					Expire:       int64(3 * month),
					MessageType:  broadcast.Type,
					FromUsername: broadcast.Client.FromUsername,
					ToUsername:   broadcast.Client.ToUsername,
				}
				go func(insert *InsertMysql) {
					mux.Lock()
					RM.InsertContent <- insert
					mux.Unlock()
					//_ = InsertMsg2(newInsert)
				}(newInsert)
			} else {
				fmt.Println("对方不在线")
				replyMsg := &ReplyMsg{
					Code:    e.WebsocketOfflineReply,
					Content: "对方不在线回答",
				}
				msg, _ := json.Marshal(replyMsg)
				RWMux.Lock()
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				RWMux.Unlock()
				newInsert := &InsertMysql{
					Id:           id,
					Content:      message.Message,
					Read:         0,
					Expire:       int64(3 * month),
					MessageType:  broadcast.Type,
					FromUsername: broadcast.Client.FromUsername,
					ToUsername:   broadcast.Client.ToUsername,
				}
				//建立goroutine向不在线但登录的用户推送消息提醒
				go func(fromUser, content string, messageType int) {
					fmt.Println("异步消息推送")
					for _, v := range ReceiveMsgManager.Clients {
						if broadcast.Client.ToUsername == v.ID {
							publishMsg := PublishMsg{
								FromUser:       fromUser,
								MessageContent: content,
								MessageType:    messageType,
							}
							pubMsg, _ := json.Marshal(publishMsg)
							mux.Lock()
							ReceiveMsgManager.Clients[v.ID].Send <- pubMsg
							mux.Unlock()
						}
					}
				}(broadcast.Client.FromUsername, message.Message, message.MessageType)
				go func(insert *InsertMysql) {
					mux.Lock()
					RM.InsertContent <- insert
					mux.Unlock()
					//_ = InsertMsg2(newInsert)
				}(newInsert)
			}
		}
	}
}

func (r *ClientRecMsg) PushMsg() {
	defer func() {
		fmt.Println("close succ")
		_ = r.Socket.Close()
	}()
	for {
		r.Socket.PongHandler()
		select {
		case msg, ok := <-r.Send:
			if !ok {
				_ = r.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				fmt.Println("msgPush close succ")
				return
			}
			RWMux.Lock()
			_ = r.Socket.WriteMessage(websocket.TextMessage, msg)
			RWMux.Unlock()
		}
	}
}

func (c *ClientRecMsg) CheckOnline() {
	defer func() {
		mux.Lock()
		ReceiveMsgManager.Unregister <- c
		mux.Unlock()
		_ = c.Socket.Close()
	}()

	for {
		PushMsg := struct {
			HeartBeat int
		}{1}
		msg, _ := json.Marshal(PushMsg)
		err := c.Socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("check websocket close")
			break
		}
		log.Println("websocket heartbeat")
		time.Sleep(20 * time.Second)
	}
}
