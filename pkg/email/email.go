package email

import (
	"GoBlog/pkg/logger"
	mail "github.com/xhit/go-simple-mail/v2"
	"time"
)

type (
	SMTPServer struct {
		*mail.SMTPServer
	}

	// MapData 邮件模板格式预定义
	MapData  map[string]emailTpl
	emailTpl struct {
		from  string
		title string
		body  string
	}
)

var smtpClient *mail.SMTPClient

// const is e-mail account config info.
const (
	emailHost       = "mail.gandi.net"
	emailPort       = 587
	emailUsername   = "support@tcp.so"
	emailPassword   = "3KS7B5O8PM0NUTE9"
	emailEncryption = mail.EncryptionSTARTTLS
)

// init 初始化 Email 客户端
func init() {
	server := &SMTPServer{&mail.SMTPServer{
		Host:           emailHost,
		Port:           emailPort,
		Username:       emailUsername,
		Password:       emailPassword,
		Encryption:     emailEncryption,
		ConnectTimeout: 10 * time.Second,
		SendTimeout:    10 * time.Second,
		KeepAlive:      false,
	}}
	client, err := server.Connect()
	logger.LogError(err)
	smtpClient = client
}

// SendEmail 发送邮件
func (c *SMTPServer) SendEmail(tplName string, toEmail string) bool {
	// 获取邮件模板
	tpl := c.getEmailTpls(tplName, "https://xxxx.com")

	email := mail.NewMSG()
	email.SetFrom(tpl.from).
		AddTo(toEmail).
		SetSubject(tpl.title)
	email.SetBody(mail.TextPlain, tpl.body)
	err := email.Send(smtpClient)
	logger.LogError(err)
	return logger.HasError(err)
}

// getEmailTpls 获取邮箱模板
func (c *SMTPServer) getEmailTpls(tplName string, urls ...string) emailTpl {
	data := MapData{
		"verifyEmail": emailTpl{
			from:  "GoBlog 技术支持 <" + emailUsername + ">",
			title: "邮箱验证",
			body:  "你好，这是一封邮箱验证邮件。\n如果不是您的操作，请忽略。\n验证链接：",
		},
		"retrievePassword": emailTpl{
			from:  "GoBlog 技术支持 <" + emailUsername + ">",
			title: "找回密码",
			body:  "你好，这是一封找回密码邮件。\n如果不是您的操作，请忽略。\n点击链接重置密码：",
		},
	}
	result := data[tplName]
	for _, url := range urls {
		result.body += url
	}
	return result
}
