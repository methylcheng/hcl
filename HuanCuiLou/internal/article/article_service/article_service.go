package article_service

import (
	"fmt"
	"huancuilou/common"
	"huancuilou/internal/article/article_model"
	"huancuilou/internal/article/article_repository"
	"time"
)

//article的业务逻辑处理层

type ArticleService struct {
	articleRepository *article_repository.ArticleRepository
}

func NewArticleService(articleRepository *article_repository.ArticleRepository) *ArticleService {
	return &ArticleService{
		articleRepository: articleRepository,
	}
}

func (as ArticleService) AddArticle(article article_model.Article) (string, error) {
	//接入微信审查接口
	openId := fmt.Sprintf("%v", article.ManagerId)
	checkResult := common.CheckContentSecurity(article.Content, openId, 1, article.Title)

	//返回true执行插入逻辑
	if checkResult {
		article.CreatedAt = time.Now()

		result, err := as.articleRepository.AddArticle(article)
		if err != nil {
			return "审核通过但Repository层出错", err
		}
		if result {
			return "发布成功", nil
		}
	}
	return "微信审核未通过", nil
}

func (as ArticleService) PageGet(pageget article_model.PageGet) ([]article_model.Article, error) {
	articles, err := as.articleRepository.PageGet(pageget)

	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (as ArticleService) ArticleGet(article article_model.Article) (article_model.Article, error) {
	var err error
	article, err = as.articleRepository.ArticleGet(article)
	if err != nil {
		return article, err
	}

	return article, nil
}
