package controllers

import (
	"GoBlog/app/models/article"
	"GoBlog/app/models/category"
	"GoBlog/app/requests"
	"GoBlog/pkg/config"
	"GoBlog/pkg/route"
	"GoBlog/pkg/view"
	"fmt"
	"net/http"
)

// CategoriesController 文章分类控制器
type CategoriesController struct {
	BaseController
}

// Create 文章分类创建页面
func (cc *CategoriesController) Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, "categories.create")
}

// Store 保存文章分类
func (cc *CategoriesController) Store(w http.ResponseWriter, r *http.Request) {
	// 1. 初始化数据
	_category := category.Category{
		Name: r.PostFormValue("name"),
	}
	// 2. 表单验证
	errs := requests.ValidateCategoryForm(_category)
	// 3. 检测错误
	if len(errs) == 0 {
		// 创建文章分类
		_category.Create()
		if _category.ID > 0 {
			fmt.Fprint(w, "创建成功！")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建文章分类失败，请联系管理员")
		}
	} else {
		view.Render(w, view.D{
			"Category": _category,
			"Errors":   errs,
		}, "categories.create")
	}
}

// Show 显示分类下的文章列表
func (cc *CategoriesController) Show(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的数据
	_category, err := category.Get(id)

	// 3. 获取结果集
	articles, pagerData, err := article.GetByCategoryID(_category.GetStringID(), r, config.GetInt("pagination.perpage"))
	if err != nil {
		cc.ResponseForSQLError(w, err)
	} else {
		// 加载模板
		view.Render(w, view.D{
			"Articles":  articles,
			"PagerData": pagerData,
		}, "articles.index", "articles._article_meta")
	}
}
