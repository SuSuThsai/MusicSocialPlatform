package Model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	UserInfo    UserInfo      `gorm:"foreignKey:UserId"`
	MusicList   []MusicList   `gorm:"foreignKey:UserId"`
	FollowList  []FollowList  `gorm:"foreignKey:UserId"`
	ConcernList []ConcernList `gorm:"foreignKey:UserId"`
	MusicLike   []MusicLike   `gorm:"foreignKey:UserId"`
	SongList    []SongList    `gorm:"foreignKey:UserId"`
	Article     []Article     `gorm:"foreignKey:UserId;references:UserId"`
	gorm.Model
	UserId    string    `gorm:"type:varchar(20);uniqueIndex;not null"json:"user_id" validate:"required,min=1,max=20"`
	Username  string    `gorm:"type:varchar(20);not null"json:"username"validate:"required,min=2,max=20"`
	Password  string    `gorm:"type:varchar(20);not null"json:"password"validate:"required,min=6,max=20"`
	Role      int       `gorm:"type:int;DEFAULT:2"json:"role"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp"`
}

type UserInfo struct {
	gorm.Model
	UserId     uint   `gorm:"type:uint;Index"json:"user_id"`
	Sex        int    `gorm:"type:int;DEFAULT:3"json:"sex"`
	Desc       string `gorm:"type:varchar(20)"json:"desc"`
	Pfp        string `gorm:"type:varchar(255)"json:"pfp"`
	Background string `gorm:"type:varchar(255)"json:"background"`
}

type LoginMemory struct {
	gorm.Model
	UserId    uint      `gorm:"type:uint;Index"json:"user_id"`
	LoginTime time.Time `gorm:"type:uint;Index"json:"login_time"`
}

type UserListenTypeCount struct {
	UserId      string `gorm:"type:varchar(20);Index;not null"json:"user_id"`
	Habits      string `gorm:"type:varchar(20);Index"json:"habits"`
	ListenCount int    `gorm:"type:int;not null;default:0"json:"listen_count"`
}

type UserListenMusicCount struct {
	gorm.Model
	UserId      string    `gorm:"type:varchar(20);Index;not null"json:"user_id"`
	MusicId     uint      `gorm:"type:uint;Index"json:"music_id"`
	ListenCount int       `gorm:"type:int;not null;default:0"json:"listen_count"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamp"`
}

type UserListeningHabits struct {
	UserId string `gorm:"type:varchar(20);Index;not null"json:"user_id"`
	Habits string `gorm:"type:varchar(20);Index"json:"habits"`
}

type MusicList struct {
	SongList      []SongList      `gorm:"foreignKey:ListId"`
	MusicListLike []MusicListLike `gorm:"foreignKey:ListId"`
	tips          []Tips          `gorm:"foreignKey:ListId"`
	gorm.Model
	UserId    uint      `gorm:"type:uint;Index"json:"user_id"`
	Desc      string    `gorm:"type:varchar(20)"json:"desc"`
	LName     string    `gorm:"type:varchar(15);not null;Index;"json:"l_name" validate:"required,min=1,max=15"`
	LikeCount int       `gorm:"type:int;DEFAULT=0;"json:"like_count"`
	Img       string    `gorm:"type:varchar(255)"json:"img"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp"`
}

type Tips struct {
	ListId    uint      `gorm:"type:uint;Index"json:"list_id"`
	TipName   string    `gorm:"type:varchar(15);Index;"json:"tip_name"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp"`
}

type MusicCommandIsListen struct {
	UserId   string `gorm:"type:varchar(20);Index"json:"user_id"`
	MusicId  uint   `gorm:"type:uint;Index"json:"music_id"`
	IsListen int    `gorm:"default:0";json:"is_listen"`
}

type CommandMusicListenCount struct {
	UserId    string    `gorm:"type:varchar(20);Index"json:"user_id"`
	Year      string    `gorm:"index;not null" json:"year"`
	Month     string    `gorm:"index;not null" json:"month"`
	Week      string    `gorm:"index;not null" json:"week"`
	Day       string    `gorm:"index;not null" json:"day"`
	Number    string    `gorm:"index;not null" json:"number"`
	IsListen  int       `gorm:"default:0";json:"is_listen"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp"`
}

type SongList struct {
	gorm.Model
	UserId    uint      `gorm:"type:uint;Index"json:"user_id"`
	ListId    uint      `gorm:"type:uint;Index"json:"list_id"`
	MusicId   uint      `gorm:"unique"json:"music_id"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp"`
}

type CommandMusicCount struct {
	UserId    string    `gorm:"type:varchar(20);Index"json:"user_id"`
	MusicId   uint      `gorm:"type:uint;Index"json:"music_id"`
	Year      string    `gorm:"index;not null" json:"year"`
	Month     string    `gorm:"index;not null" json:"month"`
	Week      string    `gorm:"index;not null" json:"week"`
	Day       string    `gorm:"index;not null" json:"day"`
	Number    string    `gorm:"index;not null" json:"number"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp"`
}

type MusicListLike struct {
	gorm.Model
	UserId uint `gorm:"type:uint;Index"json:"user_id"`
	ListId uint `gorm:"type:uint;Index"json:"list_id"`
	Like   bool `gorm:"type:bool"json:"like"`
}

type MusicLike struct {
	gorm.Model
	UserId  uint `gorm:"type:uint;Index"json:"user_id"`
	MusicId uint `json:"music_id"`
	Like    bool `gorm:"type:bool"json:"like"`
}

type FollowList struct {
	gorm.Model
	UserId     uint   `gorm:"type:uint;Index"json:"user_id"`
	FollowerId string `gorm:"type:varchar(13);unique"json:"follower_id"`
	Like       bool   `gorm:"type:bool"json:"like"`
}

type ConcernList struct {
	gorm.Model
	UserId    string `gorm:"type:varchar(20);Index"json:"user_id"`
	ConcernId string `gorm:"type:varchar(20);"json:"concern_id"`
	Like      bool   `gorm:"type:bool"json:"like"`
}
