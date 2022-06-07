package routers

import (
	"drop_os_back/controller"
	"drop_os_back/middleware"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"net/http"
)

//routers
func SetupRouter() *gin.Engine {
	r := gin.Default()
	//r.Use(cors.Default())
	r.Use(middleware.Cors(), middleware.LogerMiddleware()) //,middleware.TlsHandler()
	r.StaticFS("/uploadPic", http.Dir("./uploadPic"))
	r.StaticFS("/appDownload", http.Dir("./appDownload"))
	//告诉gin框架模板引用的静态文件去哪里找
	//r.Static("/receptionDist", "receptionDist")
	//告诉gin框架去哪里找模板文件
	//r.LoadHTMLGlob("templates/*")
	//r.GET("/reception", controller.IndexHandler)
	r.Use(static.Serve("/", static.LocalFile("dist", true)))

	// --------------后台--------------------
	r.POST("/dev/admin_login", controller.AdminLogin)
	r.POST("/dev/logout", controller.AdminLogout)
	r.POST("/dev/query_student_time", controller.QueryStudentTime)
	r.POST("/dev/query_student", controller.QueryStudent)
	r.POST("/dev/query_unusual", controller.QueryUnusual)
	r.POST("/dev/query_euser")
	r.POST("/dev/query_student_time_single", controller.QueryStuTimeSingle)
	r.POST("/dev/query_all_time", controller.QueryAllTime)
	r.POST("/dev/end_sign", controller.EndSign)
	r.POST("/dev/sign_supply", controller.SignSupply)
	r.POST("/dev/sent_unusual_email", controller.SendUnusualEmail)
	r.POST("/dev/reply_lab", controller.ReplyLab)
	r.POST("/dev/apply_ok", controller.ApplyOK)
	r.POST("/dev/show_administrators", controller.ShowAdministrators)
	r.POST("/dev/add_administrator", controller.AddAdministrator)
	r.POST("/dev/query_student_date_time", controller.QueryStuDateTime)
	r.POST("/dev/remake_pass_delete", controller.RemakeORDelete)
	r.POST("/dev/delete_sign_data", controller.DeleteSignData)
	r.POST("/dev/upIdentity", controller.UpIdentity)
	//r.POST("/admin_login", controller.AdminLogin)
	//r.POST("/logout", controller.AdminLogout)
	//r.POST("/query_student_time", controller.QueryStudentTime)
	//r.POST("/query_student", controller.QueryStudent)
	//r.POST("/query_unusual", controller.QueryUnusual)
	//r.POST("/query_euser")
	//r.POST("/query_student_time_single", controller.QueryStuTimeSingle)
	//r.POST("/query_all_time", controller.QueryAllTime)
	//r.POST("/end_sign", controller.EndSign)
	//r.POST("/sign_supply", controller.SignSupply)
	//r.POST("/sent_unusual_email", controller.SendUnusualEmail)
	//r.POST("/reply_lab", controller.ReplyLab)
	//r.POST("/apply_ok", controller.ApplyOK)
	//r.POST("/show_administrators", controller.ShowAdministrators)
	//r.POST("/add_administrator", controller.AddAdministrator)
	//r.POST("/query_student_date_time", controller.QueryStuDateTime)
	//r.POST("/remake_pass_delete", controller.RemakeORDelete)
	//r.POST("/delete_sign_data", controller.DeleteSignData)
	//r.POST("/upIdentity", controller.UpIdentity)

	// --------------前台--------------------
	r.POST("/login", controller.UserLogin)
	r.POST("/sendEmail", controller.SendEmail)
	r.POST("/register", controller.UserRegister)
	r.POST("/endSign", controller.UserEndSign)
	r.POST("/startSign", controller.UserStartSign)
	r.POST("/resetPassword", controller.ResetPassword) //
	r.POST("/commentPublishArticle", middleware.CheckToken, controller.UserPublishCommentArt)
	r.POST("/comment", middleware.CheckToken, controller.UserDiscuss)
	r.GET("/getUnreadNumber", controller.GetUnreadMessage)
	r.GET("/notifications", controller.GetNotificationInfo)
	r.GET("/timeInfo", controller.UserGetTimeInfo)
	r.GET("/getDiscussionArea/:id/:pageNum", controller.GetDiscussionArea)
	r.GET("/getAllArticles/:classification/:pageNum", controller.GetAllArticles)
	r.GET("/getOneArticle/:articleId", controller.GetArticle)
	r.GET("/getStatus", controller.GetStatus)
	r.POST("/GetStudyRec", middleware.CheckToken, controller.UserGetStudyRec)
	r.GET("/getUserInfo", controller.UserGetPersonalInfo)
	r.POST("/modifyPersonalInfo", middleware.CheckToken, controller.ModifyPersonalInfo)
	r.POST("/uploadHeadPic", middleware.CheckToken, controller.UserUploadHeadPic)
	r.GET("/getTodoList", controller.UserGetTodo)
	r.POST("/UserAdminDeleteTodo", middleware.CheckToken, controller.UserAdminDeleteTodo)
	r.POST("/UserAdminAddTodo", middleware.CheckToken, controller.UserAdminAddTodo)
	r.GET("/GetUserHeadPic/:username", controller.GetUserHeadPic)
	r.POST("/UserLikeCon", middleware.CheckToken, controller.UserLikeCon)
	r.GET("/GetLikeStatus/:id", middleware.CheckToken, controller.UserGetThumbsUp)
	r.GET("/GetMoreDis/:artId/:parentId/:pageNum", controller.GetMoreDis)
	r.GET("/GetHeadPicBT", controller.GetHeadPicBT)
	//r.GET("/GetLikeStatus/:id", middleware.CheckToken, controller.UserGetThumbsUp)
	r.GET("/ws", middleware.WsCheckToken, controller.Handler) //
	r.POST("/postPic", middleware.CheckToken, controller.PostPic)
	r.GET("/pushMsg", middleware.WsCheckToken, controller.ReceiveMsgWebsocket)
	r.GET("/getPersonalInfo/:username", controller.GetPersonalInfo)
	r.GET("/getPersonalAllArticles/:username/:pageNum", controller.GetPersonalAllArticles)
	r.POST("/checkChatStatus", middleware.CheckToken, controller.CheckChatStatus)
	r.GET("/getMsgList", controller.GetMsgList)
	r.GET("/getContactList", middleware.CheckToken, controller.GetContactList)
	r.POST("/resetTokenFromMsgList", middleware.CheckToken, controller.ResetTokenFromMsgList)
	r.POST("/deleteMsgUser", middleware.CheckToken, controller.DeleteOneMsgUser)
	r.POST("/deleteArticle", middleware.CheckToken, controller.DeleteArticle)
	r.POST("/postFile", middleware.CheckToken, controller.PostFile)
	// --------------测试--------------------
	r.POST("/test", controller.TestAPI)
	r.POST("/test2", controller.TestReceiveFile)
	r.POST("/test3", controller.TestHeaders)

	return r

}
