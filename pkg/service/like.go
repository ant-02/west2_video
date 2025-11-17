package service

import (
	"errors"
	"log"
	"west2/pkg/model"
	"west2/pkg/repository"
	"west2/util"

	"gorm.io/gorm"
)

type likeService struct {
	lr repository.LikeRepository
	vr repository.VideoRepository
}

type LikeService interface {
	LikeAction(like *model.Like) error
	GetVideoListByLike(uid string, pageNum, pageSize int64) ([]*model.Video, error)
}

func NewLikeService(lr repository.LikeRepository, vr repository.VideoRepository) LikeService {
	return &likeService{lr: lr, vr: vr}
}

func (ls *likeService) LikeAction(like *model.Like) error {
	l, err := ls.lr.GetLike(like.CommentId, like.VideoId, like.Uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if like.VideoId != "" {
				if like.Status == 1 {
					if err := ls.vr.AddLikeCount(like.VideoId); err != nil {
						log.Printf("failed to add like count: videoId: %s, error: %v", like.VideoId, err)
						return err
					}
				} else {
					if err := ls.vr.SubtractLikeCount(like.VideoId); err != nil {
						log.Printf("failed to suntract like count: videoId: %s, error: %v", like.VideoId, err)
						return err
					}
				}
			}

			like.Id = util.GetID()
			return ls.lr.CreateLike(like)
		}
		log.Printf("failed to get like by ids: commentId: %s, videoId: %s, error: %v", like.CommentId, like.VideoId, err)
		return err
	}

	if like.Status == 1 {
		if err := ls.vr.AddLikeCount(like.VideoId); err != nil {
			log.Printf("failed to add like count: videoId: %s, error: %v", like.VideoId, err)
			return err
		}
	} else {
		if err := ls.vr.SubtractLikeCount(like.VideoId); err != nil {
			log.Printf("failed to suntract like count: videoId: %s, error: %v", like.VideoId, err)
			return err
		}
	}
	err = ls.lr.SetLikeStatus(l.Id, like.Status)
	if err != nil {
		log.Printf("failed to set like status: likeId: %s, %v", l.Id, err)
	}
	return err
}

func (ls *likeService) GetVideoListByLike(uid string, pageNum, pageSize int64) ([]*model.Video, error) {
	ids, err := ls.lr.GetVideoLikeList(uid, pageNum, pageSize)
	if err != nil {
		log.Printf("failed to get video ids by user like: uid: %s, %v", uid, err)
	}

	videos, err := ls.vr.GetVideosByIds(ids)
	if err != nil {
		log.Printf("failed to get video by ids: uid: %v, %v", uid, err)
		return nil, err
	}
	return videos, err
}
