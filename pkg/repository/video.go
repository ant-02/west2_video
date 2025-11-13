package repository

import (
	"encoding/json"
	"west2/database"
	"west2/pkg/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	key string = "video:list:visitCount"
)

type videoRepository struct {
	db *gorm.DB
}

type VideoRepository interface {
	GetVideosByLatestTime(latestTime string) ([]*model.Video, error)
	CreateVideo(video *model.Video) error
	GetVideosByUid(uid string, pageNum, pageSize int64) ([]*model.Video, error)
	GetVideosGroupByVisitCount(pageNum, pageSize int64) ([]*model.Video, error)
	GetVideosByKeywords(keywords, fromDate, toDate, username string, pageNum, pageSize int64) ([]*model.Video, error)
}

func NewVideoRepository(db *gorm.DB) VideoRepository {
	return &videoRepository{db: db}
}

func (vr *videoRepository) GetVideosByLatestTime(latestTime string) ([]*model.Video, error) {
	var videos []*model.Video

	err := vr.db.Where("created_at > ?", latestTime).
		Where("deleted_at IS NULL").
		Find(&videos).Error
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (vr *videoRepository) CreateVideo(video *model.Video) error {
	err := vr.db.Create(video).Error
	if err != nil {
		return err
	}

	return database.Del([]string{key})
}

func (vr *videoRepository) GetVideosByUid(uid string, pageNum, pageSize int64) ([]*model.Video, error) {
	var videos []*model.Video
	err := vr.db.Where("uid = ?", uid).
		Where("deleted_at IS NULL").
		Offset((int(pageNum) - 1) * int(pageSize)).
		Limit(int(pageSize)).
		Find(&videos).Error
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (vr *videoRepository) GetVideosGroupByVisitCount(pageNum, pageSize int64) ([]*model.Video, error) {
	var videos []*model.Video
	first, end := (pageNum-1)*pageSize, pageNum*pageSize
	videoStrings, err := database.LRange(key, first, end)
	if err == nil {
		for _, s := range videoStrings {
			var v model.Video
			if err := json.Unmarshal([]byte(s), &v); err != nil {
				return nil, err
			}
			videos = append(videos, &v)
		}
		return videos, nil
	}
	if err != redis.Nil {
		return nil, err
	}

	err = vr.db.Where("deleted_at IS NULL").
		Order("visit_count desc").
		Offset(int(first)).
		Limit(int(pageSize)).
		Find(&videos).Error
	if err != nil {
		return nil, err
	}

	err = database.RPush(key, videos)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (vr *videoRepository) GetVideosByKeywords(keywords, fromDate, toDate, username string, pageNum, pageSize int64) ([]*model.Video, error) {
	var videos []*model.Video
	tx := vr.db.Where("title LIKE ? or description LIKE ?", "%"+keywords+"%", "%"+keywords+"%")

	if fromDate != "" {
		tx = tx.Where("from_date > ?", fromDate)
	}
	if toDate != "" {
		tx = tx.Where("to_date < ?", toDate)
	}
	if username != "" {
		tx = tx.Where("username = ?", username)
	}

	err := tx.Offset((int(pageNum) - 1) * int(pageSize)).
		Limit(int(pageSize)).
		Find(&videos).Error

	if err != nil {
		return nil, err
	}

	return videos, nil
}
