package model

import (
	"time"
	"west2/biz/model/video"
)

var dateFormat string = "2006-01-02T15:04:05.000Z"

type Video struct {
	Id           string    `gorm:"type:varchar(100);primaryKey"`
	Uid          string    `gorm:"type:varchar(100)"`
	Title        string    `gorm:"type:varchar(100);not null"`
	Description  string    `gorm:"type:varchar(256);not null"`
	VideoUrl     string    `gorm:"type:varchar(256);unique;not null"`
	CoverUrl     string    `gorm:"type:varchar(256)"`
	VisitCount   int64     `gorm:"type:int;default:0"`
	LikeCount    int64     `gorm:"type:int;default:0"`
	CommentCount int64     `gorm:"type:int;default:0"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	DeletedAt    time.Time `gorm:"type:datetime;default:null"`
}

func VideoToResVideo(v *Video) *video.Video {
	visitCount := v.VisitCount
	likeCount := v.LikeCount
	commentCount := v.CommentCount
	return &video.Video{
		Id:           v.Id,
		Uid:          v.Uid,
		VideoUrl:     v.VideoUrl,
		CoverUrl:     v.CoverUrl,
		Title:        v.Title,
		Description:  v.Description,
		VisitCount:   &visitCount,
		LikeCount:    &likeCount,
		CommentCount: &commentCount,
		CreatedAt:    v.CreatedAt.Format(dateFormat),
		UpdatedAt:    v.UpdatedAt.Format(dateFormat),
		DeletedAt:    v.DeletedAt.Format(dateFormat),
	}
}

func VideosToResVideos(videos []*Video) []*video.Video {
	var videosRes []*video.Video
	for _, v := range videos {
		videosRes = append(videosRes, VideoToResVideo(v))
	}
	return videosRes
}
