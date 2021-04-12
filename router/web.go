package router

import (
	"GoBlog/app/http/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterWebRoutes(r *mux.Router) {

	// 静态页面
	pc := new(controllers.PagesController)
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)
	r.HandleFunc("/", pc.Home).Name("home")
	r.HandleFunc("/about", pc.About).Name("about")
}
