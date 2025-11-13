package model

import "time"

type Video struct {
	Id           string    `gorm:"type:varchar(100);primaryKey"`
	Uid          string    `gorm:"type:varchar(100)"`
	Title        string    `gorm:"type:varchar(100);not null"`
	Description  string    `gorm:"type:varchar(256);not null"`
	VideoUrl     string    `gorm:"type:varchar(256);unique;not null"`
	CoverUrl     string    `gorm:"type:varchar(256);unique"`
	VisitCount   int64     `gorm:"type:int;default:0"`
	LikeCount    int64     `gorm:"type:int;default:0"`
	CommentCount int64     `gorm:"type:int;default:0"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	DeletedAt    time.Time `gorm:"type:datetime;default:null"`
}
