package policies

import (
	"GoBlog/app/models/article"
	"GoBlog/pkg/auth"
)

// CanModifyArticle 是否允许修改文章
func CanModifyArticle(_article article.Article) bool {
	return auth.User().ID == _article.UserID
}
