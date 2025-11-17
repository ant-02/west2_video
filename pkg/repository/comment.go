package repository

import (
	"west2/pkg/model"

	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

type CommentRepository interface {
	CreateComment(comment *model.Comment) error
	GetCommentListByVideoId(videoId string, pageNum, pageSize int64) ([]*model.Comment, error)
	GetCommentListByCommentId(commentId string, pageNum, pageSize int64) ([]*model.Comment, error)
	DeleteCommentsByVideoId(videoId string) error
	DeleteCommentById(id string) error
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (cr *commentRepository) CreateComment(comment *model.Comment) error {
	return cr.db.Create(comment).Error
}

func (cr *commentRepository) GetCommentListByVideoId(videoId string, pageNum, pageSize int64) ([]*model.Comment, error) {
	var comments []*model.Comment
	err := cr.db.Where("video_id = ?", videoId).
		Where("parent_id IS NULL").
		Offset((int(pageNum) - 1) * int(pageSize)).
		Limit(int(pageSize)).
		Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (cr *commentRepository) GetCommentListByCommentId(commentId string, pageNum, pageSize int64) ([]*model.Comment, error) {
	var comments []*model.Comment
	err := cr.db.Where("parent_id = ?", commentId).
		Offset((int(pageNum) - 1) * int(pageSize)).
		Limit(int(pageSize)).
		Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (cr *commentRepository) DeleteCommentsByVideoId(videoId string) error {
	return cr.db.Where("video_id = ?", videoId).
		Delete(&model.Comment{}).Error
}

func (cr *commentRepository) DeleteCommentById(id string) error {
	return cr.db.Delete(&model.Comment{}, id).Error
}
