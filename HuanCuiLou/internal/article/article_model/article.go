package article_model

import "time"

// Article 文章表
type Article struct {
	ID        int       `json:"id"`
	Title     string    `json:"title" gorm:"column:username" `
	Content   string    `json:"content"`
	ManagerId int       `json:"manager_id"`
	CreatedAt time.Time `json:"created_at"`
	Kind      string    `json:"kind"`
	LikeCount string    `json:"like_count"`
}

// PageGet 分页查询时使用的结构体
type PageGet struct {
	PageNum  int    `json:"page_num"`  //页码数
	PageSize int    `json:"page_size"` //每页显示的个数
	Kind     string `json:"kind"`
}

// TableName 自定义表名
func (Article) TableName() string {
	return "article"
}
