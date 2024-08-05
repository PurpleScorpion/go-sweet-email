package sweetEmail

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/PurpleScorpion/go-sweet-email/logger"
	"net/smtp"
)

var conf EmailConf
var logFlag bool = false
var regFlag bool = false

func SetPort(port int) {
	if regFlag {
		return
	}
	if port < 0 || port > 65535 {
		panic("port out of range ( 0-65535 )")
	}
	conf.Port = port
}

func SetHost(host string) {
	if regFlag {
		return
	}
	if isEmpty(host) {
		panic("host cannot be empty")
	}
	conf.Host = host
}

func SetUserName(userName string) {
	if regFlag {
		return
	}
	if isEmpty(userName) {
		panic("userName cannot be empty")
	}
	conf.UserName = userName
}

func SetPassword(password string) {
	if regFlag {
		return
	}
	if isEmpty(password) {
		panic("password cannot be empty")
	}
	conf.Password = password
}

func SetEmailName(emailName string) {
	conf.EmailName = emailName
}

func OpenLog() {
	logFlag = true
}

/*
自动注册邮件配置

	参数:
		flag: 是否开启日志
	要求:
		1. 需使用go-sweet框架编写的项目
		github地址: https://github.com/PurpleScorpion/go-sweet
		2. yaml格式必须为以下写法
			sweet:
			  email:
				port: 25
				host: smtp.qq.com
				username: 123456@qq.com
				password: xtjqjqjqjqjqjqjq
				emailname: Sweet
*/
func AutoRegister(flag bool) {
	conf.Port = valueInt("${sweet.email.port}")
	conf.Host = valueString("${sweet.email.host}")
	conf.UserName = valueString("${sweet.email.username}")
	conf.Password = valueString("${sweet.email.password}")
	conf.EmailName = valueString("${sweet.email.emailname}")
	logFlag = flag
	logger.Info("port: %d, host: %s, username: %s, password: %s, emailname: %s", conf.Port, conf.Host, conf.UserName, conf.Password, conf.EmailName)
	Register()
}

func Register() {
	if conf.Port < 0 || conf.Port > 65535 {
		panic("port out of range ( 0-65535 )")
	}
	if isEmpty(conf.Host) {
		panic("host cannot be empty")
	}
	if isEmpty(conf.UserName) {
		panic("userName cannot be empty")
	}
	if isEmpty(conf.Password) {
		panic("password cannot be empty")
	}
	if isEmpty(conf.EmailName) {
		conf.EmailName = conf.UserName
	}
	regFlag = true
}

/*
to: []string 收信者邮箱列表
subject: string 邮件主题
body: string 邮件内容-HTML格式
*/
func SendEmail(to []string, subject string, body string) error {
	if !regFlag {
		return errors.New("please register email first")
	}
	if to == nil {
		return errors.New("to users cannot be nil")
	}
	if len(to) == 0 {
		return errors.New("to users cannot be empty")
	}
	if isEmpty(subject) {
		return errors.New("subject cannot be empty")
	}
	if isEmpty(body) {
		return errors.New("body cannot be empty")
	}
	host := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	// 设置发件人邮箱地址、SMTP服务器地址、端口和授权码
	auth := smtp.PlainAuth("", conf.UserName, conf.Password, host)
	from := fmt.Sprintf("%s <%s>", conf.EmailName, conf.UserName)
	for i := 0; i < len(to); i++ {
		if logFlag {
			logger.Info("[sweet-email info] send email to %s", to[i])
		}
		flag := send(auth, from, to[i], subject, body)
		if logFlag && flag {
			logger.Info("[sweet-email info] send email to %s success", to[i])
		}
	}
	return nil
}

func send(auth smtp.Auth, from string, to string, subject string, body string) bool {
	// 创建邮件内容
	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		body + "\r\n")

	if logFlag {
		logger.Info("[sweet-email info] send email msg : %s", string(msg))
	}

	// SMTP连接配置，QQ邮箱使用SSL加密连接

	// 创建TLS连接
	conn, err := tls.Dial("tcp", conf.Host, nil)
	if err != nil {
		logger.Error("[sweet-email error] [%s] dial tls failed: %v", to, err)
		return false
	}
	defer conn.Close()

	// 构建SMTP客户端
	client, err := smtp.NewClient(conn, conf.Host)
	if err != nil {
		logger.Error("[sweet-email error] [%s] new client failed: %v", to, err)
		return false
	}
	defer client.Quit()

	// 开启身份验证
	if err = client.Auth(auth); err != nil {
		logger.Error("[sweet-email error] [%s] auth failed: %v", to, err)
		return false
	}

	// 设置发送邮件选项（比如：发件人）
	if err = client.Mail(conf.UserName); err != nil {
		logger.Error("[sweet-email error] [%s] mail failed: %v", to, err)
		return false
	}

	// 设置接收邮件选项（比如：收件人）
	if err = client.Rcpt(to); err != nil {
		logger.Error("[sweet-email error] [%s] rcpt failed: %v", to, err)
		return false
	}

	// 写入邮件内容
	w, err := client.Data()
	if err != nil {
		logger.Error("[sweet-email error] [%s] data failed: %v", to, err)
		return false
	}
	_, err = w.Write(msg)
	if err != nil {
		logger.Error("[sweet-email error] [%s] write failed: %v", to, err)
		return false
	}
	err = w.Close()
	if err != nil {
		logger.Error("[sweet-email error] [%s] close failed: %v", to, err)
		return false
	}
	return true
}
