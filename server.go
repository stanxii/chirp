package main

import (
	"flag"
	"fmt"
	"net/http"

	"chirp.com/api"
	"chirp.com/email"
	"chirp.com/errors"
	"chirp.com/middleware"
	"chirp.com/models"
	"github.com/gorilla/mux"
	// "old/chirp.com/models"
)

func main() {
	boolPtr := flag.Bool("prod", false, "Provide this flag in production. This ensures that a .config file is provided before the application starts.")
	flag.Parse()

	cfg := LoadConfig(*boolPtr)
	dbCfg := cfg.Database

	// load error messages
	if err := errors.LoadMessages("config/errors.yaml"); err != nil {
		panic(fmt.Errorf("Failed to read the error message file: %s", err))
	}

	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithTweet(),
		models.WithLike(),
	)
	must(err)
	defer services.Close()
	services.AutoMigrate()

	mgCfg := cfg.Mailgun
	emailer := email.NewClient(
		email.WithSender("Lenslocked.com Support", "support@mg.lenslocked.com"),
		email.WithMailgun(mgCfg.Domain, mgCfg.APIKey, mgCfg.PublicAPIKey),
	)

	router := mux.NewRouter().StrictSlash(true)
	tweetsAPI := api.NewTweets(services.Tweet, services.Like, router)
	usersAPI := api.NewUsers(services.User, services.Like, emailer)

	//init middleware
	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{
		User: userMw,
	}

	//test route
	router.HandleFunc("/ping", ping).Methods("GET")

	//v1 api routes
	subRouter := router.PathPrefix("/v1").Subrouter()
	subRouter.HandleFunc("/i/tweets", requireUserMw.ApplyFn(tweetsAPI.Index)).Methods("GET")
	subRouter.HandleFunc("/tweets", requireUserMw.ApplyFn(tweetsAPI.Create)).Methods("POST")
	subRouter.HandleFunc("/signup", usersAPI.Create).Methods("POST")
	subRouter.HandleFunc("/login", usersAPI.Login).Methods("POST")
	subRouter.HandleFunc("/logout", requireUserMw.ApplyFn(usersAPI.Logout)).Methods("POST")

	//the handler doesn't use {_username} to look up the Tweet, but the user should be redirected to the correct username if the {_username} doesn't match the Tweet's Username
	subRouter.HandleFunc("/{_username}/{id:[0-9]+}", tweetsAPI.Show).Methods("GET")
	subRouter.HandleFunc("/{_username}/{id:[0-9]+}/update", requireUserMw.ApplyFn(tweetsAPI.Update)).Methods("POST")

	subRouter.HandleFunc("/{username}/likes", usersAPI.GetLikes).Methods("GET")
	subRouter.HandleFunc("/{_username}/{id:[0-9]+}/like", requireUserMw.ApplyFn(tweetsAPI.LikeTweet)).Methods("POST")
	subRouter.HandleFunc("/{_username}/{id:[0-9]+}/like/delete", requireUserMw.ApplyFn(tweetsAPI.RemoveLike)).Methods("POST")
	subRouter.HandleFunc("/{_username}/{id:[0-9]+}/liked", tweetsAPI.GetUsers).Methods("GET")

	http.ListenAndServe(":3000", userMw.Apply(router))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Pinging the server...Success!")
}
