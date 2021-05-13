package main

import (
	"GoBlog/app/http/middlewares"
	"GoBlog/bootstrap"
	configs "GoBlog/config"
	"GoBlog/pkg/config"
	"net/http"
)

func init() {
	// 初始化配置信息
	configs.Initialize()
}

func main() {
	bootstrap.SetupDB()
	router := bootstrap.SetupRoute()

	http.ListenAndServe(":"+config.GetString("app.port"), middlewares.RemoveTrailingSlash(router))
}
