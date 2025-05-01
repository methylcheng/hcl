package routers

import (
	"huancuilou/common"
	"huancuilou/internal/article/article_controller"
	"huancuilou/internal/homePage/homePage_controller"
	"huancuilou/internal/user/user_controller"

	"github.com/gin-gonic/gin"
)

// SetupRouters 设置所有路由
func SetupRouters(
	userController *user_controller.UserController,
	articleController *article_controller.ArticleController,
	homePageController *homePage_controller.HomePageController,
	sharedResourcesController *article_controller.SharedResourcesController,
) *gin.Engine {
	r := gin.Default()

	// 用户相关路由
	userGroup := r.Group("/user")
	{
		userGroup.GET("/send-code/:phoneNumber", userController.SendCode)
		userGroup.POST("/login", userController.Login)
		userGroup.PUT("/add-administrator/:phoneNumber", common.SuperAdminOnlyMiddleware(), userController.AddAdministrator)
		userGroup.GET("", common.JwtInterceptor(), userController.GetUserInfo)
		userGroup.PUT("", common.JwtInterceptor(), userController.UpdateUserInfo)
		userGroup.POST("/add-score", common.JwtInterceptor(), userController.AddScore)
		userGroup.POST("/add-phone-record", common.AdminOnlyMiddleware(), userController.AddPhoneRecord)
		userGroup.GET("/get-phone-record/:phoneNumber", common.AdminOnlyMiddleware(), userController.GetPhoneRecordByPhone)

	}

	// 文章相关路由
	articleGroup := r.Group("/article")
	{
		articleGroup.POST("/add-article", articleController.AddArticle)
		articleGroup.POST("/page-get-article-list", articleController.PageGet)
		articleGroup.POST("/get-article", articleController.ArticleGet)
	}

	// 活动相关路由
	activityGroup := r.Group("/activity")
	{
		activityGroup.GET("/list", sharedResourcesController.GetActivityList)
		activityGroup.POST("/add", sharedResourcesController.AddActivity)
		activityGroup.PUT("/update", sharedResourcesController.UpdateActivity)
		activityGroup.DELETE("/delete/:id", sharedResourcesController.DeleteActivity)
	}

	// 首页相关路由
	homePageGroup := r.Group("/homepage")
	{
		homePageGroup.GET("/banner", homePageController.GetBanner)
		homePageGroup.POST("/banner", homePageController.AddBanner)
		homePageGroup.PUT("/banner", homePageController.UpdateBanner)
		homePageGroup.DELETE("/banner/:id", homePageController.DeleteBanner)
	}

	return r
}
