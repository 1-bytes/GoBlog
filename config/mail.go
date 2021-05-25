package config

import (
	"GoBlog/pkg/config"
)

func init() {
	config.Add("mail", config.StrMap{
		// 邮箱 SMTP 服务域名的地址
		"host": config.Env("MAIL_HOST", "mail.gandi.net"),
		// 邮箱 SMTP 服务端口
		"port": config.Env("MAIL_PORT", 587),
		// 邮箱账号
		"username": config.Env("MAIL_USERNAME"),
		// 邮箱密码
		"password": config.Env("MAIL_PASSWORD"),
		/*
			邮箱连接加密方式，可直接传入下列对应数字来选择不同的加密
			0 EncryptionNone
			1 EncryptionSSL
			2 EncryptionTLS
			3 EncryptionSSLTLS
			4 EncryptionSTARTTLS
		*/
		"encryption": config.Env("MAIL_ENCRYPTION", 4),
		// 连接超时时间（单位：秒）
		"connect_timeout": config.Env("MAIL_CONNECT_TIMEOUT", 10),
		// 发送超时时间（单位：秒）
		"send_timeout": config.Env("MAIL_SEND_TIMEOUT", 10),
		// 是否开启长连接
		"keep_alive": config.Env("MAIL_KEEP_ALIVE", false),
	})
}
