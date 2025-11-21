package service

import (
	"errors"
	"log"
	"west2/pkg/model"
	"west2/pkg/repository"
	"west2/util"

	"gorm.io/gorm"
)

type followService struct {
	fr repository.FollowRepostory
	ur repository.UserRepository
}

type FollowerService interface {
	FollowAction(follow *model.Follow) error
	GetFollowingList(followerId string, pageNum, pageSize int64) ([]*model.User, int64, error)
	GetFollowerList(followerId string, pageNum, pageSize int64) ([]*model.User, int64, error)
}

func NewFollowService(fr repository.FollowRepostory, ur repository.UserRepository) FollowerService {
	return &followService{fr: fr, ur: ur}
}

func (fs *followService) FollowAction(follow *model.Follow) error {
	f, err := fs.fr.GetFollowById(follow.FollowerId, follow.FollowingId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			follow.Id = util.GetID()
			err = fs.fr.Create(follow)
			if err != nil {
				log.Printf("failed to create follow: follow: %v, err: %v", follow, err)
				return err
			}
			return nil
		}
		log.Printf("failed to get follow by followerId and followingId: followerId: %s, followingId: %s, err: %v", follow.FollowerId, follow.FollowingId, err)
		return err
	}

	err = fs.fr.SetStatus(follow.Status, f.Id)
	if err != nil {
		log.Printf("failed to set follow's status by id: id: %s, status: %d, err: %v", f.Id, follow.Status, err)
		return err
	}
	return nil
}

func (fs *followService) GetFollowingList(followerId string, pageNum, pageSize int64) ([]*model.User, int64, error) {
	var users []*model.User
	follows, total, err := fs.fr.GetFollowingList(followerId, pageNum, pageSize)
	if err != nil {
		log.Printf("failed to get following list by follower id: followerId: %s, err: %v", followerId, err)
		return nil, 0, err
	}

	for _, f := range follows {
		u, err := fs.ur.GetUserById(f.FollowingId)
		if err != nil {
			log.Printf("failed to get user by id: uid: %s, err: %v", f.FollowingId, err)
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, nil
}

func (fs *followService) GetFollowerList(followerId string, pageNum, pageSize int64) ([]*model.User, int64, error) {
	var users []*model.User
	follows, total, err := fs.fr.GetFollowerList(followerId, pageNum, pageSize)
	if err != nil {
		log.Printf("failed to get following list by follower id: followerId: %s, err: %v", followerId, err)
		return nil, 0, err
	}

	for _, f := range follows {
		u, err := fs.ur.GetUserById(f.FollowerId)
		if err != nil {
			log.Printf("failed to get user by id: uid: %s, err: %v", f.FollowingId, err)
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, nil
}

func (fs *followService) GetFriendList(followerId string, pageNum, pageSize int64) ([]*model.User, int64, error) {
	var users []*model.User
	follows, total, err := fs.fr.GetFriendList(followerId, pageNum, pageSize)
	if err != nil {
		log.Printf("failed to get following list by follower id: followerId: %s, err: %v", followerId, err)
		return nil, 0, err
	}

	for _, f := range follows {
		u, err := fs.ur.GetUserById(f.FollowerId)
		if err != nil {
			log.Printf("failed to get user by id: uid: %s, err: %v", f.FollowingId, err)
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, nil
}
