package logger

import "log"

// LogError 检查是否存在报错，如果有报错会将信息写进log
func LogError(err error) {
	if err != nil {
		log.Println(err)
	}
}
