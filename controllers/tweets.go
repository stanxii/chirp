package controllers

import (
	"chirp.com/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"chirp.com/context"
	"chirp.com/errors"
	"chirp.com/internal/utils"
	"chirp.com/middleware"
	"chirp.com/pkg/unique"
	"github.com/gorilla/mux"
)

func ServeTweetResource(r *mux.Router, t *Tweets, m *middleware.RequireUser) {
	r.HandleFunc("/tweets", m.ApplyFn(t.Create)).Methods("POST")
	r.HandleFunc("/tweets/{_username}/{id:[0-9]+}/delete", m.ApplyFn(t.Delete)).Methods("POST")
	r.HandleFunc("/{_username}/{id:[0-9]+}", t.Show).Methods("GET")
	r.HandleFunc("/{_username}/{id:[0-9]+}/update", m.ApplyFn(t.Update)).Methods("POST")
	r.HandleFunc("/{_username}/{id:[0-9]+}/like", m.ApplyFn(t.LikeTweet)).Methods("POST")
	r.HandleFunc("/{_username}/{id:[0-9]+}/like/delete", m.ApplyFn(t.DeleteLike)).Methods("POST")
	r.HandleFunc("/{_username}/{id:[0-9]+}/liked", t.GetUsersWhoLiked).Methods("GET")
	r.HandleFunc("/{_username}/{id:[0-9]+}/retweet", m.ApplyFn(t.CreateRetweet)).Methods("POST")
}

type Tweets struct {
	us       models.UserService
	ts       models.TweetService
	ls       models.LikeService
	tagS     models.TagService
	taggingS models.TaggingService
}

func NewTweets(ts models.TweetService, ls models.LikeService, tagS models.TagService, taggingS models.TaggingService) *Tweets {
	return &Tweets{
		ts:       ts,
		ls:       ls,
		tagS:     tagS,
		taggingS: taggingS,
	}
}

type TweetForm struct {
	Post string   `json:"post"`
	Tags []string `json:"tags"`
}

/*
Creates tweet
 */
func (t *Tweets) Create(w http.ResponseWriter, r *http.Request) {
	var form TweetForm

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&form)
	if err != nil {
		utils.RenderAPIError(w, errors.InvalidData(err))
		return
	}
	user := context.User(r.Context())
	uniqueTags := unique.Strings(form.Tags, utils.NormalizeText)
	tweet := models.Tweet{
		Post:     form.Post,
		Username: user.Username,
		Tags:     uniqueTags,
	}
	err = t.ts.Create(&tweet)
	if err != nil {
		utils.RenderAPIError(w, errors.SetCustomError(err, nil, ""))
		return
	}

	t.createTags(w, &tweet)
	utils.Render(w, &tweet)
}

/*
Creates the tags in the given tweet
 */
func (t *Tweets) createTags(w http.ResponseWriter, tweet *models.Tweet) {
	for _, name := range tweet.Tags {
		tag := &models.Tag{
			Name: name,
		}
		err := t.tagS.Create(tag)

		if err != nil && err != models.ErrTagExists {
			utils.RenderAPIError(w, errors.SetCustomError(err, tag, ""))
			return
		}

		err = t.createTagging(tweet, tag)
		if err != nil {
			if pErr, ok := err.(errors.PublicError); ok {
				fmt.Println("Tagging Service: ", pErr.Public(), ": ", tag.Name)
			} else {
				log.Println(err)
			}
			return
		}

	}
}

/*
 Create taggings associated with the tweet and tag
 */
func (t *Tweets) createTagging(tweet *models.Tweet, tag *models.Tag) error {
	tagging := &models.Tagging{
		TweetID: tweet.ID,
		TagID:   tag.ID,
	}
	err := t.taggingS.Create(tagging)
	if err != nil {
		return err
	}
	return nil
}

/*
Deletes the tagging between tweet and tag
 */
func (t *Tweets) deleteTagging(tweet *models.Tweet, tag *models.Tag) error {
	err := t.taggingS.Delete(tag.ID, tweet.ID)
	if err != nil {
		return err
	}
	return nil
}

/*
Deletes the tweet
 */
func (t *Tweets) Delete(w http.ResponseWriter, r *http.Request) {
	tweet := t.tweetByID(w, r)
	if tweet == nil {
		return
	}
	user := context.User(r.Context())
	if tweet.Username != user.Username {
		utils.RenderAPIError(w, errors.Unauthorized())
		return
	}
	deletedTweet, err := t.ts.Delete(tweet.ID)
	if err != nil {
		utils.RenderAPIError(w, errors.InternalServerError(err))
	}
	utils.Render(w, deletedTweet)
}

/*
Get tweets posted by the active user
 */
func (t *Tweets) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	tweets, err := t.ts.ByUsername(user.Username)
	if err != nil {
		log.Println(err)
		utils.RenderAPIError(w, errors.InternalServerError(err))

		return
	}
	utils.Render(w, tweets)
}

/*
Get the tweet from the given user
 */
func (t *Tweets) Show(w http.ResponseWriter, r *http.Request) {
	tweet := t.tweetByID(w, r)
	if tweet == nil {
		return
	}
	utils.Render(w, tweet)
}

/*
Updates the tweet
 */
func (t *Tweets) Update(w http.ResponseWriter, r *http.Request) {
	tweet := t.tweetByID(w, r)
	if tweet == nil {
		return
	}
	user := context.User(r.Context())
	if tweet.Username != user.Username {
		utils.RenderAPIError(w, errors.Unauthorized())
		return
	}
	var form TweetForm
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&form)
	if err != nil {
		utils.RenderAPIError(w, errors.InvalidData(err))
		return
	}
	tweet.Post = form.Post
	err = t.updateTags(tweet, w, form)
	if err != nil {
		return
	}

	err = t.ts.Update(tweet)
	if err != nil {
		utils.RenderAPIError(w, errors.SetCustomError(err, nil, ""))
		return
	}
	t.createTags(w, tweet)
	utils.Render(w, tweet)
}

/*
Update the tags associated with the tweet
 */
func (t *Tweets) updateTags(tweet *models.Tweet, w http.ResponseWriter, form TweetForm) error {
	newTags := unique.Strings(form.Tags, utils.NormalizeText)
	taggings, err := t.taggingS.GetTaggings(tweet.ID)
	if err != nil {
		utils.RenderAPIError(w, errors.InternalServerError(err))
		return err
	}
	var oldTags []string
	for _, tagging := range taggings {
		tag, err := t.tagS.ByID(tagging.TagID)
		if err != nil {
			fmt.Println(tagging.TagID, ": ", err)
			continue
		}
		oldTags = append(oldTags, tag.Name)
	}
	taggingsToDelete := diff(newTags, oldTags)

	for _, tagname := range taggingsToDelete {
		tag, err := t.tagS.ByName(tagname)
		if err != nil {
			utils.RenderAPIError(w, errors.NotFound(tagname))
			return err
		}

		err = t.deleteTagging(tweet, tag)
		if err != nil {
			utils.RenderAPIError(w, errors.SetCustomError(err, nil, ""))
			return err
		}
	}
	tweet.Tags = newTags
	return nil
}

/*
Get differences between the two slices
Note: Using inefficent algo at the moment. Can be optimized using a map.
 */
func diff(a, b []string) (diff []string) {
	for _, bStr := range b {
		var match bool
		for _, aStr := range a {
			if aStr == bStr {
				match = true
				break
			}
		}
		if !match {
			diff = append(diff, bStr)
		}
	}
	return diff
}

/*
Add a like to the tweet
 */
func (t *Tweets) LikeTweet(w http.ResponseWriter, r *http.Request) {
	tweet := t.tweetByID(w, r)
	if tweet == nil {
		return
	}
	user := context.User(r.Context())
	like := models.Like{
		UserID:  user.ID,
		TweetID: tweet.ID,
	}
	err := t.ls.Create(&like)
	if err != nil {
		utils.RenderAPIError(w, errors.SetCustomError(err, &tweet, ""))
		return
	}
	err = t.updateLikesCount(w, tweet)
	if err != nil {
		utils.RenderAPIError(w, errors.InternalServerError(err))
	}
	utils.Render(w, tweet)
}

/*
Remove like from the tweet
 */
func (t *Tweets) DeleteLike(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	tweet := t.tweetByID(w, r)
	if tweet == nil {
		return
	}

	like, err := t.ls.GetLike(tweet.ID, user.ID)
	if err != nil {
		utils.RenderAPIError(w, errors.NotFound("Like on this tweet"))
		return
	}

	err = t.ls.Delete(like.TweetID, like.UserID)
	if err != nil {
		utils.RenderAPIError(w, errors.InternalServerError(err))
		return
	}
	err = t.updateLikesCount(w, tweet)
	if err != nil {
		utils.RenderAPIError(w, errors.InternalServerError(err))
	}
	utils.Render(w, tweet)

}

func (t *Tweets) GetUsersWhoLiked(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	tweet := t.tweetByID(w, r)
	if tweet == nil {
		return
	}
	users, err := t.ls.GetUsers(tweet.ID)
	if err != nil {
		utils.RenderAPIError(w, errors.NotFound("Tweet"))
		return
	}
	utils.Render(w, users)
}

// POST /tweets/:username/:id/retweet
func (t *Tweets) CreateRetweet(w http.ResponseWriter, r *http.Request) {
	tweet := t.tweetByID(w, r)
	if tweet == nil {
		return
	}
	user := context.User(r.Context())

	retweet := models.Tweet{
		Username:  user.Username,
		Retweet:   tweet,
		RetweetID: tweet.ID,
	}
	err := t.ts.Create(&retweet)
	if err != nil {
		utils.RenderAPIError(w, errors.SetCustomError(err, &retweet, ""))
		return
	}
	utils.Render(w, retweet)
}

/* HELPER METHODS */

/*
Updates number of likes on the tweet
 */
func (t *Tweets) updateLikesCount(w http.ResponseWriter, tweet *models.Tweet) error {
	tweet.LikesCount = t.ls.GetTotalLikes(tweet.ID)
	err := t.ts.Update(tweet)
	if err != nil {
		return err
	}
	return nil
}

/*
Get tweet by tweet's ID
 */
func (t *Tweets) tweetByID(w http.ResponseWriter, r *http.Request) *models.Tweet {
	vars := mux.Vars(r)
	idStr := vars["id"]
	idInt, err := strconv.Atoi(idStr)
	id := uint(idInt)
	if err != nil {
		log.Println(err)
		utils.RenderAPIError(w, errors.InvalidData(err))
		return nil
	}
	tweet, err := t.ts.ByID(id)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			// http.Error(w, "Tweet not found", http.StatusNotFound)
			utils.RenderAPIError(w, errors.NotFound("Tweet"))

		default:
			log.Println(err)
			// http.Error(w, "Whoops! Something went wrong.", http.StatusInternalServerError)
			utils.RenderAPIError(w, errors.InternalServerError(err))
		}
		return nil
	}
	return tweet
}
