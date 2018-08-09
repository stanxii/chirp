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
	Create(tagging *Tagging) error
	GetTagging(tagID uint, tweetID uint) (*Tagging, error)
	GetTaggings(tweetID uint) ([]Tagging, error)
	GetTweets(id uint) ([]Tweet, error)
	Delete(tagID, tweetID uint) error
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

type taggingGorm struct {
	db *gorm.DB
}

var _ TaggingDB = &taggingGorm{}

func (tg *taggingGorm) Create(tagging *Tagging) error {
	return tg.db.Create(tagging).Error
}

func (tg *taggingGorm) Delete(tagID, tweetID uint) error {
	tagging := Tagging{TagID: tagID, TweetID: tweetID}
	return tg.db.Delete(&tagging).Error
}

func (tg *taggingGorm) GetTagging(tagID uint, tweetID uint) (*Tagging, error) {
	var tagging Tagging
	db := tg.db.Where("tag_id = ? AND tweet_id = ? ", tagID, tweetID)
	err := first(db, &tagging)
	return &tagging, err
}

func (tg *taggingGorm) GetTaggings(tweetID uint) ([]Tagging, error) {
	var taggings []Tagging
	err := tg.db.Where("tweet_id = ?", tweetID).Find(&taggings).Error
	if err != nil {
		return nil, err
	}
	return taggings, nil
}

func (tg *taggingGorm) GetTweets(id uint) ([]Tweet, error) {
	var tweets []Tweet
	// err := lg.db.Preload("Tweet").Where("username = ?", username).Find(&likes).Error
	err := tg.db.Table("tweets").Joins("JOIN taggings ON taggings.tweet_id = tweets.id").Where("taggings.tag_id = ?", id).Find(&tweets).Error
	if err != nil {
		return nil, err
	}
	return tweets, nil
}
