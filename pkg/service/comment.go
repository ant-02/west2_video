package service

import (
	"log"
	"west2/pkg/model"
	"west2/pkg/repository"
)

type commentService struct {
	cr repository.CommentRepository
}

type CommentService interface {
	Publish(comment *model.Comment) error
	GetCommentList(videoId, commentId string, pageNum, pageSize int64) ([]*model.Comment, error)
	DeleteById(id string) error
	DeleteByVideoId(videoId string) error
}

func NewCommentService(cr repository.CommentRepository) CommentService {
	return &commentService{cr: cr}
}

func (cs *commentService) Publish(comment *model.Comment) error {
	err := cs.cr.CreateComment(comment)
	if err != nil {
		log.Printf("failed to create comment: comment: %v, err: %v", comment, err)
		return err
	}
	return nil
}

func (cs *commentService) GetCommentList(videoId, commentId string, pageNum, pageSize int64) ([]*model.Comment, error) {
	var comments []*model.Comment
	var err error
	if videoId != "" {
		comments, err = cs.cr.GetCommentListByVideoId(videoId, pageNum, pageSize)
	} else {
		comments, err = cs.cr.GetCommentListByCommentId(commentId, pageNum, pageSize)
	}
	if err != nil {
		log.Printf("failed to get comments: videoId: %s, commentId: %s, err: %v", videoId, commentId, err)
		return nil, err
	}
	return comments, nil
}

func (cs *commentService) DeleteById(id string) error {
	err := cs.cr.DeleteCommentById(id)
	if err != nil {
		log.Printf("failed to delete comment by id: id: %s, err: %v", id, err)
		return err
	}
	return nil
}

func (cs *commentService) DeleteByVideoId(videoId string) error {
	err := cs.cr.DeleteCommentsByVideoId(videoId)
	if err != nil {
		log.Printf("failed to delete comment by videoId: videoId: %s, err: %v", videoId, err)
		return err
	}
	return nil
}
