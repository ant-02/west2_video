package service

import (
	"log"
	"west2/pkg/model"
	"west2/pkg/repository"
	"west2/util"
)

type videoService struct {
	vr repository.VideoRepository
}

type VideoService interface {
	GetVideoStream(latestTime string) ([]*model.Video, error)
	Publish(title, description, data, uid string) error
	GetVideosByUid(uid string, pageNum, pageSize int64) ([]*model.Video, error)
	GetVideosByVisitCount(pageNum, pageSize int64) ([]*model.Video, error)
	Search(keywords, fromDate, toDate, username string, pageNum, pageSize int64) ([]*model.Video, error)
}

func NewVideoService(vr repository.VideoRepository) VideoService {
	return &videoService{vr: vr}
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

func (vs *videoService) GetVideosByUid(uid string, pageNum, pageSize int64) ([]*model.Video, error) {
	videos, err := vs.vr.GetVideosByUid(uid, pageNum, pageSize)
	if err != nil {
		log.Printf("failed to get videos by uid: %s, error: %v", uid, err)
		return nil, err
	}

	return videos, nil
}

func (vs *videoService) GetVideosByVisitCount(pageNum, pageSize int64) ([]*model.Video, error) {
	videos, err := vs.vr.GetVideosGroupByVisitCount(pageNum, pageSize)
	if err != nil {
		log.Printf("failed to get videos by visit count: error: %v", err)
		return nil, err
	}

	return videos, nil
}

func (vs *videoService) Search(keywords, fromDate, toDate, username string, pageNum, pageSize int64) ([]*model.Video, error) {
	videos, err := vs.vr.GetVideosByKeywords(keywords, fromDate, toDate, username, pageNum, pageSize)
	if err != nil {
		log.Printf("failed to search videos: error: %v", err)
		return nil, err
	}

	return videos, nil
}
