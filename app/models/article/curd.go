package article

import (
	"GoBlog/pkg/model"
	"GoBlog/pkg/types"
)

// Get 通过 ID 获取文章
func Get(idstr string) (Article, error) {
	var article Article
	id := types.StringToInt(idstr)
	if err := model.DB.First(&article, id).Error; err != nil {
		return article, err
	}
	return article, nil
}

// GetAll 获取全部文章
func GetAll() ([]Article, error) {
	var articles []Article
	if err := model.DB.Find(&articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}
