package mail

import (
	"GoBlog/pkg/config"
	"GoBlog/pkg/logger"
	email "github.com/xhit/go-simple-mail/v2"
	"time"
)

type (
	Mail struct{}
	// MapData 邮件模板格式预定义
	MapData       map[string]mailTemplates
	mailTemplates struct {
		from  string
		title string
		body  string
	}
)

var smtpClient *email.SMTPClient

func (m *Mail) getClient() {
	var err error
	server := &email.SMTPServer{
		Host:           config.GetString("mail.host"),
		Port:           config.GetInt("mail.port"),
		Username:       config.GetString("mail.username"),
		Password:       config.GetString("mail.password"),
		Encryption:     email.Encryption(config.GetInt("mail.encryption")),
		ConnectTimeout: time.Duration(config.GetInt64("mail.connect_timeout")) * time.Second,
		SendTimeout:    time.Duration(config.GetInt64("mail.send_timeout")) * time.Second,
		KeepAlive:      config.GetBool("mail.keep_alive"),
	}
	smtpClient, err = server.Connect()
	logger.LogError(err)
}

// Send 发送邮件
func (m *Mail) Send(tplName string, toEmail string, arg string) error {
	tpl := m.getTemplate(tplName, arg)
	mailMSG := email.NewMSG()
	mailMSG.SetFrom(tpl.from).
		AddTo(toEmail).
		SetSubject(tpl.title)
	mailMSG.SetBody(email.TextPlain, tpl.body)
	err := mailMSG.Send(smtpClient)
	logger.LogError(err)
	return err
}

// getTemplate 获取邮箱模板
func (m *Mail) getTemplate(tplName string, arg string) mailTemplates {
	emailUsername := config.GetString("mail.username")
	data := MapData{
		"verifyMail": mailTemplates{
			from:  "GoBlog 技术支持 <" + emailUsername + ">",
			title: "邮箱验证",
			body:  "你好，这是一封邮箱验证邮件（如果不是您的操作，请忽略。）\n验证码：" + arg,
		},
		"lostPassword": mailTemplates{
			from:  "GoBlog 技术支持 <" + emailUsername + ">",
			title: "找回密码",
			body:  "你好，这是一封找回密码邮件（如果不是您的操作，请忽略。）\n重置密码链接：" + arg,
		},
	}
	return data[tplName]
}
