package model

import "time"

type Like struct {
	Id        string    `gorm:"type:varchar(100);primaryKey"`
	Uid       string    `gorm:"type:varchar(100)"`
	VideoId   string    `gorm:"type:varchar(100);default:null"`
	CommentId string    `gorm:"type:varchar(100);default:null"`
	Status    int64     `gorm:"type:int;default:1"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt time.Time `gorm:"type:datetime;default:null"`
}
