package model

import (
	"time"
	"west2/biz/model/follow"
)

type Follow struct {
	Id          string    `gorm:"type:varchar(100);primaryKey"`
	FollowingId string    `gorm:"type:varchar(100);not null"`
	FollowerId  string    `gorm:"type:varchar(100);not null"`
	Status      int64     `gorm:"type:int(2);noy null;default:0"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	DeletedAt   time.Time `gorm:"default:null"`
}

func UserToFollowUser(u *User) *follow.User {
	return &follow.User{
		Id:        u.Id,
		Username:  u.Username,
		AvatarUrl: u.AvatarUrl,
	}
}

func UsersToFollowUsers(u []*User) []*follow.User {
	if u == nil {
		return nil
	}

	var users []*follow.User
	for _, x := range u {
		users = append(users, UserToFollowUser(x))
	}
	return users
}
