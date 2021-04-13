package logger

import "log"

// LogError 检查数据库是否报错，如果有报错会将信息写进log.
func LogError(err error) {
	if err != nil {
		log.Println(err)
	}
}
