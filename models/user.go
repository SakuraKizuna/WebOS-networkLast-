package models

import (
	"drop_os_back/dao"
	"drop_os_back/util"
	"fmt"
	"strconv"
	"time"
)

type User struct {
	Userid     int    `json:"userid" gorm:"primary_key"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Nickname   string `json:"nickname"`
	Grade      int    `json:"grade"`
	Qq         string `json:"qq"`
	Profession string `json:"profession"`
	Realname   string `json:"realname"`
	Gender     string `json:"gender"`
	Role       string `json:"role"`
	State      string `json:"state"`
	Belong     string `json:"belong"`
	Head_pic   string `json:"head_pic"`
	Createtime string `json:"createtime"`
}

func (user *User) getPersonalInfo(username string) (user2 User) {
	dao.DB.Where("username=?", username).Find(&user)
	return *user
}

func GetUserLists(level int, belong string) (uL []interface{}, err error) {
	var userList []*User
	if level == 0 {
		err = dao.DB.Find(&userList).Error
	} else {
		err = dao.DB.Where("grade=?", level).Where("belong=?", belong).Find(&userList).Error
	}
	temp_intf := []interface{}{}
	for _, b := range userList {
		dict := map[string]interface{}{}
		dict["email"] = b.Qq
		dict["major"] = b.Profession
		dict["name"] = b.Realname
		dict["schoolNumber"] = b.Username
		dict["sex"] = b.Gender
		if b.State == "1" {
			temp_intf = append(temp_intf, dict)
		}
	}
	//fmt.Println(temp_intf)
	return temp_intf, nil
}

//定义获取学号（username）切片
func GetUserSli(level int, belong string) (stuSli []string, err error) {
	var userList []*User
	if level == 0 {
		err = dao.DB.Where("belong=?", belong).Find(&userList).Error
	} else {
		err = dao.DB.Where("grade=?", level).Where("belong=?", belong).Find(&userList).Error
	}
	if err != nil {
		fmt.Println("mysql err:", err)
		return nil, err
	}
	schoolNumSli := []string{}
	for _, k := range userList {
		schoolNumSli = append(schoolNumSli, k.Username)
		//fmt.Println(k)
	}
	return schoolNumSli, nil
}

//定义获取个人信息的姓名和性别
func GetUserSingle(username string) (dict map[string]interface{}, err error) {
	var userList []*User
	err = dao.DB.Where("username=?", username).Last(&userList).Error
	if err != nil {
		fmt.Println("mysql err:", err)
		return nil, err
	}
	dict = map[string]interface{}{}
	for _, b := range userList {
		dict["name"] = b.Realname
		dict["sex"] = b.Gender
	}
	return dict, nil
}

//
func GetStuNumAName(level int) (StuNameAName map[string]string, err error) {
	var userList []*User
	if level == 0 {
		err = dao.DB.Find(&userList).Error
	} else {
		err = dao.DB.Where("grade=?", level).Find(&userList).Error
	}
	if err != nil {
		fmt.Println("GetStuNumAName failed,err:", err)
		return nil, err
	}
	SNName := map[string]string{}
	for _, b := range userList {
		SNName[b.Username] = b.Realname
	}
	return SNName, nil
}

func GetEmail(usernameSli []interface{}) (emailSli []string, err error) {
	var userList []*User
	err = dao.DB.Find(&userList).Error
	if err != nil {
		fmt.Println("mysql err:", err)
		return nil, err
	}
	emailSli = []string{}
	for _, stuNum := range usernameSli {
		for _, UserInfo := range userList {
			if stuNum == UserInfo.Username {
				emailSli = append(emailSli, UserInfo.Qq)
			}
		}
	}
	return emailSli, nil
}

func GetApplyLists() (AL []interface{}, err error) {
	var userList []*User
	err = dao.DB.Find(&userList).Error
	if err != nil {
		fmt.Println("GetApplyLists failed,err:", err)
		return nil, err
	}
	temp_intf := []interface{}{}
	for _, b := range userList {
		dict := map[string]interface{}{}
		dict["email"] = b.Qq
		dict["schoolNumber"] = b.Username
		if b.State == "0" {
			temp_intf = append(temp_intf, dict)
		}
	}
	return temp_intf, nil
}

func ChangeStatus(username, belong string) (err error) {
	user := User{}
	err1 := dao.DB.Where("username=?", username).Last(&user).Error
	if err1 != nil {
		fmt.Println("ChangeStatus failed,err:", err1)
		return err1
	}
	user.State = "1"
	user.Belong = belong
	dao.DB.Save(&user)
	return
}

func ResetPasswordORDelete(username, behave string) (err error) {
	user := User{}
	err1 := dao.DB.Where("username=?", username).Last(&user).Error
	if err1 != nil {
		return err1
	}
	if behave == "repass" {
		user.Password = "123456"
	}
	if behave == "delete" {
		user.State = "0"
	}
	dao.DB.Save(&user)
	return nil
}

//拿取map[stuNum]name
func GetMapStuNumName(level int) (userM map[string]string, err error) {
	var userList []*User
	userMap := map[string]string{}
	if level == 0 {
		err = dao.DB.Find(&userList).Error
	} else {
		err = dao.DB.Where("grade=?", level).Find(&userList).Error
	}
	if err != nil {
		return nil, err
	}
	for _, item := range userList {
		userMap[item.Username] = item.Realname
	}
	//fmt.Println(userMap)
	return userMap, nil
}

func UserLogin(username, password string) (statusMsg string) {
	userInfo := User{}
	err := dao.DB.Where("username=?", username).Last(&userInfo).Error
	if err != nil {
		return "noRecord"
	}
	//fmt.Println(userInfo)
	if userInfo.State == "0" {
		return "noAgree"
	}
	if userInfo.Password == password {
		return "success"
	} else {
		return "passwordErr"
	}
}

func UserRegister(username, password, email, grade, realname string) {
	userInfo := User{}
	userInfo.Username = username
	userInfo.Password = password
	userInfo.Qq = email
	//userInfo.State = "0"
	userInfo.State = "1"
	userInfo.Belong = "519"

	userInfo.Realname = realname
	userInfo.Role = "user"
	gradeInt, _ := strconv.Atoi(grade)
	userInfo.Grade = gradeInt
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	userInfo.Createtime = nowTime
	dao.DB.Save(&userInfo)
}

func ResetPassword(email, newPassword, username string) {
	userInfo := User{}
	dao.DB.Where("qq=?", email).Where("username=?", username).Last(&userInfo)
	userInfo.Password = newPassword
	dao.DB.Save(&userInfo)
}

func QueryGrade(username string) (grade int) {
	userInfo := User{}
	dao.DB.Where("username=?", username).First(&userInfo)
	return userInfo.Grade
}

func QueryPersonalInfo(username string) (userInfoMap map[string]interface{}) {
	user := User{}
	userInfo := user.getPersonalInfo(username)
	uim := map[string]interface{}{}
	uim["realname"] = userInfo.Realname
	uim["gender"] = userInfo.Gender
	uim["profession"] = userInfo.Profession
	uim["email"] = userInfo.Qq
	uim["nickname"] = userInfo.Nickname
	uim["grade"] = userInfo.Grade
	uim["username"] = userInfo.Username
	uim["role"] = userInfo.Role
	if userInfo.Head_pic == "" {
		uim["head_pic"] = "null"
	} else {
		uim["head_pic"] = userInfo.Head_pic
	}
	return uim
}

func ModifyPersonalInfo(username, realname, gender, profession, email, nickname, grade string) error {
	user := User{}
	userInfo := user.getPersonalInfo(username)
	userInfo.Realname = realname
	userInfo.Gender = gender
	userInfo.Profession = profession
	userInfo.Qq = email
	userInfo.Nickname = nickname
	gradeInt, _ := strconv.Atoi(grade)
	userInfo.Grade = gradeInt
	err := dao.DB.Save(&userInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func ModifyHeadPic(username, headPic string) error {
	user := User{}
	userInfo := user.getPersonalInfo(username)
	userInfo.Head_pic = headPic
	err := dao.DB.Save(&userInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func CheckRole(username string) (role string) {
	user := User{}
	userInfo := user.getPersonalInfo(username)
	return userInfo.Role
}

func QueryHeadPic(username string) (headPic string) {
	user := User{}
	userInfo := user.getPersonalInfo(username)
	return userInfo.Head_pic
}

func UpUserRole(username string) {
	user := User{}
	userInfo := user.getPersonalInfo(username)
	userInfo.Role = "admin"
	dao.DB.Save(&userInfo)
}

func GetContactListIn() interface{} {
	type TightUser struct {
		Userid   int    `json:"userid"`
		Username string `json:"username"`
		Realname string `json:"realname"`
		Grade    int    `json:"grade"`
	}
	TightUsers := []TightUser{}
	sql2 := "select userid,username,realname,grade from users where state = 1"
	dao.DB.Raw(sql2).Scan(&TightUsers)
	var gradeSLi []int
	for _, user := range TightUsers {
		gradeSLi = append(gradeSLi, user.Grade)
	}
	newGradeList := util.RemoveRepeatedElement(gradeSLi)
	contactList := map[int]interface{}{}
	for _, grade := range newGradeList {
		var userSLi []interface{}
		for _, user := range TightUsers {
			if user.Grade == grade {
				userSLi = append(userSLi, user)
			}
		}
		contactList[grade] = userSLi
	}
	return contactList
}
