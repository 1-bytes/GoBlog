package article

import (
	"GoBlog/app/models"
	"GoBlog/pkg/route"
)

// Article 文章模型
type Article struct {
	models.BaseModel
	Title string `gorm:"type:varchar(255);not null" valid:"title"`
	Body  string `gorm:"type:longtext" valid:"body"`
}

// Link 方法用来生成文章链接
func (article Article) Link() string {
	return route.Name2URL("articles.show", "id", article.GetStringID())
}
