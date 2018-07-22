package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"chirp.com/context"
	"chirp.com/errors"
	"chirp.com/internal/utils"
	"chirp.com/models"
)

type User struct {
	userService models.UserService
}

func (mw *User) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

//attaches User to context
func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// If the user is requesting a static asset or image
		// we will not need to lookup the current user so we skip
		// doing that.
		if strings.HasPrefix(path, "/assets/") ||
			strings.HasPrefix(path, "/images/") {
			next(w, r)
			return
		}
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			fmt.Println("No cookie!")

			next(w, r)
			return
		}
		fmt.Println(cookie.Value)

		user, err := mw.userService.ByRemember(cookie.Value)
		if err != nil {
			fmt.Println(" nothing on remember token")
			fmt.Println(err)

			next(w, r)
			return
		}
		fmt.Println(user)

		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next(w, r)

	})
}

// RequireUser assumes that User middleware has already been run
// otherwise it will not work correctly.
type RequireUser struct {
	User
}

// Apply assumes that User middleware has already been run
// otherwise it will no work correctly.
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn assumes that User middleware has already been run
// otherwise it will no work correctly.
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			fmt.Println("No user on require user")

			utils.RenderAPIError(w, errors.Unauthorized("You must be logged in to perform this action."))
			return
		}
		next(w, r)
	})
}

func NewUserMw(u models.UserService) User {
	return User{
		userService: u,
	}
}

func NewRequireUserMw(userMw User) RequireUser {
	return RequireUser{
		User: userMw,
	}
}
