package activity_model

import (
	"time"
)

type Activity struct {
	Id            int       `gorm:"primary_key;auto_increment" json:"id" form:"id"`
	Title         string    `json:"title" form:"title"`
	Content       string    `json:"content" form:"content"`
	CoverImageUrl string    `json:"coverImageUrl" form:"coverImageUrl"`
	PublisherId   string    `gorm:"foreign_key:id" json:"publisherId" form:"publisherId"`
	CreatedAt     time.Time `json:"createdAt"`
}

func (Activity) TableName() string {
	return "community_activity"
}

// BriefActivity 活动的简要信息, 包括title, 图片链接和活动id
type BriefActivity struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`
	CoverImageUrl string `json:"coverImageUrl"`
}

type ActivityResponse struct {
	Activities []BriefActivity `json:"activities"`
	Page       int             `json:"page"`
	TotalPages int             `json:"totalPages"`
}
