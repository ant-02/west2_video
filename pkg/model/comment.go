package model

import (
	"time"
	"west2/biz/model/comment"
)

type Comment struct {
	Id         string    `gorm:"type:varchar(100);primaryKey"`
	VideoId    string    `gorm:"type:varchar(100)"`
	Uid        string    `gorm:"type:varchar(100)"`
	ParentId   string    `gorm:"type:varchar(100);default:null"`
	LikeCount  int64     `gorm:"type:int;default:0"`
	ChildCount int64     `gorm:"type:int;default:0"`
	Content    string    `gorm:"type:varchar(1000);null not"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
	DeletedAt  time.Time `gorm:"type:datetime;default:null"`
}

func CommentToresComment(c *Comment) *comment.Comment {
	if c == nil {
		return nil
	}
	likeCount, childCount := c.LikeCount, c.ChildCount

	return &comment.Comment{
		Id:         c.Id,
		Uid:        c.Uid,
		VideoId:    c.Uid,
		ParentId:   c.ParentId,
		LikeCount:  &likeCount,
		ChildCount: &childCount,
		Content:    c.Content,
		CreatedAt:  c.CreatedAt.Format(dateFormat),
		UpdatedAt:  c.UpdatedAt.Format(dateFormat),
		DeletedAt:  c.DeletedAt.Format(dateFormat),
	}
}

func CommentsToResComments(cts []*Comment) []*comment.Comment {
	var comments []*comment.Comment
	for _, c := range cts {
		comments = append(comments, CommentToresComment(c))
	}
	return comments
}
