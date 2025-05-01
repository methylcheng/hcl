package activity_routers

import (
	"github.com/gin-gonic/gin"
	"huancuilou/common"
	"huancuilou/internal/article/article_controller"
	"huancuilou/internal/homePage/homePage_controller"
)

func SetActivityRouters(sharedResourcesController *article_controller.SharedResourcesController, controller *homePage_controller.HomePageController) *gin.Engine {
	r := gin.Default()
	activityGroup := r.Group("/activity")
	activityGroup.POST("/publish", common.AdminOnlyMiddleware(), sharedResourcesController.PostActivity)
	activityGroup.GET("activityList/:page", controller.GetActivityList)
	activityGroup.GET("activity/:id", controller.GetActivity)
	return r
}
