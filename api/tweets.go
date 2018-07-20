package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"chirp.com/context"
	"chirp.com/errors"
	"chirp.com/internal/utils"
	"chirp.com/middleware"
	"chirp.com/models"
	"chirp.com/pkg/unique"
	"github.com/gorilla/mux"
)

type Tweets struct {
	us       models.UserService
	ts       models.TweetService
	ls       models.LikeService
	tagS     models.TagService
	taggingS models.TaggingService
	// r        *mux.Router
}

func NewTweets(ts models.TweetService, ls models.LikeService, tagS models.TagService, taggingS models.TaggingService) *Tweets {
	return &Tweets{
		ts:       ts,
		ls:       ls,
		tagS:     tagS,
		taggingS: taggingS,
		// r:        r,
	}
}

func ServeTweetResource(r *mux.Router, t *Tweets, m *middleware.RequireUser) {
	r.HandleFunc("/i/tweets", m.ApplyFn(t.Index)).Methods("GET")
	r.HandleFunc("/tweets", m.ApplyFn(t.Create)).Methods("POST")
	r.HandleFunc("/tweets/{_username}/{id:[0-9]+}/delete", m.ApplyFn(t.Delete)).Methods("POST")
	r.HandleFunc("/{_username}/{id:[0-9]+}", t.Show).Methods("GET")
	r.HandleFunc("/{_username}/{id:[0-9]+}/update", m.ApplyFn(t.Update)).Methods("POST")
	r.HandleFunc("/{_username}/{id:[0-9]+}/like", m.ApplyFn(t.LikeTweet)).Methods("POST")
	r.HandleFunc("/{_username}/{id:[0-9]+}/like/delete", m.ApplyFn(t.DeleteLike)).Methods("POST")
	r.HandleFunc("/{_username}/{id:[0-9]+}/liked", t.GetUsersWhoLiked).Methods("GET")
	r.HandleFunc("/{_username}/{id:[0-9]+}/retweet", m.ApplyFn(t.CreateRetweet)).Methods("POST")
}

type TweetForm struct {
	Post string   `json:"post"`
	Tags []string `json:"tags"`
}

// POST /tweets
func (t *Tweets) Create(w http.ResponseWriter, r *http.Request) {
	var form TweetForm

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&form)
	if err != nil {
		utils.RenderAPIError(w, errors.InvalidData(err))
		return
	}
	user := context.User(r.Context())
	tweet := models.Tweet{
		Post:     form.Post,
		Username: user.Username,
		Tags:     form.Tags,
	}
	err = t.ts.Create(&tweet)
	if err != nil {
		utils.RenderAPIError(w, errors.SetCustomError(err))
		return
	}
	fmt.Println(tweet.ID)

	//create slice of unique and normalized tag names so we don't waste resources
	//querying duplicate tag names
	tweet.Tags = unique.Strings(tweet.Tags, utils.NormalizeText)
	fmt.Println(tweet.Tags)

	for _, name := range tweet.Tags {
		tag := &models.Tag{
			Name: name,
		}
		err := t.tagS.Create(tag)

		if err != nil && err != models.ErrTagExists {
			// fmt.Println(err)
			// tweet.Tags[i] = name + " is invalid: " + errors.SetCustomError(err).Message
			utils.RenderAPIError(w, errors.SetCustomError(err, tag))
			return
		}
		fmt.Println(tweet.ID)

		tagging := &models.Tagging{
			TweetID: tweet.ID,
			TagID:   tag.ID,
		}
		err = t.taggingS.Create(tagging)
		if err != nil {
			utils.RenderAPIError(w, errors.SetCustomError(err))
		}
	}
	utils.Render(w, &tweet)
}

// POST /tweets/:username/:id/delete
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

// Get /i/tweets
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

//GET /tweets/:username/:id
func (t *Tweets) Show(w http.ResponseWriter, r *http.Request) {
	tweet := t.tweetByID(w, r)
	if tweet == nil {
		return
	}
	utils.Render(w, tweet)
}

//POST /tweets/:username/:id/update
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
	err = t.ts.Update(tweet)
	if err != nil {
		utils.RenderAPIError(w, errors.SetCustomError(err))
		return
	}
	utils.Render(w, tweet)
}

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
		utils.RenderAPIError(w, errors.SetCustomError(err, &tweet))
		return
	}
	err = t.updateLikesCount(w, tweet)
	if err != nil {
		utils.RenderAPIError(w, errors.InternalServerError(err))
	}
	utils.Render(w, tweet)
}

//POST /:username/:id/like/delete
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
		utils.RenderAPIError(w, errors.SetCustomError(err, &retweet))
		return
	}
	utils.Render(w, retweet)
}

// func (t *Tweets) createTag(w http.RespnoseWriter, r *http.Request) {
// 	var tag models.Tag
// 	t.tagS.Create(&tag)
// }

/* HELPER METHODS */

func (t *Tweets) updateLikesCount(w http.ResponseWriter, tweet *models.Tweet) error {
	tweet.LikesCount = t.ls.GetTotalLikes(tweet.ID)
	err := t.ts.Update(tweet)
	if err != nil {
		return err
	}
	return nil
}

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
