package article

import (
	"GoBlog/app/models"
	"GoBlog/app/models/user"
	"GoBlog/pkg/route"
)

// Article 文章模型
type Article struct {
	models.BaseModel
	Title string `gorm:"type:varchar(255);not null" valid:"title"`
	Body  string `gorm:"type:longtext" valid:"body"`
	// Article 属于 User，UserID 是外键
	UserID uint64 `gorm:"not null;index"`
	// 默认情况下 UserID 会隐式的作用于 Article 表和 User 表之间创建外键关系
	// 因此必须包含在 Article 结构中，以便填充 User 的内部结构
	User user.User
}

// Link 方法用来生成文章链接
func (article Article) Link() string {
	return route.Name2URL("articles.show", "id", article.GetStringID())
}

// CreateAtDate 创建日期
func (article Article) CreateAtDate() string {
	return article.CreateAt.Format("2006-01-02")
}
