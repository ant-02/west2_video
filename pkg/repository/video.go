package repository

import (
	"encoding/json"
	"strconv"
	"time"
	"west2/database"
	"west2/pkg/model"

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
	GetVideosByUid(uid string, pageNum, pageSize int64) ([]*model.Video, int64, error)
	GetVideosGroupByVisitCount(pageNum, pageSize int64) ([]*model.Video, error)
	GetVideosByKeywords(keywords, fromDate, toDate, username string, pageNum, pageSize int64) ([]*model.Video, int64, error)
}

func NewVideoRepository(db *gorm.DB) VideoRepository {
	return &videoRepository{db: db}
}

func (vr *videoRepository) GetVideosByLatestTime(latestTime string) ([]*model.Video, error) {
	var videos []*model.Video

	t, err := strconv.ParseInt(latestTime, 10, 64)
	if err != nil {
		return nil, err
	}
	err = vr.db.Where("created_at > ?", time.Unix(t, 0)).
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

func (vr *videoRepository) GetVideosByUid(uid string, pageNum, pageSize int64) ([]*model.Video, int64, error) {
	var videos []*model.Video
	var total int64
	var err error

	tx := vr.db.Model(&model.Video{}).
		Where("uid = ?", uid).
		Where("deleted_at IS NULL")

	err = tx.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = vr.db.Where("uid = ?", uid).
		Where("deleted_at IS NULL").
		Offset((int(pageNum) - 1) * int(pageSize)).
		Limit(int(pageSize)).
		Find(&videos).Error
	if err != nil {
		return nil, 0, err
	}

	return videos, total, nil
}

func (vr *videoRepository) GetVideosGroupByVisitCount(pageNum, pageSize int64) ([]*model.Video, error) {
	var videos []*model.Video
	first, end := (pageNum-1)*pageSize, pageNum*pageSize
	videoStrings, err := database.LRange(key, first, end)
	if err != nil || videoStrings == nil {
		return nil, err
	}
	if len(videoStrings) > 0 {
		for _, s := range videoStrings {
			var v model.Video
			if err := json.Unmarshal([]byte(s), &v); err != nil {
				return nil, err
			}
			videos = append(videos, &v)
		}
		return videos, nil
	}

	err = vr.db.Where("deleted_at IS NULL").
		Order("visit_count desc").
		Offset(int(first)).
		Limit(int(pageSize)).
		Find(&videos).Error
	if err != nil {
		return nil, err
	}

	for _, v := range videos {
		j, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		err = database.RPush(key, j)
		if err != nil {
			return nil, err
		}
	}

	return videos, nil
}

func (vr *videoRepository) GetVideosByKeywords(keywords, fromDate, toDate, uid string, pageNum, pageSize int64) ([]*model.Video, int64, error) {
	var videos []*model.Video
	var total int64
	var err error
	tx := vr.db.Model(&model.Video{}).
		Where("title LIKE ? or description LIKE ?", "%"+keywords+"%", "%"+keywords+"%")

	if fromDate != "" {
		tx = tx.Where("from_date > ?", fromDate)
	}
	if toDate != "" {
		tx = tx.Where("to_date < ?", toDate)
	}
	if uid != "" {
		tx = tx.Where("uid = ?", uid)
	}

	err = tx.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = tx.Offset((int(pageNum) - 1) * int(pageSize)).
		Limit(int(pageSize)).
		Find(&videos).Error

	if err != nil {
		return nil, 0, err
	}

	return videos, total, nil
}
