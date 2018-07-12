package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ServicesConfig func(*Services) error

func WithGorm(dialect, connectionInfo string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(dialect, connectionInfo)
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

func WithLogMode(mode bool) ServicesConfig {
	return func(s *Services) error {
		s.db.LogMode(mode)
		return nil
	}
}

func WithUser(pepper, hmacKey string) ServicesConfig {
	return func(s *Services) error {
		s.User = NewUserService(s.db, pepper, hmacKey)
		return nil
	}
}

func WithTweet() ServicesConfig {
	return func(s *Services) error {
		s.Tweet = NewTweetService(s.db)
		return nil
	}
}

func WithTag() ServicesConfig {
	return func(s *Services) error {
		s.Tag = NewTagService(s.db)
		return nil
	}
}

func WithTagging() ServicesConfig {
	return func(s *Services) error {
		s.Tagging = NewTaggingService(s.db)
		return nil
	}
}

func WithLike() ServicesConfig {
	return func(s *Services) error {
		s.Like = NewLikeService(s.db)
		return nil
	}
}

func WithFollow() ServicesConfig {
	return func(s *Services) error {
		s.Follow = NewFollowService(s.db)
		return nil
	}
}

// func WithImage() ServicesConfig {
// 	return func(s *Services) error {
// 		s.Image = NewImageService()
// 		return nil
// 	}
// }

func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

type Services struct {
	Tweet   TweetService
	User    UserService
	Like    LikeService
	Follow  FollowService
	Tag     TagService
	Tagging TaggingService
	// Image ImageService
	db *gorm.DB
}

// Closes the database connection
func (s *Services) Close() error {
	return s.db.Close()
}

// // DestructiveReset drops all tables and rebuilds them
// func (s *Services) DestructiveReset() error {
// 	err := s.db.DropTableIfExists(&User{}, &Tweet{}, &Like{}, &pwReset{}).Error
// 	if err != nil {
// 		return err
// 	}
// 	return s.AutoMigrate()
// }

// // AutoMigrate will attempt to automatically migrate all tables
// func (s *Services) AutoMigrate() error {
// 	return s.db.AutoMigrate(&User{}, &Tweet{}, &Like{}, &pwReset{}).Error
// }

// DestructiveReset drops all tables and rebuilds them
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Tweet{}, &Like{}, &Follow{}, &Tag{}, &Tagging{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate all tables
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Tweet{}, &Like{}, &Follow{}, &Tag{}, &Tagging{}).Error
}
