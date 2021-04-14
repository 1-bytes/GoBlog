package main

import (
	"GoBlog/bootstrap"
	"GoBlog/pkg/database"
	"GoBlog/pkg/logger"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

var router *mux.Router
var db *sql.DB

// Article 对应一条文章数据.
type Article struct {
	Title, Body string
	ID          int64
}

// Delete 方法用于从数据库中删除单条记录
func (a Article) Delete() (rowsAffected int64, err error) {
	rs, err := db.Exec("DELETE FROM articles WHERE id = ?", a.ID)
	if err != nil {
		return 0, err
	}
	// 删除成功，跳转到文章详情页
	if affected, _ := rs.RowsAffected(); affected > 0 {
		return affected, nil
	}
	return 0, nil
}

// getArticleById 通过传参ID获取博文.
func getArticleById(id string) (Article, error) {
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err
}

// articlesDeleteHandler 删除文章
func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := getRouteVariable("id", r)
	// 2. 读取对应的文章数据
	article, err := getArticleById(id)
	// 3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 服务器内部错误")
		}
	} else {
		// 4. 未出现错误，执行删除操作
		rowsAffected, err := article.Delete()
		// 4.1 发生错误
		if err != nil {
			// SQL报错
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 服务器内部错误")
		} else {
			// 4.2 未发生错误
			if rowsAffected > 0 {
				// 重定向到文章列表页
				indexURL, _ := router.Get("articles.index").URL()
				http.Redirect(w, r, indexURL.String(), http.StatusFound)
			} else {
				// Edge case
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "404 文章未找到")
			}
		}
	}
}

// forceHTMLMiddleware 中间件,用于设置返回的Header中的ContentType.
func forceHTMLMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置Header
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		// 继续处理请求
		h.ServeHTTP(w, r)
	})
}

// removeTrailingSlash 中间件,路由清理末尾的斜杠.
func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}

// GetRouteVariable 通过获取URL路由参数名称获取值.
func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

// main ...
func main() {
	database.Initialize()
	db = database.DB

	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()

	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")

	// 中间件：强制内容为 HTML
	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
