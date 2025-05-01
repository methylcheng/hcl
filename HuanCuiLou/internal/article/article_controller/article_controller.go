package article_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"huancuilou/internal/article/article_model"
	"huancuilou/internal/article/article_service"
	"huancuilou/response"
	"net/http"
)

//处理有关article的请求与返回响应

type ArticleController struct {
	articleService *article_service.ArticleService
}

func NewArticleController(articleService *article_service.ArticleService) *ArticleController {
	return &ArticleController{
		articleService: articleService,
	}
}

// AddArticle 处理添加文章的请求和响应
func (ac ArticleController) AddArticle(c *gin.Context) {
	var article article_model.Article

	if err := c.BindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message, err := ac.articleService.AddArticle(article)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Success(message))
}

// PageGet 分页查询,num表示每页显示的数据个数，page表示页数
func (ac ArticleController) PageGet(c *gin.Context) {
	var pageget article_model.PageGet
	if err := c.BindJSON(&pageget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	fmt.Println(pageget)
	articles, err := ac.articleService.PageGet(pageget)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Success(articles))
}

func (ac ArticleController) ArticleGet(c *gin.Context) {
	var article article_model.Article

	if err := c.BindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var err error
	article, err = ac.articleService.ArticleGet(article)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, response.Success(article))
}
