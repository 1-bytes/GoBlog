package main

import (
	"GoBlog/app/http/middlewares"
	"GoBlog/bootstrap"
	"net/http"
)

func main() {
	bootstrap.SetupDB()
	router := bootstrap.SetupRoute()

	http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))
}
