package activity_service

import (
	"errors"
	"huancuilou/internal/homePage/activity/activity_model"
	"huancuilou/internal/homePage/activity/activity_repository"
	"time"
)

type ActivityService struct {
	activityRepository *activity_repository.ActivityRepository
}

func NewActivityService(repository *activity_repository.ActivityRepository) *ActivityService {
	return &ActivityService{activityRepository: repository}
}

// GetActivityList 获得活动的简要信息, 包括title, 图片链接和活动id
func (as *ActivityService) GetActivityList(page int) (*activity_model.ActivityResponse, error) {
	return as.activityRepository.GetActivityList(page)
}

// GetActivity 获得活动详细信息
func (as *ActivityService) GetActivity(activityId int) (*activity_model.Activity, error) {
	return as.activityRepository.GetActivity(activityId)
}

// PostActivity 发布活动
func (as *ActivityService) PostActivity(activity *activity_model.Activity) error {
	// 检测关键字段是否为空
	if len((*activity).Title) == 0 || len((*activity).Content) == 0 {
		return errors.New("CommunityPostsService.PostActivity error: 400:the title or introduction is empty")
	}
	// 实现微信安全审核
	// TODO

	// 设定创建时间
	(*activity).CreatedAt = time.Now()
	return as.activityRepository.PostActivity(activity)
}
