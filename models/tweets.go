package models

import (
	"time"

	"chirp.com/internal/utils"
	"github.com/jinzhu/gorm"
)

type Tweet struct {
	ID            uint      `gorm:"primary_key" json:"id"`
	Post          string    `gorm:"not_null" json:"post"`
	Username      string    `gorm:"not_null;index" json:"username"`
	Tags          []string  `gorm:"-" json:"tags,omitempty"`
	Taggings      []Tagging `json:"-"`
	LikesCount    uint      `json:"likesCount"`
	RetweetsCount uint      `json:"retweetsCount"`

	// IsRetweet bool
	Retweet   *Tweet `json:"retweet,omitempty"`
	RetweetID uint   `json:"retweetID,omitempty"`

	//tags
	tags []Tag `json:"tags"`

	// Images []Image `gorm:"-"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type TweetService interface {
	TweetDB
}

type tweetService struct {
	TweetDB
}

type TweetDB interface {
	ByID(id uint) (*Tweet, error)
	ByUsername(username string) ([]Tweet, error)
	ByUsernameAndRetweetID(username string, retweetID uint) (*Tweet, error)
	Create(tweet *Tweet) error
	Update(tweet *Tweet) error
	Delete(id uint) (*Tweet, error)
}

func NewTweetService(db *gorm.DB) TweetService {
	return &tweetService{
		TweetDB: &tweetValidator{&tweetGorm{db}},
	}
}

type tweetValidator struct {
	TweetDB
}

func (tv *tweetValidator) Create(tweet *Tweet) error {
	err := runTweetValFuncs(tweet,
		// tv.userIDRequired,
		tv.usernameRequired,
		tv.postRequired,
		tv.retweetOnlyOnce)
	if err != nil {
		return err
	}
	return tv.TweetDB.Create(tweet)
}

func (tv *tweetValidator) Update(tweet *Tweet) error {
	err := runTweetValFuncs(tweet,
		tv.usernameRequired,
		tv.postRequired)
	if err != nil {
		return err
	}
	return tv.TweetDB.Update(tweet)
}

// Delete will delete the tweet with the provided ID
func (tv *tweetValidator) Delete(id uint) (*Tweet, error) {
	if id <= 0 {
		return nil, ErrIDInvalid
	}
	tweet, err := tv.TweetDB.Delete(id)
	return tweet, err

}

func (tv *tweetValidator) usernameRequired(t *Tweet) error {
	if t.Username == "" {
		return ErrUsernameRequired
	}
	return nil
}

func (tv *tweetValidator) postRequired(t *Tweet) error {
	if t.RetweetID > 0 {
		return nil
	}
	if t.Post == "" {
		return ErrPostRequired
	}
	return nil
}

func (tv *tweetValidator) retweetOnlyOnce(t *Tweet) error {
	//check if this tweet is a retweet
	if t.RetweetID <= 0 {
		return nil
	}
	existing, err := tv.ByUsernameAndRetweetID(t.Username, t.RetweetID)
	if err == ErrNotFound {
		// tweet has not been retweeted by the user
		return nil
	}
	if err != nil {
		return err
	}
	if existing != nil {
		return ErrRetweetExists
	}
	return nil
}

var _ TweetDB = &tweetGorm{}

type tweetGorm struct {
	db *gorm.DB
}

func (tg *tweetGorm) ByID(id uint) (*Tweet, error) {
	var tweet Tweet
	db := tg.db.Where("id = ?", id)
	err := first(db, &tweet)
	return &tweet, err
}

func (tg *tweetGorm) ByUsername(username string) ([]Tweet, error) {
	var tweets []Tweet
	username = utils.NormalizeText(username)
	err := tg.db.Where("username = ?", username).Find(&tweets).Error
	if err != nil {
		return nil, err
	}
	return tweets, nil
}

func (tg *tweetGorm) ByUsernameAndRetweetID(username string, retweetID uint) (*Tweet, error) {
	var tweet Tweet
	db := tg.db.Where("username = ? AND retweet_id = ?", username, retweetID)
	err := first(db, &tweet)
	return &tweet, err

}

func (tg *tweetGorm) Create(tweet *Tweet) error {
	return tg.db.Create(tweet).Error
}

func (tg *tweetGorm) Update(tweet *Tweet) error {
	return tg.db.Save(tweet).Error
}

func (tg *tweetGorm) Delete(id uint) (*Tweet, error) {
	tweet := Tweet{ID: id}
	err := tg.db.Delete(&tweet).Error
	return &tweet, err
}

type tweetValFunc func(*Tweet) error

func runTweetValFuncs(tweet *Tweet, fns ...tweetValFunc) error {
	for _, fn := range fns {
		if err := fn(tweet); err != nil {
			return err
		}
	}
	return nil
}
