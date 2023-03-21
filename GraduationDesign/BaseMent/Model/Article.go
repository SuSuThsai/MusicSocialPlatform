package Model

import (
	"gorm.io/gorm"
	"time"
)

type Article struct {
	Comments    []Comment     `gorm:"foreignKey:ArticleId"`
	_           []Forward     `gorm:"foreignKey:ArticleId"`
	ArticleLike []ArticleLike `gorm:"foreignKey:ArticleId"`
	gorm.Model
	UserId       string    `gorm:"type:varchar(13);Index;not null"json:"user_id" validate:"required,min=1,max=13"`
	Topic1       string    `gorm:"type:varchar(12)"json:"topic1"`
	Type         int       `gorm:"type:int;not null;default:0"json:"type"`
	Topic2       string    `gorm:"type:varchar(12)"json:"topic2""`
	Topic3       string    `gorm:"type:varchar(12)"json:"topic3"`
	Content      string    `gorm:"type:text"json:"content"`
	Img          string    `gorm:"type:varchar(250)"json:"img"`
	Audio        string    `gorm:"type:varchar(250)"json:"audio"`
	Forward      int       `gorm:"type:int;not null;default:0"json:"forward"`
	Link         string    `gorm:"type:varchar(256)"json:"link"`
	LikeCount    int       `gorm:"type:int;not null;default:0"json:"like_count"`
	CommentCount int       `gorm:"type:int;not null;default:0"json:"comment_count"`
	ReadCount    int       `gorm:"type:int;not null;default:0"json:"read_count"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:timestamp"`
}

type Forward struct {
	ArticleId uint   `gorm:"type:uint;not null;uniqueIndex"`
	UserId    string `gorm:";not null;uniqueIndex"`
}

type Topic struct {
	gorm.Model
	TName        string `gorm:"type:varchar(12);uniqueIndex;not null"json:"t_name"`
	Use          int    `gorm:"type:int;not null;default:0"json:"use"`
	CommentCount int    `gorm:"type:int;not null;default:0"json:"comment_count"`
	ReadCount    int    `gorm:"type:int;not null;default:0"json:"read_count"`
}

type ArticleLike struct {
	gorm.Model
	UserId    string `json:"user_id"`
	ArticleId uint   `gorm:"type:uint;not null"json:"article_id"`
	Like      bool   `json:"like"`
}
