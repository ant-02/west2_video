package repository

import (
	"west2/pkg/model"

	"gorm.io/gorm"
)

type followRepostory struct {
	db *gorm.DB
}

type FollowRepostory interface {
	Create(follow *model.Follow) error
	SetStatus(status int64, id string) error
	GetFollowingList(followerId string, pageNum, pageSize int64) ([]*model.Follow, int64, error)
	GetFollowerList(followingId string, pageNum, pageSize int64) ([]*model.Follow, int64, error)
	GetFriendList(followerId string, pageNum, pageSize int64) ([]*model.Follow, int64, error)
	GetFollowById(followerId, followingId string) (*model.Follow, error)
}

func NewFollowRepostory(db *gorm.DB) FollowRepostory {
	return &followRepostory{db: db}
}

func (fr *followRepostory) Create(follow *model.Follow) error {
	return fr.db.Create(follow).Error
}

func (fr *followRepostory) GetFollowById(followerId, followingId string) (*model.Follow, error) {
	var follow model.Follow
	err := fr.db.Where("follower_id = ?", followerId).
		Where("following_id = ?", followingId).
		First(&follow).Error
	if err != nil {
		return nil, err
	}
	return &follow, nil
}

func (fr *followRepostory) SetStatus(status int64, id string) error {
	return fr.db.Model(&model.Follow{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (fr *followRepostory) GetFollowingList(followerId string, pageNum, pageSize int64) ([]*model.Follow, int64, error) {
	var follows []*model.Follow
	var total int64
	err := fr.db.Transaction(func(tx *gorm.DB) error {
		tx = tx.Model(&model.Follow{}).
			Where("follower_id = ?", followerId).
			Where("status = 1")

		err := tx.Count(&total).Error
		if err != nil {
			return err
		}
		err = tx.Offset((int(pageNum) - 1) * int(pageSize)).
			Limit(int(pageSize)).
			Find(&follows).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}
	return follows, total, err
}

func (fr *followRepostory) GetFollowerList(followingId string, pageNum, pageSize int64) ([]*model.Follow, int64, error) {
	var follows []*model.Follow
	var total int64
	err := fr.db.Transaction(func(tx *gorm.DB) error {
		tx = tx.Model(&model.Follow{}).
			Where("following_id = ?", followingId).
			Where("status = 1")

		err := tx.Count(&total).Error
		if err != nil {
			return err
		}
		err = tx.Offset((int(pageNum) - 1) * int(pageSize)).
			Limit(int(pageSize)).
			Find(&follows).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}
	return follows, total, err
}

func (fr *followRepostory) GetFriendList(followerId string, pageNum, pageSize int64) ([]*model.Follow, int64, error) {
	var follows []*model.Follow
	var total int64
	err := fr.db.Transaction(func(tx *gorm.DB) error {
		tx = tx.Table("follows AS f1").
			Joins("INNER JOIN follows AS f2 ON f1.following_id = f2.follower_id AND f2.following_id = f1.follower_id").
			Where("f1.follower_id = ? AND f1.status = 1 AND f2.status = 1", followerId).
			Select("f1.*")
		err := tx.Count(&total).Error
		if err != nil {
			return err
		}
		err = tx.Offset((int(pageNum) - 1) * int(pageSize)).
			Limit(int(pageSize)).
			Find(&follows).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}
	return follows, total, err
}
