package models

import (
	"drop_os_back/dao"
	"fmt"
	"log"
)

type ReceiveMessage struct {
	InsertContent chan *InsertMysql
}

var RM = ReceiveMessage{
	InsertContent: make(chan *InsertMysql),
}

//处理nsq消费者接收函数
func ReceiveToInsert() {
	for {
		fmt.Println("---接收消息并处理---")
		select {
		case MSG := <-RM.InsertContent:
			CheckAndAddMsgForML(MSG)
			_ = InsertMsg2(MSG)
			log.Println("处理完成")
		}
	}
}

//消息列表在发消息的检测和添加机制
func CheckAndAddMsgForML(msg *InsertMysql) {
	msgObj1 := Msgobj{}
	err := dao.DB.Where("from_username=?", msg.FromUsername).Where("to_username=?", msg.ToUsername).Find(&msgObj1).Error
	if err != nil {
		//检索不到相关信息的情况则创建
		var msgData Msgobj
		msgData.FromUsername = msg.FromUsername
		msgData.ToUsername = msg.ToUsername
		dao.DB.Save(&msgData)
	}
	msgObj2 := Msgobj{}
	err = dao.DB.Where("from_username=?", msg.ToUsername).Where("to_username=?", msg.FromUsername).Find(&msgObj2).Error
	if err != nil {
		//检索不到相关信息的情况则创建
		var msgData Msgobj
		msgData.FromUsername = msg.ToUsername
		msgData.ToUsername = msg.FromUsername
		dao.DB.Save(&msgData)
	}
}

//msg publisher init start
func (msgManager *ClientRecMsgManager) RecMsgStart() {
	for {
		select {
		case conn := <-ReceiveMsgManager.Unregister:
			msgCount := ReceiveMsgManager.ClientCount[conn.ID]
			if msgCount == 1 {
				close(conn.Send)
				delete(ReceiveMsgManager.Clients, conn.ID)
				delete(ReceiveMsgManager.ClientCount,conn.ID)
				fmt.Printf("----%s消息推送关闭成功---\n", conn.ID)
			}else {
				msgCount--
				ReceiveMsgManager.ClientCount[conn.ID] = msgCount
			}
		}
	}
}

//根据count计数器添加msger
func AddMsgPusher(stuNum string, cliRec *ClientRecMsg) {
	msgCount := ReceiveMsgManager.ClientCount[stuNum]
	if msgCount == 0 {
		ReceiveMsgManager.ClientCount[stuNum] = 1
	} else {
		msgCount++
		ReceiveMsgManager.ClientCount[stuNum] = msgCount
	}
	ReceiveMsgManager.Clients[stuNum] = cliRec
}
