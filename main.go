package main

import (
	"drop_os_back/dao"
	"drop_os_back/models"
	"drop_os_back/routers"
	"drop_os_back/timeTask"
	"github.com/robfig/cron"
	"log"
)

var MysqlModels = []interface{}{&models.Administrator{}, &models.Time{}, &models.User{}, &models.Rank{}}

//定义gin-定时任务
func initTimer() {
	log.Println("Timer Starting...")
	c := cron.New()
	//将weektime表的时间清零
	c.AddFunc("0 0 0 ? * MON", timeTask.ResetWeekTime)
	//log.Println("前台每周时间清零，重新计算周排名，重新获取每周周一的date")
	//每周一凌晨重置周一的日期
	c.AddFunc("0 0 0 ? * MON", timeTask.ResetMonday)
	c.Start()
}


func main() {
	//设置程序内本周周一日期
	timeTask.SetMondayDateStartProcess()
	//创建数据库
	// TODO 修改部分框架语句，改用sql原生语句，优化mysql吞吐和性能开销
	// TODO 定义定时任务（每日针对学员发送邮件：您上周的实验室学习时间为_小说，请继续加油）
	// TODO 定义中间件（每周排名：每周时间清零）
	// TODO 针对总学习时间单独建库，优化运行时间
	// TODO 向各个成员发送考勤情况的邮件
	//开启定时任务
	initTimer()

	//开启监听 双线程监听
	for i := 0; i < 5; i++ {
		//websocket监听及处理线程
		go models.Manager.Start()
		//开启接收消费者动作的处理
		go models.ReceiveToInsert()
		//开启消息推送的管理线程
		go models.ReceiveMsgManager.RecMsgStart()
	}

	//连接数据库
	err := dao.InitMySQL()
	if err != nil {
		panic(err)
	}
	defer dao.Close() //程序退出关闭数据库
	//初始化redis
	err = dao.InitRedis()
	if err != nil {
		panic(err)
	}
	//模型绑定
	dao.DB.AutoMigrate(MysqlModels...)
	r := routers.SetupRouter()

	r.Run("0.0.0.0:9090")
	//r.RunTLS(":9090", "./ssl/server.crt", "./ssl/server.key")
}
