package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var router = mux.NewRouter()
var db *sql.DB

// homeHandler 首页.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello，这里是blog</h1>")
}

// aboutHandler 关于.
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:abc@example.com\">abc@example.com</a>")
}

// notFoundHandler 404.
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

// Article 对应一条文章数据
type Article struct {
	Title, Body string
	ID          int64
}

// articlesShowHandler 文章详情.
func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	// 1.获取 URL 参数
	vars := mux.Vars(r)
	id := vars["id"]

	// 2.读取对应的文章数据
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)

	// 3.如果出现错误
	if err != nil {
		// 判断是没找到数据 还是查询报错了
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "文章未找到")
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4.读取数据成功
		tmpl, err := template.ParseFiles("resources/views/articles/show.gohtml")
		checkError(err)
		tmpl.Execute(w, article)
	}
}

// articlesIndexHandler 文章列表 API接口.
func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "访问文章列表")
}

// AtriclesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

// articlesStoreHandler 创建新的文章 API接口.
func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	body := r.PostFormValue("body")

	errors := make(map[string]string)

	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if len(title) < 3 || len(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if len(body) < 10 {
		errors["body"] = "内容长度需大于或等于10个字节"
	}

	// 检查是否有错误
	if len(errors) == 0 {
		lastInsertID, err := saveArticleToDB(title, body)
		if lastInsertID > 0 {
			fmt.Fprint(w, "插入成功，ID为"+strconv.FormatInt(lastInsertID, 10))
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		storeURL, _ := router.Get("articles.store").URL()
		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}
		tmpl.Execute(w, data)
	}
}

// saveArticleToDB 向数据库中保存一篇文章，如果插入成功则返回一个主键ID.
func saveArticleToDB(title string, body string) (int64, error) {
	// 变量初始化
	var (
		id     int64
		err    error
		stmt   *sql.Stmt
		result sql.Result
	)

	// 1.获取一个 Prepare 声明语句
	stmt, err = db.Prepare("INSERT INTO articles(title, body) VALUES (?, ?)")
	if err != nil {
		checkError(err)
		return 0, err
	}

	// 2.在此函数运行结束后关闭此语句，防止占用SQL连接
	defer stmt.Close()

	// 3.执行请求，传参进入绑定的内容
	result, err = stmt.Exec(title, body)
	if err != nil {
		checkError(err)
		return 0, err
	}

	// 4.插入成功的话，会返回自增的ID
	if id, err = result.LastInsertId(); id > 0 {
		return id, nil
	}
	return 0, err
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

// articlesCreateHandler 新建文章页面.
func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	storeURL, _ := router.Get("articles.store").URL()
	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}
	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, data)
}

// initDB 初始化数据库.
func initDB() {
	var err error
	config := mysql.Config{
		User:                 "root",
		Passwd:               "FLHY1VJ0WEOBIG3N",
		Addr:                 "127.0.0.1:3306",
		Net:                  "tcp",
		DBName:               "goblog",
		AllowNativePasswords: true,
	}

	// 设置数据库连接池
	db, err = sql.Open("mysql", config.FormatDSN())
	checkError(err)

	// 设置最大连接数
	db.SetMaxOpenConns(25)

	// 设置最大空闲连接数
	db.SetMaxIdleConns(25)

	// 设置每个链接的过期时间
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	checkError(err)
}

// checkError 检查数据库是否报错，如果有报错会将信息写进log.
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// createTables 创建数据库表.
func createTables() {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
    id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    body longtext COLLATE utf8mb4_unicode_ci
); `

	_, err := db.Exec(createArticlesSQL)
	checkError(err)
}

func main() {
	initDB()
	createTables()
	router.HandleFunc("/", homeHandler).Name("home")
	router.HandleFunc("/about", aboutHandler).Name("about")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")

	// 自定义404页面
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 中间件：强制内容为 HTML
	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
