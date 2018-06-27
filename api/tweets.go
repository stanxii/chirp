package api

import (
	"log"
	"net/http"
	"strconv"

	"chirp.com/context"
	"chirp.com/errors"
	"chirp.com/models"
	"github.com/gorilla/mux"
)

type Tweets struct {
	us models.UserService
	ts models.TweetService
	ls models.LikeService
	r  *mux.Router
}

func NewTweets(ts models.TweetService, ls models.LikeService, r *mux.Router) *Tweets {
	return &Tweets{
		ts: ts,
		ls: ls,
		r:  r,
	}
}

type TweetForm struct {
	Post string `json:"post"`
}

// POST /tweets
func (t *Tweets) Create(w http.ResponseWriter, r *http.Request) {
	var form TweetForm
	err := parseJSONForm(&form, r)
	if err != nil {
		RenderAPIError(w, errors.InvalidData(err))
		return
	}
	user := context.User(r.Context())
	tweet := models.Tweet{
		Post:     form.Post,
		Username: user.Username,
	}
	err = t.ts.Create(&tweet)
	if err != nil {
		RenderAPIError(w, errors.SetCustomError(err))
		return
	}
	// RenderJSON(w, tweet, http.StatusOK)
	Render(w, &tweet)
}

// POST /tweets/:username/:id/delete
func (t *Tweets) Delete(w http.ResponseWriter, r *http.Request) {
	tweet := t.tweetByID(w, r)
	if tweet == nil {
		return
	}
	user := context.User(r.Context())
	if tweet.Username != user.Username {
		RenderAPIError(w, errors.Unauthorized())
		return
	}
	deletedTweet, err := t.ts.Delete(tweet.ID)
	if err != nil {
		RenderAPIError(w, errors.InternalServerError(err))
	}
	// RenderJSON(w, &deletedTweet, http.StatusOK)
	Render(w, deletedTweet)
}

// Get /i/tweets
func (t *Tweets) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	tweets, err := t.ts.ByUsername(user.Username)
	if err != nil {
		log.Println(err)
		RenderAPIError(w, errors.InternalServerError(err))
		return
	}
	RenderJSON(w, tweets, http.StatusOK)
}

//GET /tweets/:username/:id
func (t *Tweets) Show(w http.ResponseWriter, r *http.Request) {
	tweet := t.tweetByID(w, r)
	if tweet == nil {
		return
	}
	RenderJSON(w, tweet, http.StatusOK)

	// user := context.User(r.Context())
	// tweets, err := t.ts.ByUserID(user.ID)
	// tweet, err := t.tweetByID(w, r)
	// if err != nil {
	// 	return
	// }
}

//POST /tweets/:username/:id/update
func (t *Tweets) Update(w http.ResponseWriter, r *http.Request) {
	tweet := t.tweetByID(w, r)
	if tweet == nil {
		// RenderAPIError(w, errors.NotFound("Tweet"))
		return
	}
	user := context.User(r.Context())
	if tweet.Username != user.Username {
		RenderAPIError(w, errors.Unauthorized())
		return
	}
	var form TweetForm
	err := parseJSONForm(&form, r)
	if err != nil {
		RenderAPIError(w, errors.InvalidData(err))
		return
	}
	tweet.Post = form.Post
	err = t.ts.Update(tweet)
	if err != nil {
		RenderAPIError(w, errors.SetCustomError(err))
		return
	}
	RenderJSON(w, tweet, http.StatusOK)
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
		RenderAPIError(w, errors.SetCustomError(err, &tweet))
		return
	}
	err = t.updateLikesCount(w, tweet)
	if err != nil {
		RenderAPIError(w, errors.InternalServerError(err))
	}
	RenderJSON(w, tweet, http.StatusOK)
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
		RenderAPIError(w, errors.NotFound("Like on this tweet"))
		return
	}

	err = t.ls.Delete(like.TweetID, like.UserID)
	if err != nil {
		RenderAPIError(w, errors.InternalServerError(err))
		return
	}
	err = t.updateLikesCount(w, tweet)
	if err != nil {
		RenderAPIError(w, errors.InternalServerError(err))
	}
	RenderJSON(w, tweet, http.StatusOK)

}

func (t *Tweets) GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	tweet := t.tweetByID(w, r)
	if tweet == nil {
		return
	}
	users, err := t.ls.GetUsers(tweet.ID)
	if err != nil {
		RenderAPIError(w, errors.NotFound("Tweet"))
		return
	}
	RenderJSON(w, users, http.StatusOK)
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
		RenderAPIError(w, errors.SetCustomError(err, &retweet))
		return
	}
	RenderJSON(w, retweet, http.StatusOK)
}

// // POST /tweets/:username/:id/retweet/delete
// func (t *Tweets) DeleteRetweet(w http.ResponseWriter, r *http.Request) {
// 	tweet := t.tweetByID(w, r)
// 	if tweet == nil {
// 		return
// 	}
// 	user := context.User(r.Context())
// 	err := t.ts.Delete(&tweet)

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
	// username := utils.NormalizeText(vars["username"])
	idStr := vars["id"]
	idInt, err := strconv.Atoi(idStr)
	id := uint(idInt)
	if err != nil {
		log.Println(err)
		// http.Error(w, "Invalid tweet ID", http.StatusNotFound)
		RenderAPIError(w, errors.InvalidData(err))
		return nil
	}
	tweet, err := t.ts.ByID(id)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			// http.Error(w, "Tweet not found", http.StatusNotFound)
			RenderAPIError(w, errors.NotFound("Tweet"))

		default:
			log.Println(err)
			// http.Error(w, "Whoops! Something went wrong.", http.StatusInternalServerError)
			RenderAPIError(w, errors.InternalServerError(err))
		}
		return nil
	}
	// images, _ := t.is.ByTweetID(tweet.ID)
	// tweet.Images = images
	return tweet
}
