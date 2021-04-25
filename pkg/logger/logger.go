package logger

import "log"

// LogError 检查是否存在报错，如果有报错会将信息写进log
func LogError(err error) {
	if HasError(err) {
		log.Println(err)
	}
}

// HasError 检查是否有错误，存在错误返回 true，没有错误返回 false
func HasError(err error) bool {
	return err != nil
}
