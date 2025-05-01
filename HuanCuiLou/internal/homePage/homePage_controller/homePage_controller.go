package homePage_controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"huancuilou/common/error_handler"
	"huancuilou/internal/homePage/activity/activity_service"
	"huancuilou/response"
	"strconv"
)

//处理首页的相关接口
//4.首页活动列表
//5.首页通知列表
//6.评论、点赞、收藏、入库

type HomePageController struct {
	activityService *activity_service.ActivityService
}

func NewHomePageController(activityService *activity_service.ActivityService) *HomePageController {
	return &HomePageController{activityService: activityService}
}

// GetActivityList 首页活动列表
func (hpc *HomePageController) GetActivityList(c *gin.Context) {
	// 获得页码
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		page = 1 // 设定默认值
	}
	// 获得活动列表
	activityListResp, err := hpc.activityService.GetActivityList(page)
	if err != nil {
		error_handler.HandleUserError(c, err)
		return
	}
	c.JSON(200, response.Success(*activityListResp))
}

// GetActivity 获得活动详细信息
func (hpc *HomePageController) GetActivity(c *gin.Context) {
	// 获得activityId
	activityId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		error_handler.HandleUserError(c, errors.New("HomePageController.GetActivity error: 400: 活动id获取失败"))
		return
	}
	// 获得活动详细信息
	activityResp, err := hpc.activityService.GetActivity(activityId)
	if err != nil {
		error_handler.HandleUserError(c, err)
		return
	}
	c.JSON(200, response.Success(*activityResp))
}
