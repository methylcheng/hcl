package activity_repository

import (
	"fmt"
	"gorm.io/gorm"
	"huancuilou/internal/homePage/activity/activity_model"
)

type ActivityRepository struct {
	DB *gorm.DB
}

func NewActivityRepository(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{DB: db}
}

var pageSize int = 10 // 分页查询时每一页的数据量

// GetActivityList 获得活动的简要信息, 包括title, 图片链接和活动id
func (car *ActivityRepository) GetActivityList(page int) (*activity_model.ActivityResponse, error) {
	var activity []activity_model.BriefActivity
	db := car.DB.Model(&activity_model.Activity{})
	// 进行分页查询
	offset := (page - 1) * pageSize
	limit := page * pageSize
	// 从数据库中读取数据
	err := db.Select("id", "title", "cover_image_url").Limit(limit).Offset(offset).Find(&activity).Error
	if err != nil {
		return nil, fmt.Errorf("ActivityRepository.GetActivityList error: 404:%v", err)
	}
	var totalActivity, totalPages int64
	err = db.Count(&totalActivity).Error
	totalPages = totalActivity/10 + 1
	if err != nil {
		return nil, fmt.Errorf("ActivityRepository.GetActivityList error: 400:%v", err)
	}
	resp := &activity_model.ActivityResponse{Activities: activity, Page: page, TotalPages: int(totalPages)}
	return resp, nil
}

// GetActivity 获得活动的详细信息
func (car *ActivityRepository) GetActivity(activityId int) (*activity_model.Activity, error) {
	var activity activity_model.Activity
	db := car.DB.Model(&activity_model.Activity{})
	err := db.Where("id = ?", activityId).First(&activity).Error
	if err != nil {
		return nil, fmt.Errorf("ActivityRepository.GetActivity error: 404:%v", err)
	}
	return &activity, nil
}

// PostActivity 发布活动
func (car *ActivityRepository) PostActivity(activity *activity_model.Activity) error {
	db := car.DB.Model(&activity_model.Activity{})
	err := db.Create(activity).Error
	if err != nil {
		return fmt.Errorf("ActivityRepository.PostActivity error: 400:%v", err)
	}
	return nil
}
