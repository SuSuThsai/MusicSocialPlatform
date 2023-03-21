package Model

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	Replys      []Comment     `gorm:"foreignKey:RId"json:"replys"`
	CommentLike []CommentLike `gorm:"foreignKey:CommentId"json:"comment_like"`
	gorm.Model
	ArticleId  uint      `gorm:"type:uint;not null"json:"article_id"`
	Type       int       `gorm:"type:int;not null;default:0"json:"type"`
	ParentId   uint      `gorm:"type:uint;DEFAULT:0"json:"parent_id"`
	RId        uint      `gorm:"type:uint;DEFAULT:0"json:"r_id"`
	UserId     string    `json:"user_id"`
	Username   string    `json:"username"`
	ToUserId   string    `json:"to_user_id"`
	ToUsername string    `json:"to_username"`
	PRId       uint      `gorm:"type:uint;DEFAULT:0"json:"pr_id"`
	Content    string    `gorm:"type:varchar(500);notnull;"json:"content"`
	PageUrl    string    `gorm:"type:varchar(256)"json:"page_url"`
	LikeCount  int       `gorm:"type:int;not null;default:0"json:"like_count"`
	CreatedAt  time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"type:timestamp"`
}

type CommentLike struct {
	gorm.Model
	UserId    string `json:"user_id"`
	CommentId uint   `json:"comment_id"`
	ArticleId uint   `gorm:"type:uint;not null"json:"article_id"`
	Like      bool   `json:"like"`
}
