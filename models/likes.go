package models

import (
	"github.com/jinzhu/gorm"
)

type Like struct {
	Tweet    *Tweet `json:"tweet"`
	TweetID  uint
	Username string
}

type LikeService interface {
	LikeDB
}

type likeService struct {
	LikeDB
}

func NewLikeService(db *gorm.DB) LikeService {
	lg := &likeGorm{db}
	return &likeService{
		LikeDB: lg,
	}
}

type LikeDB interface {
	// ByID(id uint) (*Like, error)
	// ByUserID(userID uint) ([]Like, error)
	Create(like *Like) error
	ByUsername(username string) ([]Like, error)
	Delete(id uint) error
	GetTotalLikes(id uint) uint
	GetUsers(id uint) ([]User, error)
}

type likeGorm struct {
	db *gorm.DB
}

var _ LikeDB = &likeGorm{}

func (lg *likeGorm) Create(like *Like) error {
	if like.Username == "" {
		return ErrUsernameRequired
	}
	return lg.db.Create(like).Error
}

// Delete will delete the user with the provided ID
func (lg *likeGorm) Delete(id uint) error {
	like := Like{TweetID: id}
	return lg.db.Delete(&like).Error
}

func (lg *likeGorm) ByUsername(username string) ([]Like, error) {
	var likes []Like
	err := lg.db.Preload("Tweet").Where("username = ?", username).Find(&likes).Error
	if err != nil {
		return nil, err
	}
	// err = lg.db.Model(&Like{}).Related(&)
	return likes, nil
}

func (lg *likeGorm) GetTotalLikes(id uint) uint {
	var count uint
	lg.db.Model(&Like{}).Where("tweet_id = ?", id).Count(&count)

	return count
}

func (lg *likeGorm) GetUsers(id uint) ([]User, error) {
	var users []User

	err := lg.db.Table("users").
		Select("users.username, users.name").
		Joins("JOIN likes ON users.username = likes.username AND likes.tweet_id = ?", id).
		Scan(&users).
		Error

	if err != nil {
		return nil, err
	}
	return users, nil
}
