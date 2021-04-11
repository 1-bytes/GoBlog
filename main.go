package main

import (
	"GoBlog/pkg/route"
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
	"unicode/utf8"
)

var router *mux.Router
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

// Article 对应一条文章数据.
type Article struct {
	Title, Body string
	ID          int64
}

// Int64ToString 将 int64 转换为 string
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
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

// Link 方法用来生成文章链接.
func (a Article) Link() string {
	showURL, err := router.Get("articles.show").URL("id", strconv.FormatInt(a.ID, 10))
	if err != nil {
		checkError(err)
		return ""
	}
	return showURL.String()
}

// articlesShowHandler 文章详情.
func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	// 1.获取 URL 参数
	id := getRouterVariable("id", r)

	// 2.读取对应的文章数据
	article, err := getArticleById(id)

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
		// 4.读取数据成功，显示表单
		tmpl, err := template.New("show.gohtml").Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
			"Int64ToString": Int64ToString,
		}).ParseFiles("resources/views/articles/show.gohtml")
		checkError(err)
		tmpl.Execute(w, article)
	}
}

// getRouterVariable 通过获取URL路由参数名称获取值.
func getRouterVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

// getArticleById 通过传参ID获取博文.
func getArticleById(id string) (Article, error) {
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err
}

// articlesEditHandler 更新文章页面
func articlesEditHandler(w http.ResponseWriter, r *http.Request) {
	// 1.获取 URL 参数
	id := getRouterVariable("id", r)

	// 2.读取对应的文章数据
	article, err := getArticleById(id)

	// 3.如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "文章未找到")
		} else {
			// 3.2 数据库错误
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 读取成功，显示表单
		updateURL, err := router.Get("articles.update").URL("id", id)
		data := ArticlesFormData{
			Title:  article.Title,
			Body:   article.Body,
			URL:    updateURL,
			Errors: nil,
		}

		teml, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		checkError(err)

		teml.Execute(w, data)
	}
}

// articlesUpdateHandler 更新文章接口
func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// 1.获取URL参数
	id := getRouterVariable("id", r)

	// 2.读取对应的文章数据
	_, err := getArticleById(id)

	// 3.如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4.未出现错误，进行表单验证
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")

		errors := validateArticleFormData(title, body)

		if len(errors) == 0 {
			// 表单验证通过，更新数据
			query := "UPDATE articles SET title = ?, body = ? WHERE id = ?"
			rs, err := db.Exec(query, title, body, id)

			if err != nil {
				checkError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
			}

			// 更新成功，跳转到文章详情页
			if n, _ := rs.RowsAffected(); n > 0 {
				showUrl, _ := router.Get("articles.show").URL("id", id)
				http.Redirect(w, r, showUrl.String(), http.StatusFound)
			} else {
				fmt.Fprint(w, "您没有做任何更改！")
			}
		} else {
			// 4.3表单验证不通过，验证路由
			updateURL, _ := router.Get("articles.update").URL("id", id)
			data := ArticlesFormData{
				Title:  title,
				Body:   body,
				URL:    updateURL,
				Errors: errors,
			}
			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			checkError(err)
			tmpl.Execute(w, data)
		}

	}
}

// articlesDeleteHandler 删除文章
func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := getRouterVariable("id", r)
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
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 服务器内部错误")
		}
	} else {
		// 4. 未出现错误，执行删除操作
		rowsAffected, err := article.Delete()
		// 4.1 发生错误
		if err != nil {
			// SQL报错
			checkError(err)
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

// articlesIndexHandler 文章列表.
func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	// 1.执行查询语句，返回一个结果集
	rows, err := db.Query("SELECT * FROM articles")
	checkError(err)
	defer rows.Close()

	// 2.循环读取结果
	var articles []Article
	for rows.Next() {
		var article Article
		// 2.1扫描每一行的结果并赋值到一个 article 对象中
		err := rows.Scan(&article.ID, &article.Title, &article.Body)
		checkError(err)
		// 2.2将 article 追加到 articles 这个切片当中
		articles = append(articles, article)
	}

	// 2.3检测遍历时是否发生错误
	err = rows.Err()
	checkError(err)

	// 3.加载模板
	tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
	checkError(err)

	// 4.渲染模板，将所有文章的数据传输进去
	tmpl.Execute(w, articles)
}

// AtriclesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

func validateArticleFormData(title string, body string) map[string]string {
	errors := make(map[string]string)

	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString("title") > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度不能大于或等于10个字节"
	}
	return errors
}

// articlesStoreHandler 创建新的文章 API接口.
func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	body := r.PostFormValue("body")

	errors := validateArticleFormData(title, body)

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
	route.Initialize()
	router = route.Router

	router.HandleFunc("/", homeHandler).Name("home")
	router.HandleFunc("/about", aboutHandler).Name("about")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")
	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")
	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")

	// 自定义404页面
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 中间件：强制内容为 HTML
	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
