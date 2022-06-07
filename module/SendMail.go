package module

import (
	"gopkg.in/gomail.v2"
)

func SendMailSli(recieverList []string) (err error) {
	m := gomail.NewMessage()
	for _, b := range recieverList {
		//发送人
		m.SetHeader("From", "2633122565@qq.com")
		//接收人
		m.SetHeader("To", b)
		//抄送人
		//m.SetAddressHeader("Cc", "xxx@qq.com", "xiaozhujiao")
		//主题
		m.SetHeader("Subject", "来自樱岛麻衣的提示")
		//内容
		body := "<h1>您的签到时间存在问题，请联系所属管理员</h1>"
		m.SetBody("text/html", body)
		//拿到token，并进行连接,第4个参数是填授权码
		d := gomail.NewDialer("smtp.qq.com", 587, "2633122565@qq.com", "raaepljblyzydiaa")
		// 发送邮件
		err = d.DialAndSend(m)
	}
	return
}

func SendEmail(reciever,code string) (err error) {
	m := gomail.NewMessage()
	//发送人
	m.SetHeader("From", "2633122565@qq.com")
	//接收人
	m.SetHeader("To", reciever)
	//抄送人
	//m.SetAddressHeader("Cc", "xxx@qq.com", "xiaozhujiao")
	//主题
	m.SetHeader("Subject", "Code")
	//内容
	body := "<h3>您的验证码为："+code+"，该验证码在十分钟内有效，请注意时间哦！</h3>"
	m.SetBody("text/html", body)
	//拿到token，并进行连接,第4个参数是填授权码
	d := gomail.NewDialer("smtp.qq.com", 587, "2633122565@qq.com", "raaepljblyzydiaa")
	// 发送邮件
	err = d.DialAndSend(m)
	return err
}
