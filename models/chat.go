package models

import (
	"drop_os_back/dao"
	"strconv"
	"sync"
	"time"
)

//TODO 添加message_type 字段，1为文本消息，2为图片，考虑增加可发送的文档（3）
type Trainer struct {
	Id           int    `json:"id" gorm:"primary_key"`
	Userid       string `json:"-"`            //用户名
	Content      string `json:"content"`      // 内容
	Start_time   int64  `json:"start_time"`   // 创建时间
	End_time     int64  `json:"-"`            // 过期时间
	Read         int    `json:"read"`         // 已读
	Message_type int    `json:"message_type"` //消息类型
	FromUsername string `json:"-"`            //发送者
	ToUsername   string `json:"-"`            //接受者
}

type InsertMysql struct {
	Id           string `json:"id"`
	Content      string `json:"content"`
	Read         int    `json:"read"`
	Expire       int64  `json:"expire"`
	MessageType  int    `json:"message_type"`
	FromUsername string `json:"from_username"`
	ToUsername   string `json:"to_username"`
}

func InsertMsg(userid string, content string, read int, expire int64) (err error) {
	comment := Trainer{
		Userid:     userid,
		Content:    content,
		Start_time: time.Now().Unix(),
		End_time:   time.Now().Unix() + expire,
		Read:       read,
	}
	err = dao.DB.Save(&comment).Error
	if err != nil {
		return err
	}
	return
}

func InsertMsg2(msg *InsertMysql) (err error) {
	comment := Trainer{
		Userid:       msg.Id,
		Content:      msg.Content,
		Start_time:   time.Now().Unix(),
		End_time:     time.Now().Unix() + msg.Expire,
		Read:         msg.Read,
		Message_type: msg.MessageType,
		FromUsername: msg.FromUsername,
		ToUsername:   msg.ToUsername,
	}
	err = dao.DB.Save(&comment).Error
	if err != nil {
		return err
	}
	return
}

var wg sync.WaitGroup

//添加头像url
type QELM struct {
	Trainer
	ToUserHeadPic  string `json:"to_user_head_pic"`
	ToUsername     string `json:"to_username"`
	ToUserId       int    `json:"to_user_id"`
	ToUserRealname string `json:"to_user_realname"`
	//ToUserNickname string `json:"to_user_nickname"`
}

//sql优化，不用in
//SELECT * FROM `trainers` where userid = 'oZ65W5TklL3gWTCLTllMfiXu97ig->20062111' UNION ALL SELECT * FROM `trainers` where userid = '20062111->oZ65W5TklL3gWTCLTllMfiXu97ig'  ORDER BY id desc LIMIT 0,5
func FindMany(sendID, id string, pageNum int) (results []Result, err error) {
	pageSize := 5
	pageSizeStr := strconv.Itoa(pageSize)
	pageNumStr := strconv.Itoa((pageNum - 1) * pageSize)
	var resultAll []Trainer //存放id和sendid的一些信息
	sql := "SELECT * FROM `trainers` where userid = '" + id + "' UNION ALL SELECT * FROM `trainers` where userid = '" + sendID + "'  ORDER BY id desc  LIMIT " + pageNumStr + "," + pageSizeStr
	//sql := "SELECT * FROM `trainers` where userid in ('" + id + "','" + sendID + "') ORDER BY id desc LIMIT " + pageNumStr + "," + pageSizeStr
	//fmt.Println(sql)
	dao.DB.Raw(sql).Scan(&resultAll)
	results, _ = AppendAndSort(resultAll, sendID, id)
	return
}

func AppendAndSort(resultAll []Trainer, sendID, id string) (results []Result, err error) {
	for _, r := range resultAll {
		start_time := time.Unix(r.Start_time, 0).Format("2006-01-02 15:04:05")
		sendSort := SendSortMsg{ //构造返回的msg
			Content:     r.Content,
			Read:        r.Read,
			CreatAt:     start_time,
			MessageType: r.Message_type,
		}
		var result Result
		if r.Userid == id {
			result = Result{ //构造返回所有的内容，包括传送者
				Start_time: r.Start_time,
				Msg:        sendSort,
				From:       "me",
			}
		} else {
			result = Result{ //构造返回所有的内容，包括传送者
				Start_time: r.Start_time,
				Msg:        sendSort,
				From:       "you",
			}
		}
		results = append(results, result)
	}
	return
}

//批量设置已读
func BatchRead(fromUsername, toUsername string) {
	dao.DB.Table("trainers").Where("from_username=?", fromUsername).Where("to_username=?", toUsername).Where("`read`=?", 0).Update(map[string]interface{}{"read": 1})
}

func QueryEveryLastMsg(personList []Msgobj) []QELM {
	var MsgList []QELM
	wgInt := len(personList)
	wg.Add(wgInt)
	for _, person := range personList {
		go func(person Msgobj) {
			personChat := QELM{}
			sql := "SELECT * FROM trainers WHERE from_username = '" + person.ToUsername + "' AND to_username ='" + person.FromUsername + "' UNION SELECT * FROM trainers WHERE from_username = '" + person.FromUsername + "' AND to_username ='" + person.ToUsername + "' ORDER BY id DESC LIMIT 1"
			dao.DB.Raw(sql).Scan(&personChat)
			personInfo := User{}
			personInfoSql := "select userid,head_pic,realname from users where username = '" + person.ToUsername + "'"
			dao.DB.Raw(personInfoSql).Scan(&personInfo)
			personChat.Id = person.Id
			personChat.ToUserHeadPic = personInfo.Head_pic
			personChat.ToUserId = personInfo.Userid
			personChat.ToUsername = person.ToUsername
			personChat.ToUserRealname = personInfo.Realname
			MsgList = append(MsgList, personChat)
			wg.Done()
		}(person)
	}
	wg.Wait()
	return MsgList
}

//查找过期聊天图片
func FindExpiredChatPic(nowUnixStr string) []*Trainer {
	var Pics []*Trainer
	sql := "select content from trainers where message_type = 1 and end_time < " + nowUnixStr
	dao.DB.Raw(sql).Scan(&Pics)
	return Pics
}

//删除过期聊天记录
func DeleteExpiredMsg(nowUnixStr string) {
	sql := "DELETE FROM trainers WHERE end_time < " + nowUnixStr + " ORDER BY id LIMIT 10000;"
	dao.DB.Raw(sql)
}
