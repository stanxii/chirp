package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"chirp.com/context"
	"chirp.com/email"
	"chirp.com/errors"
	"chirp.com/models"
	"chirp.com/pkg/rand"
	"github.com/gorilla/mux"
)

type Users struct {
	us      models.UserService
	ts      models.TweetService
	ls      models.LikeService
	emailer *email.Client
}

// NewUsers is used to create a new Users controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup.
func NewUsers(us models.UserService, ls models.LikeService, emailer *email.Client) *Users {
	return &Users{
		us:      us,
		ls:      ls,
		emailer: emailer,
	}
}

type SignUpForm struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`

	Password string `json:"password"`
}

func parseJSONForm(form interface{}, r *http.Request) error {
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&form)
	if err != nil {
		return err
	}
	return nil
}

// Create is used to process the signup form when a user
// submits it. This is used to create a new user account.
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignUpForm
	err := parseJSONForm(&form, r)
	if err != nil {
		RenderAPIError(w, errors.InvalidData(err))
		return
	}
	user := models.User{
		Name:     form.Name,
		Username: form.Username,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		RenderAPIError(w, errors.SetCustomError(err))
		return
	}
	u.emailer.Welcome(user.Name, user.Email)
	err = u.signIn(w, &user)
	if err != nil {
		RenderAPIError(w, errors.SetCustomError(err))
		return
	}
	RenderJSON(w, user, http.StatusOK)
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
	Render(w, user)
}

func (u *Users) getUser(w http.ResponseWriter, r *http.Request) *models.User {
	vars := mux.Vars(r)
	username := vars["username"]
	user, err := u.us.ByUsername(username)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			RenderAPIError(w, errors.NotFound("User"))
		default:
			log.Println(err)
			RenderAPIError(w, errors.InternalServerError(err))
		}
		return nil
	}
	// userWithLikes, err := u.us.UserDB.GetLikes(user)
	// err = u.us.AttachAssociations(user)
	return user
}

// Login is used to verify the provided email address and
// password and then log the user in if they are correct.
//
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var form LoginForm
	err := parseJSONForm(&form, r)
	if err != nil {
		RenderAPIError(w, errors.InvalidData(err))
		return
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			RenderAPIError(w, errors.NotFound("Email"))
		// need to add case where password is incorrect
		default:
			RenderAPIError(w, errors.SetCustomError(err))
		}
		return
	}

	err = u.signIn(w, user)
	if err != nil {
		RenderAPIError(w, errors.SetCustomError(err))
		return
	}

	// http.Redirect(w, r, "/tweets", http.StatusFound)
	RenderJSON(w, user, http.StatusOK)
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
	http.SetCookie(w, &cookie)
	return nil
}

// Logout is used to delete a users session cookie (remember_token)
// and then will update the user resource with a new remmeber
// token.
//
// POST /logout
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello from LOGOUT")
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
		RenderAPIError(w, errors.NotFound("User"))
		return
	}
	token, _ := rand.RememberToken()
	user.Remember = token
	u.us.Update(user)
}

// Get /:username
// func (u *Users) GetUser(w, http.ResponseWriter, r *http.Request) {

// }

// GET /:username/likes
func (u *Users) GetLikes(w http.ResponseWriter, r *http.Request) {
	user := u.getUser(w, r)
	if user == nil {
		return
	}
	likes, err := u.ls.ByUsername(user.Username)
	if err != nil {
		RenderAPIError(w, errors.SetCustomError(err))
	}
	user.Likes = likes
	RenderJSON(w, user, http.StatusOK)
}

/*
// ResetPwForm is used to process the forgot password form
// and the reset password form.
type ResetPwForm struct {
	Email    string `schema:"email"`
	Token    string `schema:"token"`
	Password string `schema:"password"`
}

// POST /forgot
func (u *Users) InitiateReset(w http.ResponseWriter, r *http.Request) {
	// TODO: Process the forgot password form and iniiate that process
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	token, err := u.us.InitiateReset(form.Email)
	if err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	err = u.emailer.ResetPw(form.Email, token)
	if err != nil {
		vd.SetAlert(err)
		u.ForgotPwView.Render(w, r, vd)
		return
	}

	views.RedirectAlert(w, r, "/reset", http.StatusFound, views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Instructions for resetting your password have been emailed to you.",
	})
}

// ResetPw displays the reset password form and has a method
// so that we can prefill the form data with a token provided
// via the URL query params.
//
// GET /reset
func (u *Users) ResetPw(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseURLParams(r, &form); err != nil {
		vd.SetAlert(err)
	}
	u.ResetPwView.Render(w, r, vd)
}

// CompleteReset processed the reset password form
//
// POST /reset
func (u *Users) CompleteReset(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPwForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.ResetPwView.Render(w, r, vd)
		return
	}

	user, err := u.us.CompleteReset(form.Token, form.Password)
	if err != nil {
		vd.SetAlert(err)
		u.ResetPwView.Render(w, r, vd)
		return
	}

	u.signIn(w, user)
	views.RedirectAlert(w, r, "/tweets", http.StatusFound, views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Your password has been reset and you have been logged in!",
	})
}


*/
