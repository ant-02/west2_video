package service

import (
	"errors"
	"log"
	"west2/pkg/model"
	"west2/pkg/repository"
	"west2/util"

	"gorm.io/gorm"
)

type videoService struct {
	vr repository.VideoRepository
	ur repository.UserRepository
}

type VideoService interface {
	GetVideoStream(latestTime string) ([]*model.Video, error)
	Publish(title, description, data, uid string) error
	GetVideosByUid(uid string, pageNum, pageSize int64) ([]*model.Video, int64, error)
	GetVideosByVisitCount(pageNum, pageSize int64) ([]*model.Video, error)
	Search(keywords, fromDate, toDate, username string, pageNum, pageSize int64) ([]*model.Video, int64, error)
}

func NewVideoService(vr repository.VideoRepository, ur repository.UserRepository) VideoService {
	return &videoService{vr: vr, ur: ur}
}

func (vs *videoService) GetVideoStream(latestTime string) ([]*model.Video, error) {
	videos, err := vs.vr.GetVideosByLatestTime(latestTime)
	if err != nil {
		log.Printf("failed to get video steam: latestTime: %s, error: %v", latestTime, err)
		return nil, err
	}

	return videos, nil
}

func (vs *videoService) Publish(title, description, data, uid string) error {
	id := util.GetID()

	if err := util.Base64ToVideo(data, "./static/video/"+id+".mp4"); err != nil {
		log.Printf("failed to save video file: id: %s, error: %v", id, err)
		return err
	}

	if err := vs.vr.CreateVideo(&model.Video{
		Id:          id,
		Uid:         uid,
		Title:       title,
		Description: description,
		VideoUrl:    "/static/video/" + id + ".mp4",
	}); err != nil {
		log.Printf("failed to create video: error: %v", err)
		return err
	}

	return nil
}

func (vs *videoService) GetVideosByUid(uid string, pageNum, pageSize int64) ([]*model.Video, int64, error) {
	videos, total, err := vs.vr.GetVideosByUid(uid, pageNum, pageSize)
	if err != nil {
		log.Printf("failed to get videos by uid: %s, error: %v", uid, err)
		return nil, 0, err
	}

	return videos, total, nil
}

func (vs *videoService) GetVideosByVisitCount(pageNum, pageSize int64) ([]*model.Video, error) {
	videos, err := vs.vr.GetVideosGroupByVisitCount(pageNum, pageSize)
	if err != nil {
		log.Printf("failed to get videos by visit count: error: %v", err)
		return nil, err
	}

	return videos, nil
}

func (vs *videoService) Search(keywords, fromDate, toDate, username string, pageNum, pageSize int64) ([]*model.Video, int64, error) {
	var u *model.User
	var err error
	if username != "" {
		u, err = vs.ur.GetUserByUsername(username)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, 0, nil
			}
			log.Printf("failed to get user by username: %v", err)
			return nil, 0, err
		}
	}

	var uid string
	if u != nil {
		uid = u.Id
	}
	videos, total, err := vs.vr.GetVideosByKeywords(keywords, fromDate, toDate, uid, pageNum, pageSize)
	if err != nil {
		log.Printf("failed to search videos: error: %v", err)
		return nil, 0, err
	}

	return videos, total, nil
}
