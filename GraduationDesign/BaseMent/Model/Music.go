package Model

import (
	"gorm.io/gorm"
	"time"
)

type Music struct {
	gorm.Model
	MusicTopic  []MusicTopic `gorm:"foreignKey:Id"`
	MusicListen MusicListen  `gorm:"foreignKey:Id"`
	Id          uint         `gorm:"primary_key;BIGSERIAL;not null" json:"id"`
	AlbumId     uint         `json:"album_id"`
	Name        string       `gorm:"not null"json:"name"validate:"required,min=1"`
	Source      string       `json:"source"`
	WordsSource string       `json:"words_source"`
	Singer      string       `gorm:"not null;DEFAULT:'未知'"json:"singer"`
	//Like        uint   `gorm:"type:int;DEFAULT:0"json:"like"`
	Time      string    `gorm:"type:varchar(8)"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp"`
}

type MusicTopic struct {
	Id   uint   `gorm:"index;not null" json:"id"`
	Name string `gorm:"index;not null"json:"name"`
	Tip  string `gorm:"index;not null"json:"tip"`
}

type MusicListen struct {
	Id     uint `gorm:"primary_key;not null" json:"id"`
	Listen int  `gorm:"type:int;not null;default:0"json:"listen"`
}

type Album struct {
	gorm.Model
	Id        uint      `gorm:"primary_key;BIGSERIAL;not null" json:"id"`
	SName     string    `gorm:"not null;DEFAULT:'未知'"json:"s_name"`
	Name      string    `gorm:"not null"json:"name"validate:"required,min=1"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp"`
}

type MusicRankYear struct {
	Id    uint   `gorm:"index;not null" json:"id"`
	Year  string `gorm:"index;not null" json:"year"`
	Count int    `gorm:"index;not null" json:"count"`
}

type MusicRankMonth struct {
	Id    uint   `gorm:"index;not null" json:"id"`
	Year  string `gorm:"index;not null" json:"year"`
	Month string `gorm:"index;not null" json:"month"`
	Count int    `gorm:"index;not null" json:"count"`
}

type MusicRankWeek struct {
	Id    uint   `gorm:"index;not null" json:"id"`
	Year  string `gorm:"index;not null" json:"year"`
	Month string `gorm:"index;not null" json:"month"`
	Week  string `gorm:"index;not null" json:"week"`
	Count int    `gorm:"index;not null" json:"count"`
}

type MusicRankDay struct {
	Id     uint   `gorm:"index;not null" json:"id"`
	Year   string `gorm:"index;not null" json:"year"`
	Month  string `gorm:"index;not null" json:"month"`
	Week   string `gorm:"index;not null" json:"week"`
	Day    string `gorm:"index;not null" json:"day"`
	Number string `gorm:"index;not null" json:"number"`
	Count  int    `gorm:"index;not null" json:"count"`
}

type MusicListRankYear struct {
	Id    uint   `gorm:"index;not null" json:"id"`
	Year  string `gorm:"index;not null" json:"year"`
	Count int    `gorm:"index;not null" json:"count"`
}

type MusicListRankMonth struct {
	Id    uint   `gorm:"index;not null" json:"id"`
	Year  string `gorm:"index;not null" json:"year"`
	Month string `gorm:"index;not null" json:"month"`
	Count int    `gorm:"index;not null" json:"count"`
}

type MusicListRankWeek struct {
	Id    uint   `gorm:"index;not null" json:"id"`
	Year  string `gorm:"index;not null" json:"year"`
	Month string `gorm:"index;not null" json:"month"`
	Week  string `gorm:"index;not null" json:"week"`
	Count int    `gorm:"index;not null" json:"count"`
}

type MusicListRankDay struct {
	Id     uint   `gorm:"index;not null" json:"id"`
	Year   string `gorm:"index;not null" json:"year"`
	Month  string `gorm:"index;not null" json:"month"`
	Week   string `gorm:"index;not null" json:"week"`
	Day    string `gorm:"index;not null" json:"day"`
	Number string `gorm:"index;not null" json:"number"`
	Count  int    `gorm:"index;not null" json:"count"`
}
