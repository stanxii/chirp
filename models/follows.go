package models

import (
	"github.com/jinzhu/gorm"
)

type Follow struct {
	// Tweet    *Tweet `json:"tweet"`
	// TweetID  uint   `json:"-"`
	// Username string `json:"-"`
	FollowerID uint  `json:"follower_id" gorm:"primary_key"`
	UserID     uint  `json:"user_id" gorm:"primary_key"`
	User       *User `json:"user"`
}

type FollowService interface {
	FollowDB
}

type followService struct {
	FollowDB
}

func NewFollowService(db *gorm.DB) FollowService {
	fg := &followGorm{db}
	return &followService{
		FollowDB: &followValidator{fg},
	}
}

type followValidator struct {
	FollowDB
}

type FollowDB interface {
	// ByID(id uint) (*Follow, error)
	Create(follow *Follow) error
	GetFollow(userID uint, followerID uint) (*Follow, error)
	GetUserFollowers(id uint) ([]User, error)
	GetUserFollowing(id uint) ([]User, error)
	Delete(userID uint, followerID uint) error
	GetTotalFollowers(id uint) uint
	GetTotalFollowing(id uint) uint
	// GetUsers(id uint) ([]User, error)
	// GetUserFollows(username string) ([]Tweet, error)
}

type followValFunc func(*Follow) error

func runFollowValFuncs(follow *Follow, fns ...followValFunc) error {
	for _, fn := range fns {
		if err := fn(follow); err != nil {
			return err
		}
	}
	return nil
}

func (fv *followValidator) Create(follow *Follow) error {
	err := runFollowValFuncs(follow, fv.noDuplicates)
	if err != nil {
		return err
	}
	return fv.FollowDB.Create(follow)
}

func (fv *followValidator) noDuplicates(f *Follow) error {
	existing, err := fv.GetFollow(f.UserID, f.FollowerID)
	if err == ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}

	if existing.UserID == f.UserID && existing.FollowerID == f.FollowerID {
		return ErrFollowExists
	}
	return nil
}

type followGorm struct {
	db *gorm.DB
}

var _ FollowDB = &followGorm{}

func (fg *followGorm) GetUserFollowers(userID uint) ([]User, error) {
	var users []User
	err := fg.db.Table("users").Joins("JOIN follows ON follows.follower_id = id AND follows.user_id = ?", userID).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (fg *followGorm) GetUserFollowing(userID uint) ([]User, error) {
	var users []User
	err := fg.db.Table("users").Joins("JOIN follows ON follows.user_id = id AND follows.follower_id = ?", userID).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (fg *followGorm) GetFollow(userID uint, followerID uint) (*Follow, error) {
	var follow Follow
	db := fg.db.Where("user_id = ? AND follower_id = ?", userID, followerID)
	err := first(db, &follow)
	return &follow, err
}

func (fg *followGorm) Create(follow *Follow) error {
	return fg.db.Create(follow).Error
}

func (fg *followGorm) Delete(userID uint, followerID uint) error {
	follow := Follow{UserID: userID, FollowerID: followerID}
	return fg.db.Delete(&follow).Error
}

func (fg *followGorm) GetTotalFollowers(id uint) uint {
	var count uint
	fg.db.Model(&Follow{}).Where("user_id = ?", id).Count(&count)
	return count
}

func (fg *followGorm) GetTotalFollowing(id uint) uint {
	var count uint
	fg.db.Model(&Follow{}).Where("follower_id = ?", id).Count(&count)
	return count
}
