package repository

import (
	"west2/pkg/model"

	"gorm.io/gorm"
)

type likeRepository struct {
	db *gorm.DB
}

type LikeRepository interface {
	CreateLike(like *model.Like) error
	GetLike(commentId, videoId, uid string) (*model.Like, error)
	SetLikeStatus(id string, status int64) error
	GetVideoLikeList(uid string, pageNum, pageSize int64) ([]*string, error)
}

func NewLikeReposirty(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

func (lr *likeRepository) CreateLike(like *model.Like) error {
	return lr.db.Create(like).Error
}

func (lr *likeRepository) GetLike(commentId, videoId, uid string) (*model.Like, error) {
	tx := lr.db.Where("uid = ?", uid).
		Where("deleted_at IS NULL")

	if commentId == "" {
		tx = tx.Where("video_id = ?", videoId)
	} else {
		tx = tx.Where("comment_id = ?", commentId)
	}

	var like model.Like
	err := tx.First(&like).Error
	if err != nil {
		return nil, err
	}
	return &like, nil
}

func (lr *likeRepository) SetLikeStatus(id string, status int64) error {
	return lr.db.Model(&model.Like{}).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Update("status", status).Error
}

func (lr *likeRepository) GetVideoLikeList(uid string, pageNum, pageSize int64) ([]*string, error) {
	var videoIds []*string
	err := lr.db.Model(&model.Like{}).
		Select("video_id").
		Where("uid = ?", uid).
		Where("video_id IS NOT NULL").
		Where("status = 1").
		Where("deleted_at IS NULL").
		Offset((int(pageNum) - 1) * int(pageSize)).
		Limit(int(pageSize)).
		Find(&videoIds).Error
	if err != nil {
		return nil, err
	}

	return videoIds, nil
}
