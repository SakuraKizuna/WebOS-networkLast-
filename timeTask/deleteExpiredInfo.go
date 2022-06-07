package timeTask

import (
	"drop_os_back/models"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup


//删除过期消息，消息有效时长，三个月
//删除数据库对应的本地图片
func deleteExpiredInfo() {
	nowUnix := time.Now().Unix()
	nowUnixStr := strconv.Itoa(int(nowUnix))
	PicData := models.FindExpiredChatPic(nowUnixStr)
	if len(PicData) == 0{
		DeleteExpiredMsgPic(PicData)
		models.DeleteExpiredMsg(nowUnixStr)
	}
}


func DeleteExpiredMsgPic(picData []*models.Trainer) {
	wg.Add(len(picData))
	for _, pic := range picData {
		picUrl := pic.Content
		pos := strings.Index(picUrl, "/uploadPic")
		finalPicUrl := "." + picUrl[pos:]
		go func(finalPicUrl string) {
			err := os.Remove(finalPicUrl)
			if err != nil {
				log.Println("file remove Error!")
				log.Printf("%s", err)
			} else {
				log.Println("file remove OK!")
			}
			wg.Done()
		}(finalPicUrl)
	}
	wg.Wait()
}
