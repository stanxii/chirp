package models

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"chirp.com/internal/utils"
	"github.com/jinzhu/gorm"
)

type Tag struct {
	ID        uint       `gorm:"primary_key" json:"-"`
	Name      string     `gorm:"not null;unique_index" json:"tagName,omitempty"`
	tweets    []Tweet    `json:"tweets"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type TagService interface {
	TagDB
}

type tagService struct {
	TagDB
}

type TagDB interface {
	// ByID(id uint) (*Tweet, error)
	// ByUsername(username string) ([]Tweet, error)
	// ByUsernameAndRetweetID(username string, retweetID uint) (*Tweet, error)
	// ByUsernameAndID(username string, id uint) (*Tweet, error)
	// ByUserID(userID uint) ([]Tweet, error)
	Create(tag *Tag) error
	ByName(name string) (*Tag, error)
	// Update(tweet *Tweet) error
	// Delete(id uint) (*Tweet, error)
	// CreateRetweet(tweet *Tweet) error
}

func NewTagService(db *gorm.DB) TagService {
	return &tagService{
		TagDB: &tagValidator{&tagGorm{db}},
	}
}

type tagValidator struct {
	TagDB
}

func (tv *tagValidator) Create(tag *Tag) error {
	err := runTagValFuncs(tag,
		tv.normalizeName,
		tv.nameRequired,
		tv.noSpecialCharacters,
		tv.noDuplicates,
	)

	if err != nil {
		return err
	}
	return tv.TagDB.Create(tag)
}

func (tv *tagValidator) ByName(name string) (*Tag, error) {
	tag := Tag{
		Name: name,
	}
	if err := runTagValFuncs(&tag, tv.normalizeName); err != nil {
		return nil, err
	}
	return tv.TagDB.ByName(tag.Name)
}

type tagValFunc func(*Tag) error

func runTagValFuncs(tag *Tag, fns ...tagValFunc) error {
	for _, fn := range fns {
		if err := fn(tag); err != nil {
			return err
		}
	}
	return nil
}

func (tv *tagValidator) nameRequired(t *Tag) error {
	if t.Name == "" {
		return ErrNameRequired
	}
	return nil
}

func (tv *tagValidator) normalizeName(t *Tag) error {
	t.Name = utils.NormalizeText(t.Name)
	return nil
}

func (tv *tagValidator) noDuplicates(t *Tag) error {
	existing, err := tv.ByName(t.Name)
	if err == ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}

	if existing.Name == t.Name {
		//sets tag to the existing tag
		fmt.Println("set")

		*t = *existing
		return ErrTagExists
	}
	return nil
}

func (tv *tagValidator) noSpecialCharacters(t *Tag) error {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	matched := reg.MatchString(t.Name)
	if matched {
		return ErrTagNoSpecialChar
	}
	return nil
}

var _ TagDB = &tagGorm{}

type tagGorm struct {
	db *gorm.DB
}

func (tg *tagGorm) Create(tag *Tag) error {
	return tg.db.Create(tag).Error
}

func (tg *tagGorm) ByName(name string) (*Tag, error) {
	var tag Tag
	db := tg.db.Where("name = ?", name)
	err := first(db, &tag)
	return &tag, err
}
