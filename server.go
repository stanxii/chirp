package main

import (
	"fmt"
	"net/http"

	"chirp.com/api"
	"chirp.com/app"
	"chirp.com/email"
	"chirp.com/internal/utils"
	"chirp.com/middleware"
	"chirp.com/models"
)

func main() {
	cfg := app.Init()
	dbCfg := cfg.Database
	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithTweet(),
		models.WithTag(),
		models.WithTagging(),
		models.WithLike(),
		models.WithFollow(),
	)
	utils.Must(err)
	defer services.Close()
	services.AutoMigrate()

	mgCfg := cfg.Mailgun
	emailer := email.NewClient(
		email.WithSender("Lenslocked.com Support", "support@mg.lenslocked.com"),
		email.WithMailgun(mgCfg.Domain, mgCfg.APIKey, mgCfg.PublicAPIKey),
	)

	router := app.NewRouter()

	tweetsAPI := api.NewTweets(services.Tweet, services.Like, services.Tag, services.Tagging)
	tagsAPI := api.NewTags(services.Tag, services.Tagging)
	usersAPI := api.NewUsers(services.User, services.Like, services.Follow, emailer)

	//init middleware
	userMw := middleware.NewUserMw(services.User)
	requireUserMw := middleware.NewRequireUserMw(userMw)

	//test route
	router.HandleFunc("/ping", ping).Methods("GET")
	//v1 api routes
	subRouter := router.PathPrefix("/v1").Subrouter()
	api.ServeUserResource(subRouter, usersAPI, &requireUserMw)
	api.ServeTweetResource(subRouter, tweetsAPI, &requireUserMw)
	api.ServeTagResource(subRouter, tagsAPI, &requireUserMw)

	http.ListenAndServe(":3000", userMw.Apply(router))
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Pinging the server...Success!")
}
