package models

import (
	"drop_os_back/dao"
	"strconv"
	"time"
)

type Discussion struct {
	Id              int       `json:"id" gorm:"primary_key"`
	Username        string    `json:"username"`
	Datetime        time.Time `json:"datetime"`
	Content         string    `json:"content"`
	Parent_id       int       `json:"parent_id"` //父id
	Level           int       `json:"level"`
	Belong_id       int       `json:"belong_id"` //所属父亲文章id
	Pic_address     string    `json:"pic_address"`
	Title           string    `json:"title"`
	Parent_username string    `json:"parent_username"`
	Classification  string    `json:"classification"`
	Like_num        int       `json:"like_num"`
	Content_num     int       `json:"content_num"`
	Father_con_id   int       `json:"father_con_id"`
}

//通过文章id获取文章对应的评论
func GetDiscussionByID(articleID int) (discussionSli []*Discussion) {
	newDiscussion := []*Discussion{}
	dao.DB.Where("belong_id=?", articleID).Find(&newDiscussion)
	return newDiscussion
}

//制作评论信息
func MakeDiscussionData(articleId int) (dataList []interface{}) {
	discussionSli := GetDiscussionByID(articleId)
	dataList = getJsonTree(discussionSli, 0, articleId)
	return dataList
}

//制作评论区返回的json数据
func getJsonTree(data []*Discussion, level, fatherId int) (dataList []interface{}) {
	level++
	dataList = []interface{}{}
	datalength := len(data)
	//fmt.Println(level)
	for i := 0; i < datalength; i++ {
		if data[i].Level == level && data[i].Parent_id == fatherId {
			discussion := map[string]interface{}{}
			discussion["content"] = data[i].Content
			discussion["id"] = data[i].Id
			discussion["userid"] = data[i].Username
			discussion["parent_id"] = data[i].Parent_id
			discussion["datetime"] = data[i].Datetime
			discussion["parent_username"] = data[i].Parent_username
			discussion["like_num"] = data[i].Like_num
			discussion["children"] = getJsonTree(data, level, data[i].Id)
			dataList = append(dataList, discussion)
		}
	}
	return dataList
}

//发表帖子(文章部分)
func PublishArt(stuNum, comment, classification, artPicAdd, artTitle string, nowTime time.Time) (err error) {
	disData := Discussion{}
	disData.Username = stuNum
	disData.Content = comment
	disData.Datetime = nowTime
	disData.Parent_id = 0
	disData.Level = 0
	disData.Belong_id = 0
	disData.Like_num = 0
	disData.Content_num = 0
	disData.Father_con_id = 0
	disData.Title = artTitle
	disData.Classification = classification
	disData.Pic_address = artPicAdd
	err = dao.DB.Save(&disData).Error
	if err != nil {
		return err
	}
	return nil
}

//获取用户名对应的最新一条记录
func GetNewest(username string) (data Discussion) {
	disData := Discussion{}
	dao.DB.Where("username=?", username).Last(&disData)
	return disData
}

//发表帖子（图片部分）
func PublishArtPic(artId int, saveAddress string) (err error) {
	disData := Discussion{}
	err1 := dao.DB.Where("id=?", artId).Find(&disData).Error
	if err1 != nil {
		return err1
	}
	disData.Pic_address = saveAddress
	err2 := dao.DB.Save(&disData).Error
	if err2 != nil {
		return err2
	}
	return nil
}

//发表评论
func PublishComment(id int, username, comment string, time time.Time) (err error) {
	fatherData := Discussion{}
	err1 := dao.DB.Where("id=?", id).Find(&fatherData).Error
	if err1 != nil {
		return err1
	}
	sonData := Discussion{}
	if fatherData.Belong_id == 0 {
		sonData.Belong_id = fatherData.Id
	} else {
		sonData.Belong_id = fatherData.Belong_id
	}
	sonData.Parent_username = fatherData.Username
	sonData.Parent_id = fatherData.Id
	if fatherData.Level != 2 {
		sonData.Level = fatherData.Level + 1
	} else {
		sonData.Level = 2
	}
	switch fatherData.Level {
	case 1:
		sonData.Father_con_id = fatherData.Id
		break
	case 2:
		sonData.Father_con_id = fatherData.Father_con_id
		break
	case 0:
		sonData.Father_con_id = 0
		break
	default:
		sonData.Father_con_id = 0
	}
	sonData.Username = username
	sonData.Content = comment
	sonData.Datetime = time
	fatherData.Content_num += 1
	err2 := dao.DB.Save(&fatherData).Error
	if err2 != nil {
		return err2
	}
	err3 := dao.DB.Save(&sonData).Error
	if err3 != nil {
		return err3
	}
	return nil
}

//获取所有帖子的内容(存在分页)
func GetAllArt(pageNum int, classification string) (data []interface{}) {
	var allDis []*Discussion
	pageNumStr := strconv.Itoa((pageNum - 1) * 15)
	sql := ""
	if classification == "全部" {
		sql = "SELECT * FROM `discussions`  WHERE level = 0 ORDER BY id desc LIMIT " + pageNumStr + ",15; "
	} else {
		sql = "SELECT * FROM `discussions`  WHERE level = 0 and classification = '" + classification + "' ORDER BY id desc LIMIT " + pageNumStr + ",30; "
	}
	dao.DB.Raw(sql).Scan(&allDis)
	dataList := []interface{}{}
	for _, item := range allDis {
		dis := map[string]interface{}{}
		dis["id"] = item.Id
		dis["username"] = item.Username
		dis["title"] = item.Title
		dis["datetime"] = item.Datetime.Format("2006-01-02 15:04:05")
		contentRune := []rune(item.Content)
		if len(contentRune) >= 80 {
			contentStr := string(contentRune[:80])
			dis["content"] = contentStr
		} else {
			dis["content"] = item.Content
		}
		dis["pic_address"] = item.Pic_address
		dis["classification"] = item.Classification
		dis["like_num"] = item.Like_num
		dis["content_num"] = item.Content_num
		dis["popup"] = false
		dataList = append(dataList, dis)
	}
	return dataList
}

//用户获取所有帖子的内容(存在分页)
func UserGetAllArt(pageNum int, username string) (data []interface{}) {
	var allDis []*Discussion
	pageNumStr := strconv.Itoa((pageNum - 1) * 15)
	sql := "SELECT * FROM `discussions`  WHERE level = 0 and username = '" + username + "' ORDER BY id desc LIMIT " + pageNumStr + ",15; "
	dao.DB.Raw(sql).Scan(&allDis)
	dataList := []interface{}{}
	for _, item := range allDis {
		dis := map[string]interface{}{}
		dis["id"] = item.Id
		dis["username"] = item.Username
		dis["title"] = item.Title
		dis["datetime"] = item.Datetime.Format("2006-01-02 15:04:05")
		contentRune := []rune(item.Content)
		if len(contentRune) >= 80 {
			contentStr := string(contentRune[:80])
			dis["content"] = contentStr
		} else {
			dis["content"] = item.Content
		}
		dis["pic_address"] = item.Pic_address
		dis["classification"] = item.Classification
		dis["like_num"] = item.Like_num
		dis["content_num"] = item.Content_num
		dis["popup"] = false
		dataList = append(dataList, dis)
	}
	return dataList
}

//获取单个帖子的信息
func GetArticleIn(artId int) (data map[string]interface{}) {
	artInfo := Discussion{}
	dao.DB.Where("id=?", artId).Find(&artInfo)
	dis := map[string]interface{}{}
	dis["id"] = artInfo.Id
	dis["username"] = artInfo.Username
	dis["datetime"] = artInfo.Datetime.Format("2006-01-02 15:04:05")
	dis["title"] = artInfo.Title
	dis["content"] = artInfo.Content
	dis["pic_address"] = artInfo.Pic_address
	dis["like_num"] = artInfo.Like_num
	dis["content_num"] = artInfo.Content_num
	return dis
}

//点赞
func ControlLike(id int, beahavior string) error {
	artInfo := Discussion{}
	err1 := dao.DB.Where("id=?", id).Find(&artInfo).Error
	if err1 != nil {
		return err1
	}
	if beahavior == "like" {
		artInfo.Like_num = artInfo.Like_num + 1
	} else {
		artInfo.Like_num = artInfo.Like_num - 1
	}
	err2 := dao.DB.Save(&artInfo).Error
	if err2 != nil {
		return err2
	}
	return nil
}

func GetDiscussionArea(pageNum, belong int) []*Discussion {
	newDiscussion := []*Discussion{}
	belongStr := strconv.Itoa(belong)
	pageNumStr := strconv.Itoa((pageNum - 1) * 10)
	sql := "SELECT * FROM `discussions`  WHERE belong_id = " + belongStr + " and level = 1 ORDER BY id desc LIMIT " + pageNumStr + ",10; "
	//fmt.Println(sql)
	dao.DB.Raw(sql).Scan(&newDiscussion)
	return newDiscussion
}

func DealFatherDiscussionArea(pageNum, belong int) []interface{} {
	fatherDis := GetDiscussionArea(pageNum, belong)
	dataList := []interface{}{}
	for _, faItem := range fatherDis {
		discussion := map[string]interface{}{}
		discussion["content"] = faItem.Content
		discussion["id"] = faItem.Id
		discussion["userid"] = faItem.Username
		discussion["parent_id"] = faItem.Parent_id
		discussion["datetime"] = faItem.Datetime.Format("2006-01-02 15:04:05")
		discussion["parent_username"] = "null"
		faIdStr := strconv.Itoa(faItem.Id)
		faBeStr := strconv.Itoa(faItem.Belong_id)
		childrenInfo := DealSonDiscussionArea(faIdStr, faBeStr)
		discussion["children"] = childrenInfo
		discussion["childrenShow"] = true
		discussion["addStatus"] = false
		discussion["childPageNum"] = 1
		dataList = append(dataList, discussion)
	}
	return dataList
}

func DealSonDiscussionArea(parent_id, belong_id string) []interface{} {
	sonDis := GetDiscussionByPAB(parent_id, belong_id)
	dataList := []interface{}{}
	for _, faItem := range sonDis {
		discussion := map[string]interface{}{}
		discussion["content"] = faItem.Content
		discussion["id"] = faItem.Id
		discussion["userid"] = faItem.Username
		discussion["parent_id"] = faItem.Parent_id
		discussion["datetime"] = faItem.Datetime.Format("2006-01-02 15:04:05")
		discussion["parent_username"] = faItem.Parent_username
		dataList = append(dataList, discussion)
	}
	return dataList
}

func GetDiscussionByPAB(parent_id, belong_id string) []*Discussion {
	newDiscussion := []*Discussion{}
	sql := "SELECT * FROM `discussions`  WHERE belong_id = " + belong_id + " and  father_con_id = " + parent_id + " ORDER BY id desc LIMIT 0,5; "
	dao.DB.Raw(sql).Scan(&newDiscussion)
	return newDiscussion
}

func GetMoreDis(pageNum int, belong, parent string) []interface{} {
	newDiscussion := []*Discussion{}
	pageNumStr := strconv.Itoa((pageNum - 1) * 5)
	sql := "SELECT * FROM discussions  WHERE belong_id = " + belong + " and father_con_id = " + parent + " ORDER BY id desc LIMIT " + pageNumStr + ",5; "
	//sql := "SELECT * FROM `times`  WHERE username = '" + username + "' LIMIT " + pageNumStr + ",10 order by userid desc; "
	//fmt.Println(sql)
	dataList := []interface{}{}
	dao.DB.Raw(sql).Scan(&newDiscussion)
	for _, item := range newDiscussion {
		discussion := map[string]interface{}{}
		discussion["content"] = item.Content
		discussion["id"] = item.Id
		discussion["userid"] = item.Username
		discussion["parent_id"] = item.Parent_id
		discussion["datetime"] = item.Datetime.Format("2006-01-02 15:04:05")
		discussion["parent_username"] = item.Parent_username
		dataList = append(dataList, discussion)
	}
	return dataList
}

func DeleteArticle(artId int) error {
	err := dao.DB.Where("id=?", artId).Delete(new(Discussion)).Error
	if err != nil {
		return err
	}
	dao.DB.Where("belong_id=?", artId).Delete(new(Discussion))
	return nil
}
