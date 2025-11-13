package model

import "time"

type User struct {
	Id        string    `gorm:"type:varchar(100);primaryKey"`
	Username  string    `gorm:"type:varchar(100);unique;not null"`
	Password  string    `gorm:"type:varchar(100);not null"`
	AvatarUrl string    `gorm:"type:varchar(256)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt time.Time `gorm:"default:null"`
	Videos    []Video   `gorm:"foreignKey:Uid;references:Id"`
}
