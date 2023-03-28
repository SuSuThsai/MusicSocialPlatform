package Sql

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	MessageId  uint   `gorm:"primaryKey;autoIncrement"json:"message_id" `
	SenderId   uint   `gorm:"type:uint;Index;not null"json:"sender_id" `
	ReceiverId uint   `gorm:"type:uint;Index;not null"json:"receiver_id"`
	Content    string `gorm:"type:varchar(255);not null"json:"content"`
	Read       bool   `gorm:"type:bool;DEFAULT:false"json:"read"`
}
