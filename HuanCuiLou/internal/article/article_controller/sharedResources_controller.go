package article_controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"huancuilou/common"
	"huancuilou/common/error_handler"
	"huancuilou/internal/homePage/activity/activity_model"
	"huancuilou/internal/homePage/activity/activity_service"
	"huancuilou/response"
	"net/http"
	"os"
	"strconv"
)

//管理员发布共享资源相关的接口

type SharedResourcesController struct {
	activityService *activity_service.ActivityService
}

func NewSharedResourcesController(activityService *activity_service.ActivityService) *SharedResourcesController {
	return &SharedResourcesController{activityService: activityService}
}

// PostActivity 活动发布接口
func (src *SharedResourcesController) PostActivity(c *gin.Context) {
	var activity activity_model.Activity
	// 获得publisherID
	publisherID, ok := c.Get("userID")
	if !ok {
		error_handler.HandleUserError(c, errors.New("sharedResourcesController.PostActivity  error: 401: can't get publisherID"))
	}
	// 获得活动的相关信息
	if err := c.ShouldBind(&activity); err != nil {
		error_handler.HandleUserError(c, fmt.Errorf("sharedResourcesController.PostActivity ShouldBindJSON error: 400:%v", err))
		return
	}
	activity.PublisherId = strconv.Itoa(publisherID.(int))
	// 获得封面图
	imageFile, err := c.FormFile("cover_image")
	// 如果存在错误&&不是未上传图片的错误
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		error_handler.HandleUserError(c, fmt.Errorf("sharedResourcesController.PostActivity post FormFile error: 400:%v", err))
		return
	}
	url := ""
	// 若用户上传了封面图
	if imageFile != nil {
		// 将封面图保存在本地临时目录
		dst := fmt.Sprintf("./cover_image/%s", imageFile.Filename)
		err = c.SaveUploadedFile(imageFile, dst)
		if err != nil {
			error_handler.HandleUserError(c, fmt.Errorf("sharedResourcesController.PostActivity saveFile error: 400:%v", err))
			return
		}
		// 结束时将临时封面图文件删除
		defer func(name *string) {
			_ = os.Remove(*name)
		}(&dst)
		// 将图片上传到OSS
		url, err = common.UploadToOSS(imageFile, dst)
		if err != nil {
			error_handler.HandleUserError(c, fmt.Errorf("sharedResourcesController.PostActivity UploadImage error: 400:%v", err))
		}
	}
	activity.CoverImageUrl = url
	// 进一步将进行service层处理
	err = src.activityService.PostActivity(&activity)
	if err != nil {
		error_handler.HandleUserError(c, err)
		return
	}
	// 服务成功
	c.JSON(200, response.SuccessWithoutData())
}
