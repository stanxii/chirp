package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"chirp.com/context"
	"chirp.com/email"
	"chirp.com/errors"
	"chirp.com/internal/utils"
	"chirp.com/middleware"
	"chirp.com/models"
	"chirp.com/pkg/rand"
	"github.com/gorilla/mux"
)

//the handler doesn't use {username} to look up the Tweet, but the user should be redirected to the correct username if the {username} doesn't match the Tweet's Username
func ServeUserResource(r *mux.Router, u *Users, m *middleware.RequireUser) {
	r.HandleFunc("/{username}", u.Show).Methods("GET")
	r.HandleFunc("/{username}/tweets", u.GetTweets).Methods("GET")
	r.HandleFunc("/{username}/likes", u.GetLikes).Methods("GET")
	r.HandleFunc("/{username}/followers", u.GetFollowers).Methods("GET")
	r.HandleFunc("/{username}/following", u.GetFollowing).Methods("GET")
	r.HandleFunc("/signup", u.Create).Methods("POST")
	r.HandleFunc("/login", u.Login).Methods("POST")
	r.HandleFunc("/logout", m.ApplyFn(u.Logout)).Methods("POST")
	r.HandleFunc("/{username}/follow", m.ApplyFn(u.FollowUser)).Methods("POST")
	r.HandleFunc("/{username}/follow/delete", m.ApplyFn(u.UnfollowUser)).Methods("POST")
}

type Users struct {
	us      models.UserService
	ts      models.TweetService
	ls      models.LikeService
	fs      models.FollowService
	emailer *email.Client
}

// NewUsers is used to create a new Users controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup.
func NewUsers(us models.UserService, ls models.LikeService, fs models.FollowService, ts models.TweetService, emailer *email.Client) *Users {
	return &Users{
		us:      us,
		ls:      ls,
		fs:      fs,
		ts:      ts,
		emailer: emailer,
	}
}

type SignUpForm struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Create is used to process the signup form when a user
// submits it. This is used to create a new user account.
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignUpForm
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&form)
	if err != nil {
		utils.RenderAPIError(w, errors.InvalidData(err))
		return
	}
	user := models.User{
		Name:     form.Name,
		Username: form.Username,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		utils.RenderAPIError(w, errors.SetCustomError(err, &user, ""))
		return
	}
	// u.emailer.Welcome(user.Name, user.Email)
	err = u.signIn(w, &user)
	if err != nil {
		utils.RenderAPIError(w, errors.SetCustomError(err, &user, ""))
		return
	}
	utils.Render(w, user)
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// /GET /:username
func (u *Users) Show(w http.ResponseWriter, r *http.Request) {
	user := u.getUser(w, r)
	if user == nil {
		return
	}
	utils.Render(w, user)
}

//GET /tweets/:username/tweets
func (u *Users) GetTweets(w http.ResponseWriter, r *http.Request) {
	user := u.getUser(w, r)
	if user == nil {
		return
	}
	tweets, err := u.ts.ByUsername(user.Username)
	if err != nil {
		utils.RenderAPIError(w, errors.InternalServerError(err))
		return
	}
	utils.Render(w, tweets)
}

func (u *Users) getUser(w http.ResponseWriter, r *http.Request) *models.User {
	vars := mux.Vars(r)
	username := vars["username"]
	user, err := u.us.ByUsername(username)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			utils.RenderAPIError(w, errors.NotFound("User"))
		default:
			log.Println(err)
			utils.RenderAPIError(w, errors.InternalServerError(err))
		}
		return nil
	}
	return user
}

// Login is used to verify the provided email address and
// password and then log the user in if they are correct.
//
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var form LoginForm
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&form)
	if err != nil {
		utils.RenderAPIError(w, errors.InvalidData(err))
		return
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			utils.RenderAPIError(w, errors.NotFound("Email"))
		// need to add case where password is incorrect
		default:
			utils.RenderAPIError(w, errors.SetCustomError(err, nil, ""))
		}
		return
	}

	err = u.signIn(w, user)
	if err != nil {
		utils.RenderAPIError(w, errors.SetCustomError(err, nil, ""))
		return
	}

	utils.Render(w, user)
}

// signIn is used to sign the given user in via cookies
func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	// fmt.Println("remember token: ", user.Remember)
	// fmt.Println("remember hash: ", user.RememberHash)

	http.SetCookie(w, &cookie)
	return nil
}

// Logout is used to delete a users session cookie (remember_token)
// and then will update the user resource with a new remmeber
// token.
//
// POST /logout
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	//updating the user's rememember token to ensure the user in inaccessable through the expired cookie
	user := context.User(r.Context())
	if user == nil {
		utils.RenderAPIError(w, errors.NotFound("User"))
		return
	}
	token, _ := rand.RememberToken()
	user.Remember = token
	u.us.Update(user)
}

// GET /:username/likes
func (u *Users) GetLikes(w http.ResponseWriter, r *http.Request) {
	user := u.getUser(w, r)
	if user == nil {
		return
	}
	likedTweets, err := u.ls.GetUserLikes(user.ID)
	if err != nil {
		utils.RenderAPIError(w, errors.SetCustomError(err, nil, ""))
	}
	user.LikedTweets = likedTweets
	utils.Render(w, user)
}

// POST /:username/follow
func (u *Users) FollowUser(w http.ResponseWriter, r *http.Request) {
	followee := u.getUser(w, r)
	if followee == nil {
		return
	}
	follower := context.User(r.Context())
	// //can't follow yourself
	if followee.ID == follower.ID {
		utils.RenderAPIError(w, errors.SetCustomError(models.ErrFollowSelf, followee, ""))
		return
	}
	follow := models.Follow{
		UserID:     followee.ID,
		User:       followee,
		FollowerID: follower.ID,
	}
	err := u.fs.Create(&follow)
	if err != nil {
		utils.RenderAPIError(w, errors.SetCustomError(err, followee, ""))
		return
	}
	// err = u.updateFollowCount(w, followee, follower)
	// if err != nil {
	// 	utils.RenderAPIError(w, errors.InternalServerError(err))
	// 	return
	// }
	utils.Render(w, &follow)

}

// POST /:username/follow/delete
func (u *Users) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	follower := context.User(r.Context())
	followee := u.getUser(w, r)
	if followee == nil {
		return
	}
	follow, err := u.fs.GetFollow(followee.ID, follower.ID)
	if err != nil {
		utils.RenderAPIError(w, errors.NotFound("Follow on this user"))
		return
	}
	err = u.fs.Delete(follow.UserID, follower.ID)
	if err != nil {
		utils.RenderAPIError(w, errors.InternalServerError(err))
		return
	}
	// err = u.updateFollowCount(w, followee, follower)
	// if err != nil {
	// 	utils.RenderAPIError(w, errors.InternalServerError(err))
	// 	return
	// }
	follow.User = followee
	utils.Render(w, followee)
}

// GET /:username/followers
func (u *Users) GetFollowers(w http.ResponseWriter, r *http.Request) {
	user := u.getUser(w, r)
	if user == nil {
		return
	}

	followers, err := u.fs.GetUserFollowers(user.ID)

	if err != nil {
		utils.RenderAPIError(w, errors.SetCustomError(err, nil, ""))
		return
	}
	user.Followers = followers
	utils.Render(w, user)
}

// GET /:username/following
func (u *Users) GetFollowing(w http.ResponseWriter, r *http.Request) {
	user := u.getUser(w, r)
	if user == nil {
		return
	}
	following, err := u.fs.GetUserFollowing(user.ID)
	if err != nil {
		utils.RenderAPIError(w, errors.SetCustomError(err, nil, ""))
		return
	}
	user.Following = following
	utils.Render(w, user)
}
