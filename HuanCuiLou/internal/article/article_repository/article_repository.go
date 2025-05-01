package article_repository

//article_repository ：article的数据访问层

import (
	"fmt"
	"gorm.io/gorm"
	"huancuilou/internal/article/article_model"
)

type ArticleRepository struct {
	DB *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{
		DB: db,
	}
}

func (ar ArticleRepository) AddArticle(article article_model.Article) (bool, error) {
	result := ar.DB.Create(&article)

	if result.Error != nil {
		return false, fmt.Errorf("ArticleRepository.AddArticle err:%w", result.Error)
	}
	return true, nil
}

func (ar ArticleRepository) PageGet(pageget article_model.PageGet) ([]article_model.Article, error) {
	var articles []article_model.Article
	fmt.Println(pageget.PageSize, pageget.PageNum, pageget.Kind)
	result := ar.DB.Where("kind = ?", pageget.Kind).
		Limit(pageget.PageSize).
		Offset((pageget.PageNum - 1) * pageget.PageSize).
		Find(&articles)
	if result.Error != nil {
		return nil, fmt.Errorf("ArticleRepository.PageGet err:%w", result.Error)
	}
	return articles, nil
}

func (ar ArticleRepository) ArticleGet(article article_model.Article) (article_model.Article, error) {
	result := ar.DB.Where("title = ? AND kind = ?", article.Title, article.Kind).
		First(&article)

	if result.Error != nil {
		return article, fmt.Errorf("ArticleRepository.ArticleGet err:%w", result.Error)
	}
	return article, nil
}
