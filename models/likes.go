package models

import (
	"github.com/jinzhu/gorm"
)

type Like struct {
	Tweet   *Tweet `json:"tweet"`
	TweetID uint   `gorm:"primary_key" json:"-"`
	UserID  uint   `gorm:"primary_key" json:"-"`
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
		LikeDB: &likeValidator{lg},
	}
}

type likeValidator struct {
	LikeDB
}

type LikeDB interface {
	// ByID(id uint) (*Like, error)
	GetLike(id uint, userID uint) (*Like, error)
	Create(like *Like) error
	Delete(id uint, userID uint) error
	GetTotalLikes(id uint) uint
	GetUsers(id uint) ([]User, error)
	GetUserLikes(userID uint) ([]Tweet, error)
}

type likeValFunc func(*Like) error

func runLikeValFuncs(like *Like, fns ...likeValFunc) error {
	for _, fn := range fns {
		if err := fn(like); err != nil {
			return err
		}
	}
	return nil
}

func (lv *likeValidator) Create(like *Like) error {
	err := runLikeValFuncs(like, lv.noDuplicates)
	if err != nil {
		return err
	}
	return lv.LikeDB.Create(like)
}

func (lv *likeValidator) noDuplicates(l *Like) error {
	existing, err := lv.GetLike(l.TweetID, l.UserID)
	if err == ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}

	if existing.TweetID == l.TweetID && existing.UserID == l.UserID {
		return ErrLikeExists
	}
	return nil
}

type likeGorm struct {
	db *gorm.DB
}

var _ LikeDB = &likeGorm{}

func (lg *likeGorm) Create(like *Like) error {
	return lg.db.Create(like).Error
}

// Delete will delete the user with the provided ID
func (lg *likeGorm) Delete(id uint, userID uint) error {
	like := Like{TweetID: id, UserID: userID}
	return lg.db.Delete(&like).Error
}

func (lg *likeGorm) GetLike(id uint, userID uint) (*Like, error) {
	var like Like
	db := lg.db.Where("tweet_id = ? AND user_id = ? ", id, userID)
	err := first(db, &like)
	return &like, err
}

func (lg *likeGorm) GetUserLikes(userID uint) ([]Tweet, error) {
	var tweets []Tweet
	// err := lg.db.Preload("Tweet").Where("username = ?", username).Find(&likes).Error
	err := lg.db.Table("tweets").Joins("JOIN likes ON likes.tweet_id = tweets.id AND likes.user_id = ?", userID).Find(&tweets).Error
	if err != nil {
		return nil, err
	}
	return tweets, nil
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
		Joins("JOIN likes ON users.id = likes.user_id AND likes.tweet_id = ?", id).
		Find(&users).
		Error

	if err != nil {
		return nil, err
	}
	return users, nil
}
