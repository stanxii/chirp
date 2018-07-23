package main

import (
	"flag"
	"fmt"
	"net/http"

	"chirp.com/api"
	"chirp.com/app"
	"chirp.com/email"
	"chirp.com/middleware"
)

func main() {
	boolPtr := flag.Bool("prod", false, "Provide this flag in production. This ensures that a .config file is provided before the application starts.")
	flag.Parse()
	cfg := app.LoadConfig(*boolPtr)
	services := app.Setup(cfg)
	defer services.Close()
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
