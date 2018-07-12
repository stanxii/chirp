package models

import (
	"github.com/jinzhu/gorm"
)

type Tagging struct {
	TagID   uint `gorm:"primary_key" json:"-"`
	TweetID uint `gorm:"primary_key" json:"-"`
	Tag     *Tag `json:"tag"`
}

type TaggingService interface {
	TaggingDB
}

type taggingService struct {
	TaggingDB
}

func NewTaggingService(db *gorm.DB) TaggingService {
	tg := &taggingGorm{db}
	return &taggingService{
		TaggingDB: &taggingValidator{tg},
	}
}

type taggingValidator struct {
	TaggingDB
}

type TaggingDB interface {
	// ByID(id uint) (*Tagging, error)
	// GetTag(id uint, userID uint) (*Tagging, error)
	Create(tagging *Tagging) error
	GetTagging(tagID uint, tweetID uint) (*Tagging, error)
	GetTweets(id uint) ([]Tweet, error)
	// Delete(id uint, userID uint) error
	// GetTotalTags(id uint) uint
	// GetUsers(id uint) ([]User, error)
	// GetUserTags(userID uint) ([]Tweet, error)
}

type taggingValFunc func(*Tagging) error

func runTaggingValFuncs(tagging *Tagging, fns ...taggingValFunc) error {
	for _, fn := range fns {
		if err := fn(tagging); err != nil {
			return err
		}
	}
	return nil
}

func (tv *taggingValidator) Create(tagging *Tagging) error {
	err := runTaggingValFuncs(tagging, tv.noDuplicates)
	if err != nil {
		return err
	}
	return tv.TaggingDB.Create(tagging)
}

func (tv *taggingValidator) noDuplicates(t *Tagging) error {
	existing, err := tv.GetTagging(t.TagID, t.TweetID)
	if err == ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}

	if existing.TagID == t.TagID && existing.TweetID == t.TweetID {
		return ErrTaggingExists
	}
	return nil
}

// func (tv *taggingValidator) Create(tagging *Tagging) error {
// 	err := runTagValFuncs(tagging, tv.noDuplicates)
// 	if err != nil {
// 		return err
// 	}
// 	return tv.TaggingDB.Create(tagging)
// }

// func (tv *taggingValidator) noDuplicates(l *Tagging) error {
// 	existing, err := tv.GetTagging(l.TweetID, l.UserID)
// 	if err == ErrNotFound {
// 		return nil
// 	}
// 	if err != nil {
// 		return err
// 	}

// 	if existing.TweetID == l.TweetID && existing.UserID == l.UserID {
// 		return ErrTagExists
// 	}
// 	return nil
// }

type taggingGorm struct {
	db *gorm.DB
}

var _ TaggingDB = &taggingGorm{}

func (tg *taggingGorm) Create(tagging *Tagging) error {
	return tg.db.Create(tagging).Error
}

func (tg *taggingGorm) GetTagging(tagID uint, tweetID uint) (*Tagging, error) {
	var tagging Tagging
	db := tg.db.Where("tag_id = ? AND tweet_id = ? ", tagID, tweetID)
	err := first(db, &tagging)
	return &tagging, err
}

func (tg *taggingGorm) GetTweets(id uint) ([]Tweet, error) {
	var tweets []Tweet
	// err := lg.db.Preload("Tweet").Where("username = ?", username).Find(&likes).Error
	err := tg.db.Table("tweets").Joins("JOIN taggings ON taggings.tweet_id = tweets.id AND taggings.tag_id = ?", id).Scan(&tweets).Error
	if err != nil {
		return nil, err
	}
	return tweets, nil
}
