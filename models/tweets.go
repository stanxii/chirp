package models

import (
	"time"

	"chirp.com/utils"
	"github.com/jinzhu/gorm"
)

type Tweet struct {
	// gorm.Model
	ID         uint   `gorm:"primary_key" json:"id"`
	Post       string `gorm:"not_null" json:"post"`
	Username   string `gorm:"not_null;index" json:"username"`
	LikesCount uint   `json:"likesCount"`

	// Images []Image `gorm:"-"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type TweetService interface {
	TweetDB
}

type TweetDB interface {
	ByID(id uint) (*Tweet, error)
	ByUsername(username string) ([]Tweet, error)
	// ByUsernameAndID(username string, id uint) (*Tweet, error)
	// ByUserID(userID uint) ([]Tweet, error)
	Create(tweet *Tweet) error
	Update(tweet *Tweet) error
	Delete(id uint) error
}

func NewTweetService(db *gorm.DB) TweetService {
	return &tweetService{
		TweetDB: &tweetValidator{&tweetGorm{db}},
	}
}

type tweetService struct {
	TweetDB
}

type tweetValidator struct {
	TweetDB
}

func (tv *tweetValidator) Create(tweet *Tweet) error {
	err := runTweetValFuncs(tweet,
		// tv.userIDRequired,
		tv.usernameRequired,
		tv.postRequired)
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
func (tv *tweetValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}
	return tv.TweetDB.Delete(id)

}

// func (tv *tweetValidator) userIDRequired(t *Tweet) error {
// 	if t.UserID <= 0 {
// 		return ErrUserIDRequired
// 	}
// 	return nil
// }

func (tv *tweetValidator) usernameRequired(t *Tweet) error {
	if t.Username == "" {
		return ErrUsernameRequired
	}
	return nil
}

func (tv *tweetValidator) postRequired(t *Tweet) error {
	if t.Post == "" {
		return ErrPostRequired
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

// func (tg *tweetGorm) ByUsernameAndID(username string, id uint) (*Tweet, error) {
// 	var tweet Tweet
// 	db := tg.db.Where("username = ? AND id = ?", username, id)
// 	err := first(db, &tweet)
// 	return &tweet, err
// }

func (tg *tweetGorm) ByUsername(username string) ([]Tweet, error) {
	var tweets []Tweet
	username = utils.NormalizeText(username)
	err := tg.db.Where("username = ?", username).Find(&tweets).Error
	if err != nil {
		return nil, err
	}
	return tweets, nil
}

// func (tg *tweetGorm) ByUserID(userID uint) ([]Tweet, error) {
// 	var tweets []Tweet
// 	err := tg.db.Where("user_id = ?", userID).Find(&tweets).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return tweets, nil
// }

func (tg *tweetGorm) Create(tweet *Tweet) error {
	return tg.db.Create(tweet).Error
}

func (tg *tweetGorm) Update(tweet *Tweet) error {
	return tg.db.Save(tweet).Error
}

func (tg *tweetGorm) Delete(id uint) error {
	tweet := Tweet{ID: id}
	return tg.db.Delete(&tweet).Error
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
