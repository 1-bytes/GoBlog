package route

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Router 路由对象.
var Router *mux.Router

// Initialize 初始化路由.
func Initialize() {
	Router = mux.NewRouter()
}

// Name2URL 通过路由名称来获取URL.
func Name2URL(routeName string, pairs ...string) string {
	url, err := Router.Get(routeName).URL(pairs...)
	if err != nil {
		//checkError(err)
		return ""
	}
	return url.String()
}

// GetRouteVariable 通过获取URL路由参数名称获取值.
func GetRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}
